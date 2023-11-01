package types

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"math"
	"time"
)

type Event struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	SourceID     uint      `json:"sourceID" gorm:"-"`
	Color        string    `json:"color"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Tags         []Tag     `json:"tags"`
	TimeOfStart  time.Time `json:"timeOfStart"`
	TimeOfFinish time.Time `json:"timeOfFinish"` // if equal with TimeOfStart - not for display
	UserEmail    string    `json:"-"`
	User         User      `json:"-" gorm:"references:Email;foreignKey:UserEmail;constraint:OnDelete:CASCADE;"`
}

func (event *Event) Validate() error {
	return validation.ValidateStruct(event,
		validation.Field(&event.Color, validation.Required, is.HexColor),
		validation.Field(&event.TimeOfStart, validation.Min(time.Unix(0, 0))),
		validation.Field(&event.TimeOfFinish, validation.Max(time.Unix(math.MaxInt32, 0))),
		validation.Field(&event.Tags),
	)
}
