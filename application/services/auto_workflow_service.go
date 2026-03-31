package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/infrastructure/storage"
	"github.com/drama-generator/backend/pkg/config"
	"github.com/drama-generator/backend/pkg/logger"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	workflowStagePrepareEpisodes   = "prepare_episodes"
	workflowStageCharacterExtract  = "character_extract"
	workflowStageSceneExtract      = "scene_extract"
	workflowStagePropExtract       = "prop_extract"
	workflowStageCharacterBaseline = "character_baseline_generate"
	workflowStageStoryboard        = "storyboard_generate"
	workflowStageStoryboardImage   = "storyboard_image_generate"
	workflowStageStoryboardVideo   = "storyboard_video_generate"
	workflowStageEpisodeMerge      = "episode_merge"
	workflowStageCompleted         = "completed"
)

type WorkflowEpisodeSummary struct {
	EpisodeID         uint   `json:"episode_id"`
	EpisodeNumber     int    `json:"episode_number"`
	Title             string `json:"title"`
	Status            string `json:"status"`
	StoryboardCount   int64  `json:"storyboard_count"`
	CompletedImages   int64  `json:"completed_images"`
	CompletedVideos   int64  `json:"completed_videos"`
	FailedImages      int64  `json:"failed_images"`
	FailedVideos      int64  `json:"failed_videos"`
	MergeStatus       string `json:"merge_status,omitempty"`
	FinalVideoURL     string `json:"final_video_url,omitempty"`
	NeedsManualReview bool   `json:"needs_manual_review"`
}

type workflowStageError struct {
	Stage     string
	EpisodeID *uint
	Err       error
}

func (e *workflowStageError) Error() string {
	if e == nil || e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func (e *workflowStageError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

type WorkflowStatusResponse struct {
	Run      *models.WorkflowRun      `json:"run"`
	Steps    []models.WorkflowStepRun `json:"steps"`
	Episodes []WorkflowEpisodeSummary `json:"episodes"`
}

type AutoWorkflowService struct {
	db                *gorm.DB
	log               *logger.Logger
	cfg               *config.Config
	taskService       *TaskService
	dramaService      *DramaService
	storyboardService *StoryboardService
	imageService      *ImageGenerationService
	videoService      *VideoGenerationService
	videoMergeService *VideoMergeService
	characterService  *CharacterLibraryService
	propService       *PropService
}

func NewAutoWorkflowService(db *gorm.DB, cfg *config.Config, transferService *ResourceTransferService, localStorage *storage.LocalStorage, log *logger.Logger) *AutoWorkflowService {
	aiService := NewAIService(db, log)
	taskService := NewTaskService(db, log)
	imageService := NewImageGenerationService(db, cfg, transferService, localStorage, log)

	return &AutoWorkflowService{
		db:                db,
		log:               log,
		cfg:               cfg,
		taskService:       taskService,
		dramaService:      NewDramaService(db, cfg, log),
		storyboardService: NewStoryboardService(db, cfg, log),
		imageService:      imageService,
		videoService:      NewVideoGenerationService(db, transferService, localStorage, aiService, log, NewPromptI18n(cfg)),
		videoMergeService: NewVideoMergeService(db, transferService, cfg.Storage.LocalPath, cfg.Storage.BaseURL, log),
		characterService:  NewCharacterLibraryService(db, log, cfg),
		propService:       NewPropService(db, aiService, taskService, imageService, log, cfg),
	}
}

func (s *AutoWorkflowService) StartProjectWorkflow(dramaID string) (*models.WorkflowRun, error) {
	drama, err := s.dramaService.GetDrama(dramaID)
	if err != nil {
		return nil, err
	}

	var existing models.WorkflowRun
	err = s.db.
		Where("drama_id = ? AND scope = ? AND status IN ?", drama.ID, models.WorkflowRunScopeProject, []models.WorkflowRunStatus{
			models.WorkflowRunStatusPending,
			models.WorkflowRunStatusProcessing,
		}).
		Order("created_at DESC").
		First(&existing).Error
	if err == nil {
		return &existing, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	now := time.Now()
	run := &models.WorkflowRun{
		DramaID:      drama.ID,
		Scope:        models.WorkflowRunScopeProject,
		Status:       models.WorkflowRunStatusPending,
		CurrentStage: workflowStagePrepareEpisodes,
		Progress:     0,
		StartedAt:    &now,
	}
	if err := s.db.Create(run).Error; err != nil {
		return nil, err
	}

	go s.processProjectWorkflow(run.ID)
	return run, nil
}

func (s *AutoWorkflowService) ResumeProjectWorkflow(dramaID string) (*models.WorkflowRun, error) {
	drama, err := s.dramaService.GetDrama(dramaID)
	if err != nil {
		return nil, err
	}

	var run models.WorkflowRun
	if err := s.db.
		Where("drama_id = ? AND scope = ?", drama.ID, models.WorkflowRunScopeProject).
		Order("created_at DESC").
		First(&run).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return s.StartProjectWorkflow(dramaID)
		}
		return nil, err
	}

	if run.Status == models.WorkflowRunStatusPending || run.Status == models.WorkflowRunStatusProcessing {
		return &run, nil
	}

	now := time.Now()
	if err := s.db.Model(&models.WorkflowRun{}).Where("id = ?", run.ID).Updates(map[string]interface{}{
		"status":       models.WorkflowRunStatusPending,
		"error_msg":    nil,
		"completed_at": nil,
		"updated_at":   now,
	}).Error; err != nil {
		return nil, err
	}

	run.Status = models.WorkflowRunStatusPending
	run.ErrorMsg = nil
	run.CompletedAt = nil

	go s.processProjectWorkflow(run.ID)
	return &run, nil
}

func (s *AutoWorkflowService) GetProjectWorkflowStatus(dramaID string) (*WorkflowStatusResponse, error) {
	drama, err := s.dramaService.GetDrama(dramaID)
	if err != nil {
		return nil, err
	}

	var run models.WorkflowRun
	if err := s.db.
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("drama_id = ? AND scope = ?", drama.ID, models.WorkflowRunScopeProject).
		Order("created_at DESC").
		First(&run).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &WorkflowStatusResponse{
				Run:      nil,
				Steps:    []models.WorkflowStepRun{},
				Episodes: s.buildEpisodeSummaries(drama.ID),
			}, nil
		}
		return nil, err
	}

	return &WorkflowStatusResponse{
		Run:      &run,
		Steps:    run.Steps,
		Episodes: s.buildEpisodeSummaries(drama.ID),
	}, nil
}

