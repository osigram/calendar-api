package gormstorage

import (
	"calendar-api/internal/core"
	"context"
	"fmt"
)

func (gs *GormStorage) GetUser(ctx context.Context, email string) (*core.User, error) {
	db := gs.db.WithContext(ctx)

	var user core.User
	result := db.Preload("ExtensionsData").Preload("Sessions").First(&user, "email = ?", email)

	return &user, result.Error
}

func (gs *GormStorage) AddUser(ctx context.Context, user *core.User) error {
	db := gs.db.WithContext(ctx)

	result := db.Create(user)

	return result.Error
}

func (gs *GormStorage) AddSession(ctx context.Context, session *core.Session) error {
	db := gs.db.WithContext(ctx)

	result := db.Create(session)

	return result.Error
}

func (gs *GormStorage) UpdateUser(ctx context.Context, user *core.User) error {
	db := gs.db.WithContext(ctx)

	if user == nil {
		return fmt.Errorf("empty user")
	}

	result := db.Model(user).Updates(*user)

	return result.Error
}

func (gs *GormStorage) UpdateSession(ctx context.Context, session *core.Session) error {
	db := gs.db.WithContext(ctx)

	if session == nil {
		return fmt.Errorf("empty session")
	}

	result := db.Model(session).Updates(*session)

	return result.Error
}

func (gs *GormStorage) DeleteSession(ctx context.Context, sessionID uint) error {
	db := gs.db.WithContext(ctx)

	result := db.Delete(&core.Session{ID: sessionID})

	return result.Error
}
