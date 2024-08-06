package gormstorage

import (
	types2 "calendar-api/internal/core"
	"fmt"
	"time"
)

func (gs *GormStorage) AddEvent(event *types2.Event) error {
	db := gs.db

	result := db.Create(event)

	return result.Error
}

func (gs *GormStorage) GetEventByID(id uint) (*types2.Event, error) {
	db := gs.db

	var event types2.Event
	result := db.Preload("Tags").First(&event, id)

	return &event, result.Error
}

func (gs *GormStorage) GetEventsByDate(user *types2.User, timeOfStart time.Time, timeOfFinish time.Time) ([]types2.Event, error) {
	db := gs.db

	var events []types2.Event
	result := db.Model(&types2.Event{}).
		Preload("Tags").
		Where(&types2.Event{UserEmail: user.Email}).
		Where("time_of_start > ?", timeOfStart).
		Where("time_of_finish < ?", timeOfFinish).
		Find(&events)

	return events, result.Error
}

func (gs *GormStorage) DeleteEvent(id uint) error {
	db := gs.db

	result := db.Delete(&types2.Event{ID: id})

	return result.Error
}

func (gs *GormStorage) UpdateEvent(event *types2.Event) error {
	db := gs.db

	if event == nil {
		return fmt.Errorf("empty event")
	}

	result := db.Model(event).Updates(*event)

	return result.Error
}
