package storage

import (
	"calendar-api/handlers/events"
	"calendar-api/handlers/extensions"
	"calendar-api/handlers/tags"
	"calendar-api/middlewares"
)

type EventRepository interface {
	events.Adder
	events.ByIDGetter
	events.ByDateGetter
	events.Deleter
	events.Updater
}

type UserRepository interface {
	middlewares.UserGetSetter
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