func (s *AutoWorkflowService) processProjectWorkflow(runID uint) {
	var run models.WorkflowRun
	if err := s.db.First(&run, runID).Error; err != nil {
		s.log.Errorw("Failed to load workflow run", "error", err, "run_id", runID)
		return
	}

	s.updateRun(runID, models.WorkflowRunStatusProcessing, workflowStagePrepareEpisodes, 2, "", nil)

	episodes, err := s.ensureEpisodes(run.DramaID)
	if err != nil {
		s.failRun(runID, workflowStagePrepareEpisodes, err, nil)
		return
	}
	s.completeStage(runID, nil, workflowStagePrepareEpisodes, "Episodes prepared", map[string]interface{}{"episodes": len(episodes)})

	if err := s.runCharacterExtractionStage(runID, run.DramaID, episodes); err != nil {
		s.failRun(runID, workflowStageCharacterExtract, err, nil)
		return
	}
	if err := s.runSceneExtractionStage(runID, run.DramaID, episodes); err != nil {
		s.failRun(runID, workflowStageSceneExtract, err, nil)
		return
	}
	if err := s.runPropExtractionStage(runID, run.DramaID, episodes); err != nil {
		s.failRun(runID, workflowStagePropExtract, err, nil)
		return
	}
	if err := s.runCharacterBaselineStage(runID, run.DramaID); err != nil {
		s.failRun(runID, workflowStageCharacterBaseline, err, nil)
		return
	}

	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].EpisodeNum < episodes[j].EpisodeNum
	})

	totalEpisodes := len(episodes)
	for index := range episodes {
		if err := s.runEpisodeStages(runID, &episodes[index], index, totalEpisodes); err != nil {
			var stageErr *workflowStageError
			if errors.As(err, &stageErr) {
				s.failRun(runID, stageErr.Stage, stageErr.Err, stageErr.EpisodeID)
			} else {
				s.failRun(runID, s.stageWithEpisodeNumber(workflowStageStoryboard, episodes[index].EpisodeNum), err, &episodes[index].ID)
			}
			return
		}
	}

	resultJSON, _ := json.Marshal(map[string]interface{}{
		"episodes": s.buildEpisodeSummaries(run.DramaID),
	})
	now := time.Now()
	s.db.Model(&models.WorkflowRun{}).Where("id = ?", runID).Updates(map[string]interface{}{
		"status":        models.WorkflowRunStatusCompleted,
		"current_stage": workflowStageCompleted,
		"progress":      100,
		"result_json":   datatypes.JSON(resultJSON),
		"completed_at":  &now,
		"updated_at":    now,
	})
}

func (s *AutoWorkflowService) runExtractionStages(runID uint, dramaID uint, episodes []models.Episode) error {
	if len(episodes) == 0 {
		return fmt.Errorf("no episodes available for workflow")
	}
	if err := s.runCharacterExtractionStage(runID, dramaID, episodes); err != nil {
		return err
	}
	if err := s.runSceneExtractionStage(runID, dramaID, episodes); err != nil {
		return err
	}
	if err := s.runPropExtractionStage(runID, dramaID, episodes); err != nil {
		return err
	}
	return nil
}

