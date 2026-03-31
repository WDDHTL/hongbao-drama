package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/pkg/ai"
	"github.com/drama-generator/backend/pkg/logger"
	"gorm.io/gorm"
)

type AIService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewAIService(db *gorm.DB, log *logger.Logger) *AIService {
	return &AIService{
		db:  db,
		log: log,
	}
}

type CreateAIConfigRequest struct {
	ServiceType   string            `json:"service_type" binding:"required,oneof=text image video"`
	Name          string            `json:"name" binding:"required,min=1,max=100"`
	Provider      string            `json:"provider" binding:"required"`
	BaseURL       string            `json:"base_url" binding:"required,url"`
	APIKey        string            `json:"api_key" binding:"required"`
	Model         models.ModelField `json:"model" binding:"required"`
	Endpoint      string            `json:"endpoint"`
	QueryEndpoint string            `json:"query_endpoint"`
	Priority      int               `json:"priority"`
	IsDefault     bool              `json:"is_default"`
	Settings      string            `json:"settings"`
}

type UpdateAIConfigRequest struct {
	Name          string             `json:"name" binding:"omitempty,min=1,max=100"`
	Provider      string             `json:"provider"`
	BaseURL       string             `json:"base_url" binding:"omitempty,url"`
	APIKey        string             `json:"api_key"`
	Model         *models.ModelField `json:"model"`
	Endpoint      string             `json:"endpoint"`
	QueryEndpoint string             `json:"query_endpoint"`
	Priority      *int               `json:"priority"`
	IsDefault     bool               `json:"is_default"`
	IsActive      bool               `json:"is_active"`
	Settings      string             `json:"settings"`
}

type TestConnectionRequest struct {
	BaseURL  string            `json:"base_url" binding:"required,url"`
	APIKey   string            `json:"api_key" binding:"required"`
	Model    models.ModelField `json:"model" binding:"required"`
	Provider string            `json:"provider"`
	Endpoint string            `json:"endpoint"`
}

func defaultAIEndpoints(serviceType string, provider string) (string, string) {
	switch provider {
	case "gemini", "google":
		if serviceType == "text" || serviceType == "image" {
			return "/v1beta/models/{model}:generateContent", ""
		}
	case "openai", "openclaw":
		if serviceType == "text" {
			return "/chat/completions", ""
		}
		if serviceType == "image" {
			return "/images/generations", ""
		}
		if serviceType == "video" {
			return "/videos", "/videos/{taskId}"
		}
	case "chatfire":
		if serviceType == "text" {
			return "/chat/completions", ""
		}
		if serviceType == "image" {
			return "/images/generations", ""
		}
		if serviceType == "video" {
			return "/video/generations", "/video/task/{taskId}"
		}
	case "doubao", "volcengine", "volces":
		if serviceType == "video" {
			return "/contents/generations/tasks", "/contents/generations/tasks/{taskId}"
		}
	case "minimax":
		if serviceType == "video" {
			return "/video_generation", "/query/video_generation?task_id={taskId}"
		}
	case "custom_relay":
		if serviceType == "text" {
			return "/chat/completions", ""
		}
		if serviceType == "image" {
			return "/images/generations", ""
		}
		if serviceType == "video" {
			return "/video/generations", "/video/task/{taskId}"
		}
	}

	if serviceType == "image" {
		return "/images/generations", ""
	}
	if serviceType == "video" {
		return "/videos", "/videos/{taskId}"
	}
	return "/chat/completions", ""
}

func (s *AIService) CreateConfig(req *CreateAIConfigRequest) (*models.AIServiceConfig, error) {
	endpoint := req.Endpoint
	queryEndpoint := req.QueryEndpoint

	if endpoint == "" {
		defaultEndpoint, defaultQueryEndpoint := defaultAIEndpoints(req.ServiceType, req.Provider)
		endpoint = defaultEndpoint
		if queryEndpoint == "" {
			queryEndpoint = defaultQueryEndpoint
		}
	}

	config := &models.AIServiceConfig{
		ServiceType:   req.ServiceType,
		Name:          req.Name,
		Provider:      req.Provider,
		BaseURL:       req.BaseURL,
		APIKey:        req.APIKey,
		Model:         req.Model,
		Endpoint:      endpoint,
		QueryEndpoint: queryEndpoint,
		Priority:      req.Priority,
		IsDefault:     req.IsDefault,
		IsActive:      true,
		Settings:      req.Settings,
	}

	if err := s.db.Create(config).Error; err != nil {
		s.log.Errorw("Failed to create AI config", "error", err)
		return nil, err
	}

	s.log.Infow("AI config created", "config_id", config.ID, "provider", req.Provider, "endpoint", endpoint)
	return config, nil
}

