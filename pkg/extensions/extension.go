package extensions

import (
	"calendar-api/internal/core"
	"time"
)

type Extension interface {
	GetEventByID(uint) (*core.Event, error)
	GetEventsByDate(string, time.Time, time.Time) ([]core.Event, error)
	ValidateAdditionalData(string) bool
	AdditionalDataOptions() (map[string]string, error)
	ExtensionInfo() core.ExtensionInfo
}

type Getter interface {
	Get(uint) (Extension, error)
}
