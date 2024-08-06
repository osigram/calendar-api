package gormstorage

import (
	"calendar-api/internal/core"
)

func (gs *GormStorage) GetUser(email string) (*core.User, error) {
	db := gs.db

	var user core.User
	result := db.Preload("ExtensionsData").First(&user, "email = ?", email)

	return &user, result.Error
}

func (gs *GormStorage) AddUser(user *core.User) error {
	db := gs.db

	result := db.Create(user)

	return result.Error
}