func (s *AutoWorkflowService) runCharacterExtractionStage(runID uint, dramaID uint, episodes []models.Episode) error {
	var characterCount int64
	if err := s.db.Model(&models.Character{}).Where("drama_id = ?", dramaID).Count(&characterCount).Error; err != nil {
		return err
	}
	if characterCount > 0 {
		s.completeStage(runID, nil, workflowStageCharacterExtract, "Characters already exist, skipping extraction", map[string]interface{}{"skipped": true, "count": characterCount})
		s.updateRun(runID, models.WorkflowRunStatusProcessing, workflowStageSceneExtract, 12, "", nil)
		return nil
	}

	s.updateRun(runID, models.WorkflowRunStatusProcessing, workflowStageCharacterExtract, 8, "", nil)
	taskIDs := make([]string, 0, len(episodes))
	for _, episode := range episodes {
		s.updateStep(runID, &episode.ID, workflowStageCharacterExtract, models.WorkflowRunStatusProcessing, 10, fmt.Sprintf("Extracting characters for episode %d", episode.EpisodeNum), nil, nil)
		taskID, err := s.characterService.ExtractCharactersFromScript(episode.ID)
		if err != nil {
			return fmt.Errorf("episode %d character extraction failed to start: %w", episode.EpisodeNum, err)
		}
		taskIDs = append(taskIDs, taskID)
	}

	for _, taskID := range taskIDs {
		task, err := s.waitForAsyncTask(taskID, 10*time.Minute)
		if err != nil {
			return err
		}
		if task.Status != "completed" {
			return fmt.Errorf("character extraction failed: %s", task.Error)
		}
	}

	for _, episode := range episodes {
		s.completeStage(runID, &episode.ID, workflowStageCharacterExtract, fmt.Sprintf("Episode %d character extraction completed", episode.EpisodeNum), nil)
	}
	s.updateRun(runID, models.WorkflowRunStatusProcessing, workflowStageSceneExtract, 18, "", nil)
	return nil
}

func (s *AutoWorkflowService) runSceneExtractionStage(runID uint, dramaID uint, episodes []models.Episode) error {
	var sceneCount int64
	if err := s.db.Model(&models.Scene{}).Where("drama_id = ?", dramaID).Count(&sceneCount).Error; err != nil {
		return err
	}
	if sceneCount > 0 {
		s.completeStage(runID, nil, workflowStageSceneExtract, "Scenes already exist, skipping extraction", map[string]interface{}{"skipped": true, "count": sceneCount})
		s.updateRun(runID, models.WorkflowRunStatusProcessing, workflowStagePropExtract, 24, "", nil)
		return nil
	}

	s.updateRun(runID, models.WorkflowRunStatusProcessing, workflowStageSceneExtract, 18, "", nil)
	taskIDs := make([]string, 0, len(episodes))
	for _, episode := range episodes {
		s.updateStep(runID, &episode.ID, workflowStageSceneExtract, models.WorkflowRunStatusProcessing, 10, fmt.Sprintf("Extracting scenes for episode %d", episode.EpisodeNum), nil, nil)
		taskID, err := s.imageService.ExtractBackgroundsForEpisode(fmt.Sprintf("%d", episode.ID), "", "")
		if err != nil {
			return fmt.Errorf("episode %d scene extraction failed to start: %w", episode.EpisodeNum, err)
		}
		taskIDs = append(taskIDs, taskID)
	}

	for _, taskID := range taskIDs {
		task, err := s.waitForAsyncTask(taskID, 10*time.Minute)
		if err != nil {
			return err
		}
		if task.Status != "completed" {
			return fmt.Errorf("scene extraction failed: %s", task.Error)
		}
	}

	for _, episode := range episodes {
		s.completeStage(runID, &episode.ID, workflowStageSceneExtract, fmt.Sprintf("Episode %d scene extraction completed", episode.EpisodeNum), nil)
	}
	s.updateRun(runID, models.WorkflowRunStatusProcessing, workflowStagePropExtract, 28, "", nil)
	return nil
}

func (s *AutoWorkflowService) runPropExtractionStage(runID uint, dramaID uint, episodes []models.Episode) error {
	var propCount int64
	if err := s.db.Model(&models.Prop{}).Where("drama_id = ?", dramaID).Count(&propCount).Error; err != nil {
		return err
	}
	if propCount > 0 {
		s.completeStage(runID, nil, workflowStagePropExtract, "Props already exist, skipping extraction", map[string]interface{}{"skipped": true, "count": propCount})
		return nil
	}

	s.updateRun(runID, models.WorkflowRunStatusProcessing, workflowStagePropExtract, 28, "", nil)
	taskIDs := make([]string, 0, len(episodes))
	for _, episode := range episodes {
		s.updateStep(runID, &episode.ID, workflowStagePropExtract, models.WorkflowRunStatusProcessing, 10, fmt.Sprintf("Extracting props for episode %d", episode.EpisodeNum), nil, nil)
		taskID, err := s.propService.ExtractPropsFromScript(episode.ID)
		if err != nil {
			return fmt.Errorf("episode %d prop extraction failed to start: %w", episode.EpisodeNum, err)
		}
		taskIDs = append(taskIDs, taskID)
	}

	for _, taskID := range taskIDs {
		task, err := s.waitForAsyncTask(taskID, 10*time.Minute)
		if err != nil {
			return err
		}
		if task.Status != "completed" {
			return fmt.Errorf("prop extraction failed: %s", task.Error)
		}
	}

	for _, episode := range episodes {
		s.completeStage(runID, &episode.ID, workflowStagePropExtract, fmt.Sprintf("Episode %d prop extraction completed", episode.EpisodeNum), nil)
	}
	return nil
}

