package extensions

import (
	"calendar-api/types"
	"time"
)

type Extension interface {
	GetEventById(int64) (*types.Event, error)
	GetEventsByDate(*types.User, time.Time, time.Time) ([]types.Event, error)
}
