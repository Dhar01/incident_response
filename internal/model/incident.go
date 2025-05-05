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

type SeverityType int

const (
	Low SeverityType = iota
	Medium
	High
	Critical
)

type StatusType int

const (
	Open StatusType = iota
	Acknowledged
	Closed
)
