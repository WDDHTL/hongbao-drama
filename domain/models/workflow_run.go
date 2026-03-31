package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type WorkflowRunStatus string

const (
	WorkflowRunStatusPending    WorkflowRunStatus = "pending"
	WorkflowRunStatusProcessing WorkflowRunStatus = "processing"
	WorkflowRunStatusCompleted  WorkflowRunStatus = "completed"
	WorkflowRunStatusFailed     WorkflowRunStatus = "failed"
	WorkflowRunStatusPaused     WorkflowRunStatus = "paused"
)

type WorkflowRunScope string

const (
	WorkflowRunScopeProject WorkflowRunScope = "project"
	WorkflowRunScopeEpisode WorkflowRunScope = "episode"
)

type WorkflowRun struct {
	ID           uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	DramaID      uint              `gorm:"not null;index" json:"drama_id"`
	EpisodeID    *uint             `gorm:"index" json:"episode_id,omitempty"`
	Scope        WorkflowRunScope  `gorm:"type:varchar(20);not null;index" json:"scope"`
	Status       WorkflowRunStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	CurrentStage string            `gorm:"type:varchar(100);index" json:"current_stage"`
	Progress     int               `gorm:"default:0" json:"progress"`
	ConfigJSON   datatypes.JSON    `gorm:"type:json" json:"config_json,omitempty"`
	ResultJSON   datatypes.JSON    `gorm:"type:json" json:"result_json,omitempty"`
	ErrorMsg     *string           `gorm:"type:text" json:"error_msg,omitempty"`
	StartedAt    *time.Time        `json:"started_at,omitempty"`
	CompletedAt  *time.Time        `json:"completed_at,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	DeletedAt    gorm.DeletedAt    `gorm:"index" json:"-"`

	Steps []WorkflowStepRun `gorm:"foreignKey:WorkflowRunID" json:"steps,omitempty"`
}

func (WorkflowRun) TableName() string {
	return "workflow_runs"
}

type WorkflowStepRun struct {
	ID            uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	WorkflowRunID uint              `gorm:"not null;index" json:"workflow_run_id"`
	EpisodeID     *uint             `gorm:"index" json:"episode_id,omitempty"`
	Stage         string            `gorm:"type:varchar(100);not null;index" json:"stage"`
	Status        WorkflowRunStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	Progress      int               `gorm:"default:0" json:"progress"`
	Message       string            `gorm:"type:varchar(500)" json:"message,omitempty"`
	ErrorMsg      *string           `gorm:"type:text" json:"error_msg,omitempty"`
	MetaJSON      datatypes.JSON    `gorm:"type:json" json:"meta_json,omitempty"`
	StartedAt     *time.Time        `json:"started_at,omitempty"`
	CompletedAt   *time.Time        `json:"completed_at,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `gorm:"index" json:"-"`
}

func (WorkflowStepRun) TableName() string {
	return "workflow_step_runs"
}
