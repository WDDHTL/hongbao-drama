package handlers

import (
	"github.com/drama-generator/backend/application/services"
	"github.com/drama-generator/backend/infrastructure/storage"
	"github.com/drama-generator/backend/pkg/config"
	"github.com/drama-generator/backend/pkg/logger"
	"github.com/drama-generator/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowHandler struct {
	workflowService *services.AutoWorkflowService
	log             *logger.Logger
}

type RetryVideoRequest struct {
	EpisodeID uint   `json:"episode_id" binding:"required"`
	Model     string `json:"model" binding:"required"`
}

func NewWorkflowHandler(db *gorm.DB, cfg *config.Config, transferService *services.ResourceTransferService, localStorage *storage.LocalStorage, log *logger.Logger) *WorkflowHandler {
	return &WorkflowHandler{
		workflowService: services.NewAutoWorkflowService(db, cfg, transferService, localStorage, log),
		log:             log,
	}
}

func (h *WorkflowHandler) StartProjectWorkflow(c *gin.Context) {
	dramaID := c.Param("id")
	run, err := h.workflowService.StartProjectWorkflow(dramaID)
	if err != nil {
		h.log.Errorw("Failed to start project workflow", "error", err, "drama_id", dramaID)
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, run)
}

func (h *WorkflowHandler) ResumeProjectWorkflow(c *gin.Context) {
	dramaID := c.Param("id")
	run, err := h.workflowService.ResumeProjectWorkflow(dramaID)
	if err != nil {
		h.log.Errorw("Failed to resume project workflow", "error", err, "drama_id", dramaID)
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, run)
}

func (h *WorkflowHandler) PauseProjectWorkflow(c *gin.Context) {
	dramaID := c.Param("id")
	run, err := h.workflowService.PauseProjectWorkflow(dramaID)
	if err != nil {
		h.log.Errorw("Failed to pause project workflow", "error", err, "drama_id", dramaID)
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, run)
}

func (h *WorkflowHandler) GetProjectWorkflowStatus(c *gin.Context) {
	dramaID := c.Param("id")
	status, err := h.workflowService.GetProjectWorkflowStatus(dramaID)
	if err != nil {
		h.log.Errorw("Failed to get project workflow status", "error", err, "drama_id", dramaID)
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, status)
}

// Retry video generation for a specific episode with a chosen model (after user confirmation)
func (h *WorkflowHandler) RetryEpisodeVideo(c *gin.Context) {
	dramaID := c.Param("id")
	var req RetryVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.workflowService.RetryEpisodeVideoWithModel(dramaID, req.EpisodeID, req.Model); err != nil {
		h.log.Errorw("Failed to retry episode video", "error", err, "drama_id", dramaID, "episode_id", req.EpisodeID, "model", req.Model)
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "retry started"})
}
