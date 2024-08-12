package gormstorage

import (
	"calendar-api/internal/core"
	"fmt"
	"time"
)

func (gs *GormStorage) AddEvent(event *core.Event) error {
	db := gs.db

	result := db.Create(event)

	return result.Error
}

func (gs *GormStorage) GetEventByID(id uint) (*core.Event, error) {
	db := gs.db

	var event core.Event
	result := db.Preload("Tags").First(&event, id)

	return &event, result.Error
}

func (gs *GormStorage) GetEventsByDate(user *core.User, timeOfStart time.Time, timeOfFinish time.Time) ([]core.Event, error) {
	db := gs.db

	var events []core.Event
	result := db.Model(&core.Event{}).
		Preload("Tags").
		Where(&core.Event{UserEmail: user.Email}).
		Where("time_of_start > ?", timeOfStart).
		Where("time_of_finish < ?", timeOfFinish).
		Find(&events)

	return events, result.Error
}

func (gs *GormStorage) DeleteEvent(id uint) error {
	db := gs.db

	result := db.Delete(&core.Event{ID: id})

	return result.Error
}

func (gs *GormStorage) UpdateEvent(event *core.Event) error {
	db := gs.db

	if event == nil {
		return fmt.Errorf("empty event")
	}

	result := db.Model(event).Updates(*event)

	return result.Error
}
