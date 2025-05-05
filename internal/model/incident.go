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

	Title   string
	Details string
	UserID  uint
	User    User
}