func (s *AutoWorkflowService) runCharacterBaselineStage(runID uint, dramaID uint) error {
	s.updateRun(runID, models.WorkflowRunStatusProcessing, workflowStageCharacterBaseline, 34, "", nil)
	s.updateStep(runID, nil, workflowStageCharacterBaseline, models.WorkflowRunStatusProcessing, 5, "Checking character baseline images", nil, nil)

	var characters []models.Character
	if err := s.db.Where("drama_id = ?", dramaID).Find(&characters).Error; err != nil {
		return err
	}
	if len(characters) == 0 {
		return fmt.Errorf("no characters found after extraction")
	}

	pendingIDs := make([]string, 0)
	for _, character := range characters {
		if s.characterHasBaseline(&character) {
			continue
		}
		pendingIDs = append(pendingIDs, fmt.Sprintf("%d", character.ID))
	}

	if len(pendingIDs) == 0 {
		s.completeStage(runID, nil, workflowStageCharacterBaseline, "Character baselines already exist, skipping generation", map[string]interface{}{"skipped": true})
		return nil
	}

	for _, characterID := range pendingIDs {
		if _, err := s.characterService.GenerateCharacterImage(characterID, s.imageService, "", ""); err != nil {
			return fmt.Errorf("failed to start character baseline generation for %s: %w", characterID, err)
		}
	}

	if err := s.waitForCharacterBaselines(dramaID, 20*time.Minute); err != nil {
		return err
	}

	s.completeStage(runID, nil, workflowStageCharacterBaseline, "Character baselines generated", map[string]interface{}{"generated": len(pendingIDs)})
	return nil
}

func (s *AutoWorkflowService) runEpisodeStages(runID uint, episode *models.Episode, index int, totalEpisodes int) error {
	baseProgress := 38
	perEpisode := 62 / maxInt(totalEpisodes, 1)
	stageBase := baseProgress + index*perEpisode

	if err := s.runEpisodeStoryboardStage(runID, episode, stageBase); err != nil {
		return &workflowStageError{Stage: s.stageWithEpisodeNumber(workflowStageStoryboard, episode.EpisodeNum), EpisodeID: &episode.ID, Err: err}
	}
	if err := s.runEpisodeImageStage(runID, episode, stageBase+(perEpisode/4)); err != nil {
		return &workflowStageError{Stage: s.stageWithEpisodeNumber(workflowStageStoryboardImage, episode.EpisodeNum), EpisodeID: &episode.ID, Err: err}
	}
	if err := s.runEpisodeVideoStage(runID, episode, stageBase+(perEpisode/2)); err != nil {
		return &workflowStageError{Stage: s.stageWithEpisodeNumber(workflowStageStoryboardVideo, episode.EpisodeNum), EpisodeID: &episode.ID, Err: err}
	}
	if err := s.runEpisodeMergeStage(runID, episode, stageBase+((perEpisode*3)/4)); err != nil {
		return &workflowStageError{Stage: s.stageWithEpisodeNumber(workflowStageEpisodeMerge, episode.EpisodeNum), EpisodeID: &episode.ID, Err: err}
	}

	return nil
}

func (s *AutoWorkflowService) runEpisodeStoryboardStage(runID uint, episode *models.Episode, progress int) error {
	stage := workflowStageStoryboard
	s.updateRun(runID, models.WorkflowRunStatusProcessing, s.stageWithEpisodeNumber(stage, episode.EpisodeNum), progress, "", nil)

	var storyboardCount int64
	if err := s.db.Model(&models.Storyboard{}).Where("episode_id = ?", episode.ID).Count(&storyboardCount).Error; err != nil {
		return err
	}
	if storyboardCount > 0 {
		s.completeStage(runID, &episode.ID, stage, fmt.Sprintf("Episode %d storyboards already exist, skipping generation", episode.EpisodeNum), map[string]interface{}{"skipped": true, "count": storyboardCount})
		return nil
	}

	s.updateStep(runID, &episode.ID, stage, models.WorkflowRunStatusProcessing, 10, fmt.Sprintf("Generating storyboards for episode %d", episode.EpisodeNum), nil, nil)
	taskID, err := s.storyboardService.GenerateStoryboard(fmt.Sprintf("%d", episode.ID), "")
	if err != nil {
		return err
	}

	task, err := s.waitForAsyncTask(taskID, 15*time.Minute)
	if err != nil {
		return err
	}
	if task.Status != "completed" {
		return fmt.Errorf("episode %d storyboard generation failed: %s", episode.EpisodeNum, task.Error)
	}

	s.completeStage(runID, &episode.ID, stage, fmt.Sprintf("Episode %d storyboard generation completed", episode.EpisodeNum), nil)
	return nil
}

