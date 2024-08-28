package extensions

import (
	"calendar-api/internal/core"
	"context"
	"time"
)

type Extension interface {
	GetEventByID(context.Context, uint) (*core.Event, error)
	GetEventsByDate(context.Context, string, time.Time, time.Time) ([]core.Event, error)
	ValidateAdditionalData(context.Context, string) bool
	AdditionalDataOptions(context.Context) (map[string]string, error)
	ExtensionInfo() core.ExtensionInfo
}

type Getter interface {
	Get(context.Context, uint) (Extension, error)
}
