package extensionsmapping

import (
	types2 "calendar-api/internal/core"
	"time"
)

type Extension interface {
	GetEventByID(uint) (*types2.Event, error)
	GetEventsByDate(string, time.Time, time.Time) ([]types2.Event, error)
	ValidateAdditionalData(string) bool
	AdditionalDataOptions() (map[string]string, error)
	ExtensionInfo() types2.ExtensionInfo
}
