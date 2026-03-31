package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	models "github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/pkg/ai"
	"github.com/drama-generator/backend/pkg/config"
	"github.com/drama-generator/backend/pkg/logger"
	"github.com/drama-generator/backend/pkg/utils"
	"gorm.io/gorm"
)

var minorAgePattern = regexp.MustCompile(`([0-9]{1,2})\s*岁`)

type CharacterLibraryService struct {
	db          *gorm.DB
	log         *logger.Logger
	config      *config.Config
	aiService   *AIService
	taskService *TaskService
	promptI18n  *PromptI18n
}

type characterViewConfig struct {
	frameType string
	viewHint  string
}

func containsMinorDescriptor(prompt string) bool {
	matches := minorAgePattern.FindAllStringSubmatch(prompt, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		if age := strings.TrimSpace(match[1]); age != "" {
			switch age {
			case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16":
				return true
			}
		}
	}

	minorHints := []string{
		"小女孩", "小男孩", "女孩", "男孩", "孩童", "儿童", "小朋友", "少女", "少年",
	}
	for _, hint := range minorHints {
		if strings.Contains(prompt, hint) {
			return true
		}
	}

	return false
}

func sanitizeMinorCharacterPrompt(prompt string) string {
	if !containsMinorDescriptor(prompt) {
		return prompt
	}

	replacer := strings.NewReplacer(
		"长期缺乏照顾", "安静内敛",
		"瘦弱单薄", "身形轻盈",
		"瘦小", "娇小",
		"不健康的潮红", "自然红润",
		"火辣辣的潮红", "自然红润",
		"细密的汗珠", "清爽自然的面部细节",
		"嘴唇略显干裂", "嘴唇线条柔和",
		"委屈", "安静",
		"泪水", "清澈目光",
		"微红令人心碎", "自然柔和",
		"令人心碎的凄凉感", "温柔安静的气质",
		"凄凉感", "安静气质",
		"破碎感", "柔和气质",
		"脆弱", "温柔",
	)

	safePrompt := replacer.Replace(prompt)
	safePrompt += "，整体形象健康、自然、干净，儿童形象温和友好。"
	return safePrompt
}

func getCharacterViewConfigs() []characterViewConfig {
	return []characterViewConfig{
		{frameType: "front", viewHint: "front-facing full body standing pose, face fully visible, both eyes visible, shoulders and toes facing the camera, not side view, not back view"},
		{frameType: "side", viewHint: "strict side profile full body standing pose, only one side of the face visible, body turned 90 degrees, not front view, not back view"},
		{frameType: "back", viewHint: "back-facing full body standing pose, back of head and outfit visible, no facial features visible, not front view, not side view"},
	}
}

func getCharacterViewConfig(frameType string) (characterViewConfig, bool) {
	for _, config := range getCharacterViewConfigs() {
		if config.frameType == frameType {
			return config, true
		}
	}
	return characterViewConfig{}, false
}

func (s *CharacterLibraryService) buildCharacterImagePrompt(character models.Character, drama models.Drama, style string) string {
	prompt := character.Name
	if character.Appearance != nil && *character.Appearance != "" {
		prompt = *character.Appearance
	} else if character.Description != nil && *character.Description != "" {
		prompt = *character.Description
	}

	styleHints := make([]string, 0, 2)
	if drama.Style != "" && drama.Style != "realistic" {
		styleHints = append(styleHints, drama.Style)
	}
	if strings.TrimSpace(style) != "" {
		styleHints = append(styleHints, strings.TrimSpace(style))
	}
	if len(styleHints) > 0 {
		prompt += ", " + strings.Join(styleHints, ", ")
	}

	return sanitizeMinorCharacterPrompt(prompt)
}

