package types

import (
	"time"
)

type Event struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	SourceID     uint      `json:"sourceId" gorm:"-"`
	Color        string    `json:"color"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Tags         []Tag     `json:"tags"`
	TimeOfStart  time.Time `json:"timeOfStart"`
	TimeOfFinish time.Time `json:"timeOfFinish"` // if equal with TimeOfStart - not for display
	UserEmail    string    `json:"-"`
	User         User      `json:"-" gorm:"references:Email;foreignKey:UserEmail;constraint:OnDelete:CASCADE;"`
}
