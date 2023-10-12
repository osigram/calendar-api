package extensionsmapping

import (
	"calendar-api/types"
	"time"
)

type Extension interface {
	GetEventByID(uint) (*types.Event, error)
	GetEventsByDate(string, time.Time, time.Time) ([]types.Event, error)
	ValidateAdditionalData(string) bool
	AdditionalDataOptions() ([]string, error)
}