func (s *CharacterLibraryService) generateCharacterViewImage(
	character *models.Character,
	drama *models.Drama,
	frameType string,
	modelName string,
	style string,
	imageService *ImageGenerationService,
) (*models.ImageGeneration, error) {
	viewConfig, ok := getCharacterViewConfig(frameType)
	if !ok {
		return nil, fmt.Errorf("unsupported frame type: %s", frameType)
	}

	prompt := s.buildCharacterImagePrompt(*character, *drama, style)
	negativePrompt := "complex background, scenery, furniture, props, crowd, multiple characters, text, watermark, logo, border, frame"
	viewPrompt := fmt.Sprintf(
		"%s, same character design, same clothing, same hairstyle, same accessories, %s, pure white seamless studio background, isolated single character, no environment, no props, no extra objects, clean soft studio lighting, professional character turnaround sheet, highly consistent identity",
		prompt,
		viewConfig.viewHint,
	)

	dramaIDStr := fmt.Sprintf("%d", character.DramaID)
	imageType := "character"
	req := &GenerateImageRequest{
		DramaID:        dramaIDStr,
		CharacterID:    &character.ID,
		ImageType:      imageType,
		FrameType:      &viewConfig.frameType,
		Prompt:         viewPrompt,
		NegativePrompt: &negativePrompt,
		Provider:       "openai",
		Model:          modelName,
		Size:           "2560x1440",
		Quality:        "standard",
	}

	return imageService.GenerateImage(req)
}

func (s *CharacterLibraryService) generateCharacterTurnaroundSheet(
	character *models.Character,
	drama *models.Drama,
	modelName string,
	style string,
	imageService *ImageGenerationService,
) (*models.ImageGeneration, error) {
	prompt := s.buildCharacterImagePrompt(*character, *drama, style)
	negativePrompt := "complex background, scenery, furniture, props, crowd, multiple characters, text, watermark, logo, border, frame, cropped body, close-up portrait, split panel border"
	frameType := "sheet"
	viewPrompt := fmt.Sprintf(
		"%s, one single turnaround sheet image, show the same character three times in one canvas: left strict side profile full body, center front-facing full body, right back-facing full body, identical clothing hairstyle accessories and body proportions across all three poses, pure white seamless studio background, isolated single character presentation, no environment, no props, no extra objects, clean soft studio lighting, fashion lookbook layout, full body visible from head to toe, evenly spaced composition",
		prompt,
	)

	dramaIDStr := fmt.Sprintf("%d", character.DramaID)
	req := &GenerateImageRequest{
		DramaID:        dramaIDStr,
		CharacterID:    &character.ID,
		ImageType:      "character",
		FrameType:      &frameType,
		Prompt:         viewPrompt,
		NegativePrompt: &negativePrompt,
		Provider:       "openai",
		Model:          modelName,
		// 提高分辨率，满足上游最小像素限制
		Size:           "2560x1600",
		Quality:        "standard",
	}

	return imageService.GenerateImage(req)
}

func NewCharacterLibraryService(db *gorm.DB, log *logger.Logger, cfg *config.Config) *CharacterLibraryService {
	return &CharacterLibraryService{
		db:          db,
		log:         log,
		config:      cfg,
		aiService:   NewAIService(db, log),
		taskService: NewTaskService(db, log),
		promptI18n:  NewPromptI18n(cfg),
	}
}

type CreateLibraryItemRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Category    *string `json:"category"`
	ImageURL    string  `json:"image_url" binding:"required"`
	LocalPath   *string `json:"local_path"`
	Description *string `json:"description"`
	Tags        *string `json:"tags"`
	SourceType  string  `json:"source_type"`
}

type CharacterLibraryQuery struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=20"`
	Category   string `form:"category"`
	SourceType string `form:"source_type"`
	Keyword    string `form:"keyword"`
}

// ListLibraryItems 获取用户角色库列表
func (s *CharacterLibraryService) ListLibraryItems(query *CharacterLibraryQuery) ([]models.CharacterLibrary, int64, error) {
	var items []models.CharacterLibrary
	var total int64

	db := s.db.Model(&models.CharacterLibrary{})

	// 筛选条件
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}

	if query.SourceType != "" {
		db = db.Where("source_type = ?", query.SourceType)
	}

	if query.Keyword != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		s.log.Errorw("Failed to count character library", "error", err)
		return nil, 0, err
	}

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	err := db.Order("created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&items).Error

	if err != nil {
		s.log.Errorw("Failed to list character library", "error", err)
		return nil, 0, err
	}

	return items, total, nil
}

