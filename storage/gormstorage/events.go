package gormstorage

import (
	"calendar-api/types"
	"time"
)

func (gs *GormStorage) AddEvent(event *types.Event) error {
	db := gs.db

	result := db.Create(event)

	return result.Error
}

func (gs *GormStorage) GetEventByID(id uint) (*types.Event, error) {
	db := gs.db

	var event types.Event
	result := db.Joins("Tags").First(&event, id)

	return &event, result.Error
}

func (gs *GormStorage) GetEventsByDate(user *types.User, timeOfStart time.Time, timeOfFinish time.Time) ([]types.Event, error) {
	db := gs.db

	var events []types.Event
	result := db.Model(&types.Event{}).
		Joins("Tags").
		Where(&types.Event{UserEmail: user.Email}).
		Where("time_of_start > ?", timeOfStart).
		Where("time_of_finish < ?", timeOfFinish).
		Scan(&events)

	return events, result.Error
}

func (gs *GormStorage) DeleteEvent(id uint) error {
	db := gs.db

	result := db.Delete(&types.Event{ID: id})

	return result.Error
}

func (gs *GormStorage) UpdateEvent(event *types.Event) error {
	db := gs.db

	result := db.Save(event)

	return result.Error
}