func (s *AIService) GetConfig(configID uint) (*models.AIServiceConfig, error) {
	var config models.AIServiceConfig
	err := s.db.Where("id = ? ", configID).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("config not found")
		}
		return nil, err
	}
	return &config, nil
}

func (s *AIService) ListConfigs(serviceType string) ([]models.AIServiceConfig, error) {
	var configs []models.AIServiceConfig
	query := s.db

	if serviceType != "" {
		query = query.Where("service_type = ?", serviceType)
	}

	err := query.Order("priority DESC, created_at DESC").Find(&configs).Error
	if err != nil {
		s.log.Errorw("Failed to list AI configs", "error", err)
		return nil, err
	}

	return configs, nil
}

func (s *AIService) UpdateConfig(configID uint, req *UpdateAIConfigRequest) (*models.AIServiceConfig, error) {
	var config models.AIServiceConfig
	if err := s.db.Where("id = ? ", configID).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("config not found")
		}
		return nil, err
	}

	tx := s.db.Begin()

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Provider != "" {
		updates["provider"] = req.Provider
	}
	if req.BaseURL != "" {
		updates["base_url"] = req.BaseURL
	}
	if req.APIKey != "" {
		updates["api_key"] = req.APIKey
	}
	if req.Model != nil && len(*req.Model) > 0 {
		updates["model"] = *req.Model
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}

	if req.Provider != "" && req.Endpoint == "" {
		endpoint, queryEndpoint := defaultAIEndpoints(config.ServiceType, req.Provider)
		updates["endpoint"] = endpoint
		if req.QueryEndpoint == "" {
			updates["query_endpoint"] = queryEndpoint
		}
	} else if req.Endpoint != "" {
		updates["endpoint"] = req.Endpoint
	}

	updates["query_endpoint"] = req.QueryEndpoint
	if req.Settings != "" {
		updates["settings"] = req.Settings
	}
	updates["is_default"] = req.IsDefault
	updates["is_active"] = req.IsActive

	if err := tx.Model(&config).Updates(updates).Error; err != nil {
		tx.Rollback()
		s.log.Errorw("Failed to update AI config", "error", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	s.log.Infow("AI config updated", "config_id", configID)
	return &config, nil
}

func (s *AIService) DeleteConfig(configID uint) error {
	result := s.db.Where("id = ? ", configID).Delete(&models.AIServiceConfig{})

	if result.Error != nil {
		s.log.Errorw("Failed to delete AI config", "error", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("config not found")
	}

	s.log.Infow("AI config deleted", "config_id", configID)
	return nil
}

func (s *AIService) TestConnection(req *TestConnectionRequest) error {
	s.log.Infow("TestConnection called", "baseURL", req.BaseURL, "provider", req.Provider, "endpoint", req.Endpoint, "modelCount", len(req.Model))

	model := ""
	if len(req.Model) > 0 {
		model = req.Model[0]
	}
	s.log.Infow("Using model for test", "model", model, "provider", req.Provider)

	var client ai.AIClient
	var endpoint string

	switch req.Provider {
	case "gemini", "google":
		s.log.Infow("Using Gemini client", "baseURL", req.BaseURL)
		endpoint = "/v1beta/models/{model}:generateContent"
		client = ai.NewGeminiClient(req.BaseURL, req.APIKey, model, endpoint)
	case "openai", "chatfire", "openclaw", "custom_relay":
		s.log.Infow("Using OpenAI-compatible client", "baseURL", req.BaseURL, "provider", req.Provider)
		endpoint = req.Endpoint
		if endpoint == "" {
			endpoint = "/chat/completions"
		}
		client = ai.NewOpenAIClient(req.BaseURL, req.APIKey, model, endpoint)
	default:
		s.log.Infow("Using default OpenAI-compatible client", "baseURL", req.BaseURL)
		endpoint = req.Endpoint
		if endpoint == "" {
			endpoint = "/chat/completions"
		}
		client = ai.NewOpenAIClient(req.BaseURL, req.APIKey, model, endpoint)
	}

	s.log.Infow("Calling TestConnection on client", "endpoint", endpoint)
	err := client.TestConnection()
	if err != nil {
		s.log.Errorw("TestConnection failed", "error", err)
	} else {
		s.log.Infow("TestConnection succeeded")
	}
	return err
}

func (s *AIService) GetDefaultConfig(serviceType string) (*models.AIServiceConfig, error) {
	var config models.AIServiceConfig
	err := s.db.Where("service_type = ? AND is_active = ?", serviceType, true).
		Order("priority DESC, created_at DESC").
		First(&config).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active config found")
		}
		return nil, err
	}

	return &config, nil
}

