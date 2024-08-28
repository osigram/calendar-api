package gormstorage

import (
	"calendar-api/internal/core"
	"context"
)

func (gs *GormStorage) AddTag(ctx context.Context, text string, eventID uint) error {
	db := gs.db.WithContext(ctx)

	tag := &core.Tag{
		EventID: eventID,
		TagText: text,
	}
	result := db.Create(tag)

	return result.Error
}

func (gs *GormStorage) DeleteTag(ctx context.Context, id uint) error {
	db := gs.db.WithContext(ctx)

	result := db.Delete(&core.Tag{ID: id})

	return result.Error
}
