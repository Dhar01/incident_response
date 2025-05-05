package model

import (
	"time"

	"gorm.io/gorm"
)

type Incident struct {
	AuthID    uint64         `gorm:"primaryKey" json:"authID,omitempty"`
	CreatedAt time.Time      `json:"createdAt,omitempty"`
	UpdatedAt time.Time      `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Title       string
	Description string
	Status      StatusType
	Severity    SeverityType
	AssignedTo  uint64
}

type IncidentReq struct {
	Title       string       `json:"title"`
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
