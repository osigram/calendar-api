package gormstorage

import (
	"calendar-api/internal/core"
)

func (gs *GormStorage) AddTag(text string, eventID uint) error {
	db := gs.db

	tag := &core.Tag{
		EventID: eventID,
		TagText: text,
	}
	result := db.Create(tag)

	return result.Error
}

func (gs *GormStorage) DeleteTag(id uint) error {
	db := gs.db

	result := db.Delete(&core.Tag{ID: id})

	return result.Error
}
