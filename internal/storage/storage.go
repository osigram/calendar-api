package storage

import (
	"calendar-api/internal/auth"
	"calendar-api/internal/services/events"
	"calendar-api/internal/services/extensions"
	"calendar-api/internal/services/tags"
)

type Storage interface {
	events.Storage
	auth.Storage
	tags.Storage
	extensions.Storage
}
