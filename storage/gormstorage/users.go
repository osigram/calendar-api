package gormstorage

import "calendar-api/types"

func (gs *GormStorage) GetUser(email string) (*types.User, error) {
	db := gs.db

	var user types.User
	result := db.Joins("ExtensionsData").First(&user, email)

	return &user, result.Error
}

func (gs *GormStorage) AddUser(user *types.User) error {
	db := gs.db

	result := db.Create(user)

	return result.Error
}
