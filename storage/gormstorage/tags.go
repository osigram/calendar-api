package gormstorage

import "calendar-api/types"

func (gs *GormStorage) AddTag(text string, eventID uint) error {
	db := gs.db

	tag := &types.Tag{
		EventID: eventID,
		TagText: text,
	}
	result := db.Create(tag)

	return result.Error
}

func (gs *GormStorage) DeleteTag(id uint) error {
	db := gs.db

	result := db.Delete(&types.Tag{ID: id})

	return result.Error
}
