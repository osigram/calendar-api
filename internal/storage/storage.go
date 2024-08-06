package storage

import (
	events2 "calendar-api/internal/handlers/events"
	extensions2 "calendar-api/internal/handlers/extensions"
	tags2 "calendar-api/internal/handlers/tags"
	"calendar-api/internal/middlewares"
)

type EventRepository interface {
	events2.Adder
	events2.ByIDGetter
	events2.ByDateGetter
	events2.Deleter
	events2.Updater
}

type UserRepository interface {
	middlewares.UserGetSetter
}

type TagRepository interface {
	tags2.Adder
	tags2.Deleter
}

type ExtensionRepository interface {
	extensions2.Installer
	extensions2.Deleter
}

type Storage interface {
	EventRepository
	UserRepository
	TagRepository
	ExtensionRepository
}
