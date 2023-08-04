package extensions

import "calendar-api/handlers/events"

type Extension interface {
	events.ByIdGetter
	events.ByDateGetter
}
