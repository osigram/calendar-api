package context

import (
	"calendar-api/internal/core"
	"golang.org/x/net/context"
)

type Context struct {
	context.Context
	User *core.User
}
