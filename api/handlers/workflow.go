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