func (s *AutoWorkflowService) runEpisodeImageStage(runID uint, episode *models.Episode, progress int) error {
	stage := workflowStageStoryboardImage
	s.updateRun(runID, models.WorkflowRunStatusProcessing, s.stageWithEpisodeNumber(stage, episode.EpisodeNum), progress, "", nil)

	targetCount, completedCount, err := s.getEpisodeImageProgress(episode.ID)
	if err != nil {
		return err
	}
	if targetCount > 0 && completedCount >= targetCount {
		s.completeStage(runID, &episode.ID, stage, fmt.Sprintf("Episode %d storyboard images already complete, skipping generation", episode.EpisodeNum), map[string]interface{}{"skipped": true, "count": targetCount})
		return nil
	}

	s.updateStep(runID, &episode.ID, stage, models.WorkflowRunStatusProcessing, 10, fmt.Sprintf("Generating storyboard images for episode %d", episode.EpisodeNum), nil, nil)
	if _, err := s.imageService.BatchGenerateImagesForEpisode(fmt.Sprintf("%d", episode.ID)); err != nil {
		return err
	}
	if err := s.waitForEpisodeStoryboardImages(episode.ID, 20*time.Minute); err != nil {
		return err
	}

	s.completeStage(runID, &episode.ID, stage, fmt.Sprintf("Episode %d storyboard image generation completed", episode.EpisodeNum), nil)
	return nil
}

func (s *AutoWorkflowService) runEpisodeVideoStage(runID uint, episode *models.Episode, progress int) error {
	stage := workflowStageStoryboardVideo
	s.updateRun(runID, models.WorkflowRunStatusProcessing, s.stageWithEpisodeNumber(stage, episode.EpisodeNum), progress, "", nil)

	targetCount, completedCount, err := s.getEpisodeVideoProgress(episode.ID)
	if err != nil {
		return err
	}
	if targetCount > 0 && completedCount >= targetCount {
		s.completeStage(runID, &episode.ID, stage, fmt.Sprintf("Episode %d storyboard videos already complete, skipping generation", episode.EpisodeNum), map[string]interface{}{"skipped": true, "count": targetCount})
		return nil
	}

	s.updateStep(runID, &episode.ID, stage, models.WorkflowRunStatusProcessing, 10, fmt.Sprintf("Generating storyboard videos for episode %d", episode.EpisodeNum), nil, nil)
	if _, err := s.videoService.BatchGenerateVideosForEpisode(fmt.Sprintf("%d", episode.ID)); err != nil {
		return err
	}
	if err := s.waitForEpisodeVideos(episode.ID, 35*time.Minute); err != nil {
		return err
	}

	s.completeStage(runID, &episode.ID, stage, fmt.Sprintf("Episode %d storyboard video generation completed", episode.EpisodeNum), nil)
	return nil
}

func (s *AutoWorkflowService) runEpisodeMergeStage(runID uint, episode *models.Episode, progress int) error {
	stage := workflowStageEpisodeMerge
	s.updateRun(runID, models.WorkflowRunStatusProcessing, s.stageWithEpisodeNumber(stage, episode.EpisodeNum), progress, "", nil)

	s.updateStep(runID, &episode.ID, stage, models.WorkflowRunStatusProcessing, 10, fmt.Sprintf("Merging episode %d", episode.EpisodeNum), nil, nil)
	result, err := s.videoMergeService.FinalizeEpisode(fmt.Sprintf("%d", episode.ID), nil)
	if err != nil {
		return err
	}

	mergeIDAny, ok := result["merge_id"]
	if !ok {
		return fmt.Errorf("episode %d merge task missing merge_id", episode.EpisodeNum)
	}
	mergeID, err := toUint(mergeIDAny)
	if err != nil {
		return err
	}

	if err := s.waitForEpisodeMerge(mergeID, 20*time.Minute); err != nil {
		return err
	}
	s.completeStage(runID, &episode.ID, stage, fmt.Sprintf("Episode %d merge completed", episode.EpisodeNum), nil)
	return nil
}

