package types

import "time"

type Event struct {
	Id           int64     `json:"id" db:"id"`
	SourceId     int64     `json:"sourceId" db:"-"`
	Color        string    `json:"color"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Tags         []Tag     `json:"tags"`
	TimeOfStart  time.Time `json:"timeOfStart" db:"time_of_start"`
	TimeOfFinish time.Time `json:"timeOfFinish" db:"time_of_finish"` // if equal with TimeOfStart - not for display
	User         User      `json:"-" db:"user"`
}