// CreateLibraryItem 添加到角色库
func (s *CharacterLibraryService) CreateLibraryItem(req *CreateLibraryItemRequest) (*models.CharacterLibrary, error) {
	sourceType := req.SourceType
	if sourceType == "" {
		sourceType = "generated"
	}

	item := &models.CharacterLibrary{
		Name:        req.Name,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		LocalPath:   req.LocalPath,
		Description: req.Description,
		Tags:        req.Tags,
		SourceType:  sourceType,
	}

	if err := s.db.Create(item).Error; err != nil {
		s.log.Errorw("Failed to create library item", "error", err)
		return nil, err
	}

	s.log.Infow("Library item created", "item_id", item.ID)
	return item, nil
}

// GetLibraryItem 获取角色库项
func (s *CharacterLibraryService) GetLibraryItem(itemID string) (*models.CharacterLibrary, error) {
	var item models.CharacterLibrary
	err := s.db.Where("id = ? ", itemID).First(&item).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("library item not found")
		}
		s.log.Errorw("Failed to get library item", "error", err)
		return nil, err
	}

	return &item, nil
}

// DeleteLibraryItem 删除角色库项
func (s *CharacterLibraryService) DeleteLibraryItem(itemID string) error {
	result := s.db.Where("id = ? ", itemID).Delete(&models.CharacterLibrary{})

	if result.Error != nil {
		s.log.Errorw("Failed to delete library item", "error", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("library item not found")
	}

	s.log.Infow("Library item deleted", "item_id", itemID)
	return nil
}

// ApplyLibraryItemToCharacter 将角色库形象应用到角色
func (s *CharacterLibraryService) ApplyLibraryItemToCharacter(characterID string, libraryItemID string) error {
	// 验证角色库项存在且属于该用户
	var libraryItem models.CharacterLibrary
	if err := s.db.Where("id = ? ", libraryItemID).First(&libraryItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("library item not found")
		}
		return err
	}

	// 查找角色
	var character models.Character
	if err := s.db.Where("id = ?", characterID).First(&character).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("character not found")
		}
		return err
	}

	// 查询Drama验证权限
	var drama models.Drama
	if err := s.db.Where("id = ? ", character.DramaID).First(&drama).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("unauthorized")
		}
		return err
	}

	// 更新角色的 local_path 和 image_url
	updates := map[string]interface{}{}
	if libraryItem.LocalPath != nil && *libraryItem.LocalPath != "" {
		updates["local_path"] = libraryItem.LocalPath
	}
	if libraryItem.ImageURL != "" {
		updates["image_url"] = libraryItem.ImageURL
	}
	if len(updates) > 0 {
		if err := s.db.Model(&character).Updates(updates).Error; err != nil {
			s.log.Errorw("Failed to update character image", "error", err)
			return err
		}
	}

	s.log.Infow("Library item applied to character", "character_id", characterID, "library_item_id", libraryItemID)
	return nil
}

// UploadCharacterImage 上传角色图片
func (s *CharacterLibraryService) UploadCharacterImage(characterID string, imageURL string) error {
	// 查找角色
	var character models.Character
	if err := s.db.Where("id = ?", characterID).First(&character).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("character not found")
		}
		return err
	}

	// 查询Drama验证权限
	var drama models.Drama
	if err := s.db.Where("id = ? ", character.DramaID).First(&drama).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("unauthorized")
		}
		return err
	}

	// 更新图片URL
	if err := s.db.Model(&character).Update("image_url", imageURL).Error; err != nil {
		s.log.Errorw("Failed to update character image", "error", err)
		return err
	}

	s.log.Infow("Character image uploaded", "character_id", characterID)
	return nil
}