func (s *AutoWorkflowService) ensureEpisodes(dramaID uint) ([]models.Episode, error) {
	var drama models.Drama
	if err := s.db.Preload("Episodes").First(&drama, dramaID).Error; err != nil {
		return nil, err
	}

	if len(drama.Episodes) > 0 {
		episodes := drama.Episodes
		sort.Slice(episodes, func(i, j int) bool {
			return episodes[i].EpisodeNum < episodes[j].EpisodeNum
		})
		return episodes, nil
	}

	scriptContent := ""
	if drama.Description != nil {
		scriptContent = strings.TrimSpace(*drama.Description)
	}
	if scriptContent == "" && len(drama.Metadata) > 0 {
		var metadata map[string]interface{}
		if err := json.Unmarshal(drama.Metadata, &metadata); err == nil {
			for _, key := range []string{"full_script", "script", "outline", "summary"} {
				if value, ok := metadata[key].(string); ok && strings.TrimSpace(value) != "" {
					scriptContent = strings.TrimSpace(value)
					break
				}
			}
		}
	}
	if scriptContent == "" {
		return nil, fmt.Errorf("drama description or metadata.full_script is required to auto-create episode")
	}

	title := fmt.Sprintf("%s Episode 1", drama.Title)
	episode := models.Episode{
		DramaID:       dramaID,
		EpisodeNum:    1,
		Title:         title,
		ScriptContent: &scriptContent,
		Status:        "draft",
	}
	if err := s.db.Create(&episode).Error; err != nil {
		return nil, err
	}
	s.db.Model(&models.Drama{}).Where("id = ?", dramaID).Update("total_episodes", 1)

	return []models.Episode{episode}, nil
}

func (s *AutoWorkflowService) waitForAsyncTask(taskID string, timeout time.Duration) (*models.AsyncTask, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		task, err := s.taskService.GetTask(taskID)
		if err != nil {
			return nil, err
		}
		switch task.Status {
		case "completed", "failed":
			return task, nil
		}
		time.Sleep(3 * time.Second)
	}
	return nil, fmt.Errorf("task %s timeout", taskID)
}

