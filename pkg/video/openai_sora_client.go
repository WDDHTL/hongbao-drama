package video

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"
)

type OpenAISoraClient struct {
	BaseURL       string
	APIKey        string
	Model         string
	Endpoint      string
	QueryEndpoint string
	HTTPClient    *http.Client
}

type OpenAISoraResponse struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	Model       string `json:"model"`
	Status      string `json:"status"`
	Progress    int    `json:"progress"`
	CreatedAt   int64  `json:"created_at"`
	CompletedAt int64  `json:"completed_at"`
	Size        string `json:"size"`
	Seconds     string `json:"seconds"`
	Quality     string `json:"quality"`
	VideoURL    string `json:"video_url"`
	Video       struct {
		URL string `json:"url"`
	} `json:"video"`
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

func NewOpenAISoraClient(baseURL, apiKey, model, endpoint, queryEndpoint string) *OpenAISoraClient {
	if endpoint == "" {
		endpoint = "/videos"
	}
	if queryEndpoint == "" {
		queryEndpoint = "/videos/{taskId}"
	}

	return &OpenAISoraClient{
		BaseURL:       baseURL,
		APIKey:        apiKey,
		Model:         model,
		Endpoint:      endpoint,
		QueryEndpoint: queryEndpoint,
		HTTPClient: &http.Client{
			Timeout: 300 * time.Second,
		},
	}
}

func joinURL(baseURL, endpoint string) string {
	if endpoint == "" {
		return baseURL
	}

	trimmedBase := strings.TrimRight(baseURL, "/")
	trimmedEndpoint := strings.TrimLeft(endpoint, "/")
	if trimmedBase == "" {
		return "/" + trimmedEndpoint
	}
	return trimmedBase + "/" + trimmedEndpoint
}

func resolveTaskEndpoint(endpoint, taskID string) string {
	switch {
	case strings.Contains(endpoint, "{taskId}"):
		return strings.ReplaceAll(endpoint, "{taskId}", taskID)
	case strings.Contains(endpoint, "{task_id}"):
		return strings.ReplaceAll(endpoint, "{task_id}", taskID)
	case strings.HasSuffix(endpoint, "/"+taskID):
		return endpoint
	default:
		return strings.TrimRight(endpoint, "/") + "/" + taskID
	}
}

func (c *OpenAISoraClient) GenerateVideo(imageURL, prompt string, opts ...VideoOption) (*VideoResult, error) {
	options := &VideoOptions{
		Duration: 4,
	}

	for _, opt := range opts {
		opt(options)
	}

	model := c.Model
	if options.Model != "" {
		model = options.Model
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("model", model); err != nil {
		return nil, fmt.Errorf("write model field: %w", err)
	}
	if err := writer.WriteField("prompt", prompt); err != nil {
		return nil, fmt.Errorf("write prompt field: %w", err)
	}

	if options.Duration > 0 {
		if err := writer.WriteField("seconds", fmt.Sprintf("%d", options.Duration)); err != nil {
			return nil, fmt.Errorf("write duration field: %w", err)
		}
	}

	if options.Resolution != "" {
		if err := writer.WriteField("size", options.Resolution); err != nil {
			return nil, fmt.Errorf("write size field: %w", err)
		}
	}

	primaryReference := imageURL
	if primaryReference == "" {
		switch {
		case options.FirstFrameURL != "":
			primaryReference = options.FirstFrameURL
		case len(options.ReferenceImageURLs) > 0:
			primaryReference = options.ReferenceImageURLs[0]
		case options.LastFrameURL != "":
			primaryReference = options.LastFrameURL
		}
	}

	if primaryReference != "" {
		var imageData []byte
		mimeType := "image/png"
		filename := "reference.png"

		if strings.HasPrefix(primaryReference, "data:") {
			parts := strings.SplitN(primaryReference, ",", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid data URI format")
			}

			header := parts[0]
			switch {
			case strings.Contains(header, "image/jpeg"), strings.Contains(header, "image/jpg"):
				mimeType = "image/jpeg"
				filename = "reference.jpg"
			case strings.Contains(header, "image/webp"):
				mimeType = "image/webp"
				filename = "reference.webp"
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				return nil, fmt.Errorf("failed to decode base64 image: %w", err)
			}
			imageData = decoded
		} else {
			resp, err := http.Get(primaryReference)
			if err != nil {
				return nil, fmt.Errorf("failed to download reference image: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("failed to download reference image, status: %d", resp.StatusCode)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read downloaded image: %w", err)
			}
			imageData = data

			mimeType = resp.Header.Get("Content-Type")
			if mimeType == "" || mimeType == "application/octet-stream" {
				switch strings.ToLower(filepath.Ext(primaryReference)) {
				case ".jpg", ".jpeg":
					mimeType = "image/jpeg"
				case ".webp":
					mimeType = "image/webp"
				default:
					mimeType = "image/png"
				}
			}

			base := filepath.Base(primaryReference)
			if base != "" && base != "." {
				if idx := strings.Index(base, "?"); idx != -1 {
					base = base[:idx]
				}
				filename = base
			}
		}

		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="input_reference"; filename="%s"`, filename))
		header.Set("Content-Type", mimeType)

		part, err := writer.CreatePart(header)
		if err != nil {
			return nil, fmt.Errorf("create part: %w", err)
		}
		if _, err := part.Write(imageData); err != nil {
			return nil, fmt.Errorf("write image data: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	endpoint := joinURL(c.BaseURL, c.Endpoint)
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result OpenAISoraResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if result.Error.Message != "" {
		return nil, fmt.Errorf("openai error: %s", result.Error.Message)
	}

	videoResult := &VideoResult{
		TaskID:    result.ID,
		Status:    result.Status,
		Completed: result.Status == "completed",
	}

	if result.VideoURL != "" {
		videoResult.VideoURL = result.VideoURL
	} else if result.Video.URL != "" {
		videoResult.VideoURL = result.Video.URL
	}

	return videoResult, nil
}

func (c *OpenAISoraClient) GetTaskStatus(taskID string) (*VideoResult, error) {
	endpoint := joinURL(c.BaseURL, resolveTaskEndpoint(c.QueryEndpoint, taskID))
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result OpenAISoraResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	videoResult := &VideoResult{
		TaskID:    result.ID,
		Status:    result.Status,
		Completed: result.Status == "completed",
	}

	if result.Error.Message != "" {
		videoResult.Error = result.Error.Message
	}

	if result.VideoURL != "" {
		videoResult.VideoURL = result.VideoURL
	} else if result.Video.URL != "" {
		videoResult.VideoURL = result.Video.URL
	}

	return videoResult, nil
}
