package middlewares

import (
	"calendar-api/internal/core"
)

type UserGetSetter interface {
	GetUser(email string) (*core.User, error)
	AddUser(user *core.User) error
}