func (s *AutoWorkflowService) waitForCharacterBaselines(dramaID uint, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		var characters []models.Character
		if err := s.db.Where("drama_id = ?", dramaID).Find(&characters).Error; err != nil {
			return err
		}
		allReady := true
		for _, character := range characters {
			if !s.characterHasBaseline(&character) {
				allReady = false
				break
			}
		}
		if allReady {
			return nil
		}

		var pending int64
		if err := s.db.Model(&models.ImageGeneration{}).
			Where("drama_id = ? AND character_id IS NOT NULL AND status IN ?", dramaID, []models.ImageGenerationStatus{
				models.ImageStatusPending,
				models.ImageStatusProcessing,
			}).
			Count(&pending).Error; err != nil {
			return err
		}
		if pending == 0 {
			return fmt.Errorf("character baseline generation finished without producing a baseline sheet")
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("character baseline generation timeout")
}

func (s *AutoWorkflowService) waitForEpisodeStoryboardImages(episodeID uint, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		targetCount, completedCount, err := s.getEpisodeImageProgress(episodeID)
		if err != nil {
			return err
		}
		if targetCount > 0 && completedCount >= targetCount {
			return nil
		}

		var pending int64
		if err := s.db.Model(&models.ImageGeneration{}).
			Where("storyboard_id IN (?) AND status IN ?", s.storyboardSubQuery(episodeID), []models.ImageGenerationStatus{
				models.ImageStatusPending,
				models.ImageStatusProcessing,
			}).
			Count(&pending).Error; err != nil {
			return err
		}
		if pending == 0 && completedCount < targetCount {
			return fmt.Errorf("episode %d image generation incomplete", episodeID)
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("episode %d image generation timeout", episodeID)
}

func (s *AutoWorkflowService) waitForEpisodeVideos(episodeID uint, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		targetCount, completedCount, err := s.getEpisodeVideoProgress(episodeID)
		if err != nil {
			return err
		}
		if targetCount > 0 && completedCount >= targetCount {
			return nil
		}

		var pending int64
		if err := s.db.Model(&models.VideoGeneration{}).
			Where("storyboard_id IN (?) AND status IN ?", s.storyboardSubQuery(episodeID), []models.VideoStatus{
				models.VideoStatusPending,
				models.VideoStatusProcessing,
			}).
			Count(&pending).Error; err != nil {
			return err
		}
		if pending == 0 && completedCount < targetCount {
			return fmt.Errorf("episode %d video generation incomplete", episodeID)
		}
		time.Sleep(6 * time.Second)
	}
	return fmt.Errorf("episode %d video generation timeout", episodeID)
}

func (s *AutoWorkflowService) waitForEpisodeMerge(mergeID uint, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		merge, err := s.videoMergeService.GetMerge(mergeID)
		if err != nil {
			return err
		}
		switch merge.Status {
		case models.VideoMergeStatusCompleted:
			return nil
		case models.VideoMergeStatusFailed:
			if merge.ErrorMsg != nil {
				return fmt.Errorf(*merge.ErrorMsg)
			}
			return fmt.Errorf("video merge failed")
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("episode merge timeout")
}

func (s *AutoWorkflowService) getEpisodeImageProgress(episodeID uint) (int64, int64, error) {
	var targetCount int64
	if err := s.storyboardQueryForGeneration(episodeID).Count(&targetCount).Error; err != nil {
		return 0, 0, err
	}
	if targetCount == 0 {
		return 0, 0, nil
	}

	var completedCount int64
	if err := s.db.Model(&models.ImageGeneration{}).
		Where("storyboard_id IN (?) AND status = ?", s.storyboardSubQuery(episodeID), models.ImageStatusCompleted).
		Distinct("storyboard_id").
		Count(&completedCount).Error; err != nil {
		return 0, 0, err
	}
	return targetCount, completedCount, nil
}

func (s *AutoWorkflowService) getEpisodeVideoProgress(episodeID uint) (int64, int64, error) {
	var targetCount int64
	if err := s.storyboardQueryForGeneration(episodeID).Count(&targetCount).Error; err != nil {
		return 0, 0, err
	}
	if targetCount == 0 {
		return 0, 0, nil
	}

	var completedCount int64
	if err := s.db.Model(&models.VideoGeneration{}).
		Where("storyboard_id IN (?) AND status = ?", s.storyboardSubQuery(episodeID), models.VideoStatusCompleted).
		Distinct("storyboard_id").
		Count(&completedCount).Error; err != nil {
		return 0, 0, err
	}
	return targetCount, completedCount, nil
}

func (s *AutoWorkflowService) storyboardSubQuery(episodeID uint) *gorm.DB {
	return s.storyboardQueryForGeneration(episodeID).Select("id")
}

func (s *AutoWorkflowService) storyboardQueryForGeneration(episodeID uint) *gorm.DB {
	query := s.db.Model(&models.Storyboard{}).Where("episode_id = ?", episodeID)

	var primaryCount int64
	if err := s.db.Model(&models.Storyboard{}).
		Where("episode_id = ? AND is_primary = ?", episodeID, true).
		Count(&primaryCount).Error; err == nil && primaryCount > 0 {
		query = query.Where("is_primary = ?", true)
	}

	return query
}

func (s *AutoWorkflowService) buildEpisodeSummaries(dramaID uint) []WorkflowEpisodeSummary {
	var episodes []models.Episode
	if err := s.db.Where("drama_id = ?", dramaID).Order("episode_number ASC").Find(&episodes).Error; err != nil {
		return []WorkflowEpisodeSummary{}
	}

	summaries := make([]WorkflowEpisodeSummary, 0, len(episodes))
	for _, episode := range episodes {
		var storyboardCount int64
		var completedImages int64
		var completedVideos int64
		var failedImages int64
		var failedVideos int64
		_ = s.db.Model(&models.Storyboard{}).Where("episode_id = ?", episode.ID).Count(&storyboardCount).Error
		_ = s.db.Model(&models.ImageGeneration{}).Where("storyboard_id IN (?) AND status = ?", s.storyboardSubQuery(episode.ID), models.ImageStatusCompleted).Distinct("storyboard_id").Count(&completedImages).Error
		_ = s.db.Model(&models.VideoGeneration{}).Where("storyboard_id IN (?) AND status = ?", s.storyboardSubQuery(episode.ID), models.VideoStatusCompleted).Distinct("storyboard_id").Count(&completedVideos).Error
		_ = s.db.Model(&models.ImageGeneration{}).Where("storyboard_id IN (?) AND status = ?", s.storyboardSubQuery(episode.ID), models.ImageStatusFailed).Count(&failedImages).Error
		_ = s.db.Model(&models.VideoGeneration{}).Where("storyboard_id IN (?) AND status = ?", s.storyboardSubQuery(episode.ID), models.VideoStatusFailed).Count(&failedVideos).Error

		mergeStatus := ""
		var latestMerge models.VideoMerge
		if err := s.db.Where("episode_id = ?", episode.ID).Order("created_at DESC").First(&latestMerge).Error; err == nil {
			mergeStatus = string(latestMerge.Status)
		}

		finalVideoURL := ""
		if episode.VideoURL != nil {
			finalVideoURL = *episode.VideoURL
		}

		summaries = append(summaries, WorkflowEpisodeSummary{
			EpisodeID:         episode.ID,
			EpisodeNumber:     episode.EpisodeNum,
			Title:             episode.Title,
			Status:            episode.Status,
			StoryboardCount:   storyboardCount,
			CompletedImages:   completedImages,
			CompletedVideos:   completedVideos,
			FailedImages:      failedImages,
			FailedVideos:      failedVideos,
			MergeStatus:       mergeStatus,
			FinalVideoURL:     finalVideoURL,
			NeedsManualReview: failedImages > 0 || failedVideos > 0,
		})
	}

	return summaries
}

func (s *AutoWorkflowService) characterHasBaseline(character *models.Character) bool {
	if character == nil || character.ImageURL == nil || strings.TrimSpace(*character.ImageURL) == "" {
		return false
	}

	references := make([]string, 0, 3)
	if len(character.ReferenceImages) > 0 {
		_ = json.Unmarshal(character.ReferenceImages, &references)
	}

	for _, reference := range references {
		if strings.TrimSpace(reference) != "" {
			return true
		}
	}

	return true
}

func (s *AutoWorkflowService) updateRun(runID uint, status models.WorkflowRunStatus, stage string, progress int, message string, result map[string]interface{}) {
	updates := map[string]interface{}{
		"status":        status,
		"current_stage": stage,
		"progress":      progress,
		"updated_at":    time.Now(),
	}
	if message != "" {
		payload, _ := json.Marshal(map[string]string{"message": message})
		updates["result_json"] = datatypes.JSON(payload)
	}
	if result != nil {
		payload, _ := json.Marshal(result)
		updates["result_json"] = datatypes.JSON(payload)
	}
	_ = s.db.Model(&models.WorkflowRun{}).Where("id = ?", runID).Updates(updates).Error
}

func (s *AutoWorkflowService) updateStep(runID uint, episodeID *uint, stage string, status models.WorkflowRunStatus, progress int, message string, err error, meta map[string]interface{}) {
	var step models.WorkflowStepRun
	query := s.db.Where("workflow_run_id = ? AND stage = ?", runID, stage)
	if episodeID != nil {
		query = query.Where("episode_id = ?", *episodeID)
	} else {
		query = query.Where("episode_id IS NULL")
	}
	loadErr := query.First(&step).Error
	now := time.Now()

	var metaJSON datatypes.JSON
	if meta != nil {
		payload, _ := json.Marshal(meta)
		metaJSON = datatypes.JSON(payload)
	}

	errMsg := (*string)(nil)
	if err != nil {
		msg := err.Error()
		errMsg = &msg
	}

	if loadErr == gorm.ErrRecordNotFound {
		step = models.WorkflowStepRun{
			WorkflowRunID: runID,
			EpisodeID:     episodeID,
			Stage:         stage,
			Status:        status,
			Progress:      progress,
			Message:       message,
			ErrorMsg:      errMsg,
			MetaJSON:      metaJSON,
			StartedAt:     &now,
		}
		if status == models.WorkflowRunStatusCompleted || status == models.WorkflowRunStatusFailed {
			step.CompletedAt = &now
		}
		_ = s.db.Create(&step).Error
		return
	}
	if loadErr != nil {
		return
	}

	updates := map[string]interface{}{
		"status":     status,
		"progress":   progress,
		"message":    message,
		"meta_json":  metaJSON,
		"updated_at": now,
	}
	if step.StartedAt == nil {
		updates["started_at"] = &now
	}
	if errMsg != nil {
		updates["error_msg"] = *errMsg
	} else {
		updates["error_msg"] = nil
	}
	if status == models.WorkflowRunStatusCompleted || status == models.WorkflowRunStatusFailed {
		updates["completed_at"] = &now
	}
	_ = s.db.Model(&models.WorkflowStepRun{}).Where("id = ?", step.ID).Updates(updates).Error
}

func (s *AutoWorkflowService) completeStage(runID uint, episodeID *uint, stage string, message string, meta map[string]interface{}) {
	s.updateStep(runID, episodeID, stage, models.WorkflowRunStatusCompleted, 100, message, nil, meta)
}

func (s *AutoWorkflowService) failRun(runID uint, stage string, err error, episodeID *uint) {
	if err == nil {
		return
	}
	now := time.Now()
	message := err.Error()
	if stage != "" {
		updates := map[string]interface{}{
			"status":       models.WorkflowRunStatusFailed,
			"progress":     0,
			"message":      message,
			"error_msg":    message,
			"completed_at": &now,
			"updated_at":   now,
		}
		stepQuery := s.db.Model(&models.WorkflowStepRun{}).
			Where("workflow_run_id = ? AND stage = ? AND status <> ?", runID, stage, models.WorkflowRunStatusCompleted)
		if episodeID != nil {
			stepQuery = stepQuery.Where("episode_id = ?", *episodeID)
		}
		result := stepQuery.Updates(updates)
		if result.Error != nil || result.RowsAffected == 0 {
			s.updateStep(runID, episodeID, stage, models.WorkflowRunStatusFailed, 0, message, err, nil)
		}
	}
	_ = s.db.Model(&models.WorkflowRun{}).Where("id = ?", runID).Updates(map[string]interface{}{
		"status":        models.WorkflowRunStatusFailed,
		"current_stage": stage,
		"error_msg":     message,
		"completed_at":  &now,
		"updated_at":    now,
	}).Error
}

func (s *AutoWorkflowService) stageWithEpisodeNumber(stage string, episodeNumber int) string {
	return fmt.Sprintf("%s_episode_%d", stage, episodeNumber)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func toUint(value interface{}) (uint, error) {
	switch v := value.(type) {
	case uint:
		return v, nil
	case int:
		if v < 0 {
			return 0, fmt.Errorf("invalid negative id")
		}
		return uint(v), nil
	case int64:
		if v < 0 {
			return 0, fmt.Errorf("invalid negative id")
		}
		return uint(v), nil
	case float64:
		if v < 0 {
			return 0, fmt.Errorf("invalid negative id")
		}
		return uint(v), nil
	default:
		return 0, fmt.Errorf("unsupported id type")
	}
}
