package context

import (
	"calendar-api/internal/core"
	"context"
)

type Context struct {
	context.Context
	User      *core.User
	UserAgent string
}