// AddCharacterToLibrary 将角色添加到角色库
func (s *CharacterLibraryService) AddCharacterToLibrary(characterID string, category *string) (*models.CharacterLibrary, error) {
	// 查找角色
	var character models.Character
	if err := s.db.Where("id = ?", characterID).First(&character).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("character not found")
		}
		return nil, err
	}

	// 查询Drama验证权限
	var drama models.Drama
	if err := s.db.Where("id = ? ", character.DramaID).First(&drama).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unauthorized")
		}
		return nil, err
	}

	// 检查是否有图片
	if character.ImageURL == nil || *character.ImageURL == "" {
		return nil, fmt.Errorf("角色还没有形象图片")
	}

	// 创建角色库项
	charLibrary := &models.CharacterLibrary{
		Name:        character.Name,
		ImageURL:    *character.ImageURL,
		LocalPath:   character.LocalPath,
		Description: character.Description,
		SourceType:  "character",
	}

	if err := s.db.Create(charLibrary).Error; err != nil {
		s.log.Errorw("Failed to add character to library", "error", err)
		return nil, err
	}

	s.log.Infow("Character added to library", "character_id", characterID, "library_item_id", charLibrary.ID)
	return charLibrary, nil
}

// DeleteCharacter 删除单个角色
func (s *CharacterLibraryService) DeleteCharacter(characterID uint) error {
	// 查找角色
	var character models.Character
	if err := s.db.Where("id = ?", characterID).First(&character).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("character not found")
		}
		return err
	}

	// 验证权限：检查角色所属的drama是否属于当前用户
	var drama models.Drama
	if err := s.db.Where("id = ? ", character.DramaID).First(&drama).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("unauthorized")
		}
		return err
	}

	// 删除角色
	if err := s.db.Delete(&character).Error; err != nil {
		s.log.Errorw("Failed to delete character", "error", err, "id", characterID)
		return err
	}

	s.log.Infow("Character deleted", "id", characterID)
	return nil
}

// GenerateCharacterImage AI生成角色形象
func (s *CharacterLibraryService) GenerateCharacterImage(characterID string, imageService *ImageGenerationService, modelName string, style string) (*models.ImageGeneration, error) {
	// 查找角色
	var character models.Character
	if err := s.db.Where("id = ?", characterID).First(&character).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("character not found")
		}
		return nil, err
	}

	// 查询Drama验证权限
	var drama models.Drama
	if err := s.db.Where("id = ? ", character.DramaID).First(&drama).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unauthorized")
		}
		return nil, err
	}

	emptyReferenceImages, err := json.Marshal([]string{})
	if err == nil {
		if updateErr := s.db.Model(&models.Character{}).
			Where("id = ?", character.ID).
			Update("reference_images", emptyReferenceImages).Error; updateErr != nil {
			s.log.Warnw("Failed to reset character reference image slots", "character_id", characterID, "error", updateErr)
		}
	}

	imageGen, err := s.generateCharacterTurnaroundSheet(&character, &drama, modelName, style, imageService)
	if err != nil {
		s.log.Errorw("Failed to generate character turnaround sheet", "error", err, "character_id", characterID)
		return nil, fmt.Errorf("图片生成失败: %w", err)
	}

	s.log.Infow("Character turnaround sheet generation started", "character_id", characterID, "image_gen_id", imageGen.ID)
	return imageGen, nil
}

func (s *CharacterLibraryService) GenerateSingleCharacterImage(characterID string, frameType string, imageService *ImageGenerationService, modelName string, style string) (*models.ImageGeneration, error) {
	var character models.Character
	if err := s.db.Where("id = ?", characterID).First(&character).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("character not found")
		}
		return nil, err
	}

	var drama models.Drama
	if err := s.db.Where("id = ? ", character.DramaID).First(&drama).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unauthorized")
		}
		return nil, err
	}

	imageGen, err := s.generateCharacterViewImage(&character, &drama, frameType, modelName, style, imageService)
	if err != nil {
		return nil, fmt.Errorf("图片生成失败: %w", err)
	}

	s.log.Infow("Character single-view image generation started", "character_id", characterID, "image_gen_id", imageGen.ID, "frame_type", frameType)
	return imageGen, nil
}

