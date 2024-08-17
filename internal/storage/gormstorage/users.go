package gormstorage

import (
	"calendar-api/internal/core"
	"fmt"
)

func (gs *GormStorage) GetUser(email string) (*core.User, error) {
	db := gs.db

	var user core.User
	result := db.Preload("ExtensionsData").Preload("Sessions").First(&user, "email = ?", email)

	return &user, result.Error
}

func (gs *GormStorage) AddUser(user *core.User) error {
	db := gs.db

	result := db.Create(user)

	return result.Error
}

func (gs *GormStorage) AddSession(session *core.Session) error {
	db := gs.db

	result := db.Create(session)

	return result.Error
}

func (gs *GormStorage) UpdateUser(user *core.User) error {
	db := gs.db

	if user == nil {
		return fmt.Errorf("empty user")
	}

	result := db.Model(user).Updates(*user)

	return result.Error
}

func (gs *GormStorage) UpdateSession(session *core.Session) error {
	db := gs.db

	if session == nil {
		return fmt.Errorf("empty session")
	}

	result := db.Model(session).Updates(*session)

	return result.Error
}

func (gs *GormStorage) DeleteSession(sessionID uint) error {
	db := gs.db

	result := db.Delete(&core.Session{ID: sessionID})

	return result.Error
}
