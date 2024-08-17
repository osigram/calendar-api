package storage

import (
	"calendar-api/internal/auth"
	"calendar-api/internal/services/events"
	"calendar-api/internal/services/extensions"
	"calendar-api/internal/services/tags"
)

type EventRepository interface {
	events.Adder
	events.ByIDGetter
	events.ByDateGetter
	events.Deleter
	events.Updater
}

type UserRepository interface {
	auth.Storage
}

type TagRepository interface {
	tags.Adder
	tags.Deleter
}

type ExtensionRepository interface {
	extensions.Installer
	extensions.Deleter
}

type Storage interface {
	EventRepository
	UserRepository
	TagRepository
	ExtensionRepository
}