// waitAndUpdateCharacterImage 后台异步等待图片生成完成并更新角色image_url
func (s *CharacterLibraryService) waitAndUpdateCharacterImage(characterID uint, imageGenID uint) {
	maxAttempts := 60
	pollInterval := 5 * time.Second

	for i := 0; i < maxAttempts; i++ {
		time.Sleep(pollInterval)

		// 查询图片生成状态
		var imageGen models.ImageGeneration
		if err := s.db.First(&imageGen, imageGenID).Error; err != nil {
			s.log.Errorw("Failed to query image generation status", "error", err, "image_gen_id", imageGenID)
			continue
		}

		// 检查是否完成
		if imageGen.Status == models.ImageStatusCompleted && imageGen.ImageURL != nil && *imageGen.ImageURL != "" {
			// 更新角色的image_url
			if err := s.db.Model(&models.Character{}).Where("id = ?", characterID).Update("image_url", *imageGen.ImageURL).Error; err != nil {
				s.log.Errorw("Failed to update character image_url", "error", err, "character_id", characterID)
				return
			}
			s.log.Infow("Character image updated successfully", "character_id", characterID, "image_url", *imageGen.ImageURL)
			return
		}

		// 检查是否失败
		if imageGen.Status == models.ImageStatusFailed {
			s.log.Errorw("Character image generation failed", "character_id", characterID, "image_gen_id", imageGenID, "error", imageGen.ErrorMsg)
			return
		}
	}

	s.log.Warnw("Character image generation timeout", "character_id", characterID, "image_gen_id", imageGenID)
}

type UpdateCharacterRequest struct {
	Name        *string `json:"name"`
	Role        *string `json:"role"`
	Appearance  *string `json:"appearance"`
	Personality *string `json:"personality"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image_url"`
	LocalPath   *string `json:"local_path"`
}

// UpdateCharacter 更新角色信息
func (s *CharacterLibraryService) UpdateCharacter(characterID string, req *UpdateCharacterRequest) error {
	// 查找角色
	var character models.Character
	if err := s.db.Where("id = ?", characterID).First(&character).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("character not found")
		}
		return err
	}

	// 验证权限：查询角色所属的drama是否属于该用户
	var drama models.Drama
	if err := s.db.Where("id = ? ", character.DramaID).First(&drama).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("unauthorized")
		}
		return err
	}

	// 构建更新数据
	updates := make(map[string]interface{})

	if req.Name != nil && *req.Name != "" {
		updates["name"] = *req.Name
	}
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.Appearance != nil {
		updates["appearance"] = *req.Appearance
	}
	if req.Personality != nil {
		updates["personality"] = *req.Personality
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}
	if req.LocalPath != nil {
		updates["local_path"] = *req.LocalPath
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	// 更新角色信息
	if err := s.db.Model(&character).Updates(updates).Error; err != nil {
		s.log.Errorw("Failed to update character", "error", err, "character_id", characterID)
		return err
	}

	s.log.Infow("Character updated", "character_id", characterID, "updates", updates)
	return nil
}

// BatchGenerateCharacterImages 批量生成角色图片（并发执行）
func (s *CharacterLibraryService) BatchGenerateCharacterImages(characterIDs []string, imageService *ImageGenerationService, modelName string) {
	s.log.Infow("Starting batch character image generation",
		"count", len(characterIDs),
		"model", modelName)

	// 使用 goroutine 并发生成所有角色图片
	for _, characterID := range characterIDs {
		// 为每个角色启动单独的 goroutine
		go func(charID string) {
			imageGen, err := s.GenerateCharacterImage(charID, imageService, modelName, "") // 批量生成暂不支持自定义风格，使用默认值
			if err != nil {
				s.log.Errorw("Failed to generate character image in batch",
					"character_id", charID,
					"error", err)
				return
			}

			s.log.Infow("Character image generated in batch",
				"character_id", charID,
				"image_gen_id", imageGen.ID)
		}(characterID)
	}

	s.log.Infow("Batch character image generation tasks submitted",
		"total", len(characterIDs))
}

