package storage

import (
	"calendar-api/handlers/events"
	"calendar-api/middlewares"
)

type Storage interface {
	events.Adder
	events.ByIDGetter
	events.ByDateGetter
	events.Deleter
	events.Updater
	middlewares.UserGetSetter
}