func normalizeAIProvider(provider string) string {
	return strings.ToLower(strings.TrimSpace(provider))
}

func configSupportsModel(config *models.AIServiceConfig, modelName string) bool {
	if config == nil || modelName == "" {
		return false
	}

	for _, model := range config.Model {
		if model == modelName {
			return true
		}
	}

	return false
}

func (s *AIService) GetPreferredConfig(serviceType string, provider string, modelName string) (*models.AIServiceConfig, error) {
	var configs []models.AIServiceConfig
	err := s.db.Where("service_type = ? AND is_active = ?", serviceType, true).
		Order("priority DESC, created_at DESC").
		Find(&configs).Error
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return nil, errors.New("no active config found")
	}

	normalizedProvider := normalizeAIProvider(provider)

	if normalizedProvider != "" && modelName != "" {
		for _, config := range configs {
			if normalizeAIProvider(config.Provider) == normalizedProvider && configSupportsModel(&config, modelName) {
				return &config, nil
			}
		}
	}

	if normalizedProvider != "" {
		for _, config := range configs {
			if normalizeAIProvider(config.Provider) == normalizedProvider {
				return &config, nil
			}
		}
	}

	if modelName != "" {
		for _, config := range configs {
			if configSupportsModel(&config, modelName) {
				return &config, nil
			}
		}
	}

	return &configs[0], nil
}

func (s *AIService) GetConfigForModel(serviceType string, modelName string) (*models.AIServiceConfig, error) {
	var configs []models.AIServiceConfig
	err := s.db.Where("service_type = ? AND is_active = ?", serviceType, true).
		Order("priority DESC, created_at DESC").
		Find(&configs).Error

	if err != nil {
		return nil, err
	}

	for _, config := range configs {
		for _, model := range config.Model {
			if model == modelName {
				return &config, nil
			}
		}
	}

	return nil, errors.New("no active config found for model: " + modelName)
}

func (s *AIService) GetAIClient(serviceType string) (ai.AIClient, error) {
	config, err := s.GetDefaultConfig(serviceType)
	if err != nil {
		return nil, err
	}

	model := ""
	if len(config.Model) > 0 {
		model = config.Model[0]
	}

	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint, _ = defaultAIEndpoints(serviceType, config.Provider)
	}

	switch config.Provider {
	case "gemini", "google":
		return ai.NewGeminiClient(config.BaseURL, config.APIKey, model, endpoint), nil
	default:
		return ai.NewOpenAIClient(config.BaseURL, config.APIKey, model, endpoint), nil
	}
}

func (s *AIService) GetAIClientForModel(serviceType string, modelName string) (ai.AIClient, error) {
	config, err := s.GetConfigForModel(serviceType, modelName)
	if err != nil {
		return nil, err
	}

	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint, _ = defaultAIEndpoints(serviceType, config.Provider)
	}

	switch config.Provider {
	case "gemini", "google":
		return ai.NewGeminiClient(config.BaseURL, config.APIKey, modelName, endpoint), nil
	default:
		return ai.NewOpenAIClient(config.BaseURL, config.APIKey, modelName, endpoint), nil
	}
}

func (s *AIService) GenerateText(prompt string, systemPrompt string, options ...func(*ai.ChatCompletionRequest)) (string, error) {
	client, err := s.GetAIClient("text")
	if err != nil {
		return "", fmt.Errorf("failed to get AI client: %w", err)
	}

	return client.GenerateText(prompt, systemPrompt, options...)
}

func (s *AIService) GenerateImage(prompt string, size string, n int) ([]string, error) {
	client, err := s.GetAIClient("image")
	if err != nil {
		return nil, fmt.Errorf("failed to get AI client for image: %w", err)
	}

	return client.GenerateImage(prompt, size, n)
}
