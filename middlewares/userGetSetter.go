package middlewares

import "calendar-api/types"

type UserGetSetter interface {
	GetUser(email string) (*types.User, error)
	AddUser(user *types.User) error
}
