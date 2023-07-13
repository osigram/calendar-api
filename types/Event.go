package types

import "time"

type Event struct {
	SourceName   string    `json:"sourceName"`
	Color        string    `json:"color"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Tags         []string  `json:"tags"`
	TimeOfStart  time.Time `json:"timeOfStart"`
	TimeOfFinish time.Time `json:"timeOfFinish"` // if equal with TimeOfStart - not for display
}
