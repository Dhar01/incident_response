package model

import (
	"time"

	"gorm.io/gorm"
)

type Incident struct {
	IncidentID uint64         `gorm:"primaryKey"`
	CreatedAt  time.Time      `json:"createdAt,omitempty"`
	UpdatedAt  time.Time      `json:"updatedAt,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Title       string       `gorm:"not null"`
	Description string       `gorm:"type:text"`
	Status      StatusType   `gorm:"default:'open'"`
	Severity    SeverityType `gorm:"default:'medium'"`

	AuthID     uint64 `gorm:"not null"`
	AssignedTo uint64 `gorm:"not null"`

	Creator  Auth `gorm:"foreignKey:AuthID"`
	Assignee Auth `gorm:"foreignKey:AssignedTo"`
}

type IncidentReq struct {
	Title       string       `json:"title" validate:"required"`
	Description string       `json:"description"`
	Status      StatusType   `json:"status"`
	Severity    SeverityType `json:"severity"`
	AssignedTo  uint64       `json:"assigned_to"`
}

type SeverityType string

const (
	Low      SeverityType = "low"
	Medium   SeverityType = "medium"
	High     SeverityType = "high"
	Critical SeverityType = "critical"
)

type StatusType string

const (
	Open         StatusType = "open"
	Acknowledged StatusType = "acknowledged"
	Closed       StatusType = "closed"
)