// ExtractCharactersFromScript 从分集剧本中提取角色
func (s *CharacterLibraryService) ExtractCharactersFromScript(episodeID uint) (string, error) {
	var episode models.Episode
	if err := s.db.First(&episode, episodeID).Error; err != nil {
		return "", fmt.Errorf("episode not found")
	}

	if episode.ScriptContent == nil || *episode.ScriptContent == "" {
		return "", fmt.Errorf("剧本内容为空")
	}

	task, err := s.taskService.CreateTask("character_extraction", fmt.Sprintf("%d", episode.DramaID))
	if err != nil {
		return "", fmt.Errorf("创建任务失败: %w", err)
	}

	go s.processCharacterExtraction(task.ID, episode)

	return task.ID, nil
}

func (s *CharacterLibraryService) processCharacterExtraction(taskID string, episode models.Episode) {
	s.taskService.UpdateTaskStatus(taskID, "processing", 0, "正在分析剧本...")

	script := ""
	if episode.ScriptContent != nil {
		script = *episode.ScriptContent
	}

	// 获取 drama 的 style 信息
	var drama models.Drama
	if err := s.db.First(&drama, episode.DramaID).Error; err != nil {
		s.log.Warnw("Failed to load drama", "error", err, "drama_id", episode.DramaID)
	}

	prompt := s.promptI18n.GetCharacterExtractionPrompt(drama.Style)
	userPrompt := fmt.Sprintf("【剧本内容】\n%s", script)

	response, err := s.aiService.GenerateText(userPrompt, prompt, ai.WithMaxTokens(3000))
	if err != nil {
		s.taskService.UpdateTaskError(taskID, err)
		return
	}

	s.taskService.UpdateTaskStatus(taskID, "processing", 50, "正在整理角色数据...")

	var extractedCharacters []struct {
		Name        string `json:"name"`
		Role        string `json:"role"`
		Appearance  string `json:"appearance"`
		Personality string `json:"personality"`
		Description string `json:"description"`
	}

	if err := utils.SafeParseAIJSON(response, &extractedCharacters); err != nil {
		s.log.Errorw("Failed to parse AI response for characters", "error", err, "response", response)
		s.taskService.UpdateTaskError(taskID, fmt.Errorf("解析AI响应失败"))
		return
	}

	var savedCharacters []models.Character
	for _, charData := range extractedCharacters {
		// 检查是否已存在同名角色
		var existingCharacter models.Character
		err := s.db.Where("drama_id = ? AND name = ?", episode.DramaID, charData.Name).First(&existingCharacter).Error

		if err == nil {
			// 如果存在，只关联，不更新（或者可以选更新，这里暂不更新）
			if err := s.db.Model(&episode).Association("Characters").Append(&existingCharacter); err != nil {
				s.log.Warnw("Failed to associate existing character", "error", err)
			}
			savedCharacters = append(savedCharacters, existingCharacter)
		} else {
			// 创建新角色
			newCharacter := models.Character{
				DramaID:     episode.DramaID,
				Name:        charData.Name,
				Role:        &charData.Role,
				Appearance:  &charData.Appearance,
				Personality: &charData.Personality,
				Description: &charData.Description,
			}
			if err := s.db.Create(&newCharacter).Error; err != nil {
				s.log.Errorw("Failed to create extracted character", "error", err)
				continue
			}

			// 关联到分集
			if err := s.db.Model(&episode).Association("Characters").Append(&newCharacter); err != nil {
				s.log.Warnw("Failed to associate new character", "error", err)
			}
			savedCharacters = append(savedCharacters, newCharacter)
		}
	}

	s.taskService.UpdateTaskResult(taskID, map[string]interface{}{
		"characters": savedCharacters,
		"count":      len(savedCharacters),
	})
}
