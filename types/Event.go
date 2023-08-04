package types

import "time"

type Event struct {
	Id           int64     `json:"id"`
	SourceId     int64     `json:"sourceId"`
	Color        string    `json:"color"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Tags         []Tag     `json:"tags"`
	TimeOfStart  time.Time `json:"timeOfStart"`
	TimeOfFinish time.Time `json:"timeOfFinish"` // if equal with TimeOfStart - not for display
	User         User      `json:"user"`
}
