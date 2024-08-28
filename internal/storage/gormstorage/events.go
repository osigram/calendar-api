package gormstorage

import (
	"calendar-api/internal/core"
	"context"
	"fmt"
	"time"
)

func (gs *GormStorage) AddEvent(ctx context.Context, event *core.Event) error {
	db := gs.db.WithContext(ctx)

	result := db.Create(event)

	return result.Error
}

func (gs *GormStorage) GetEventByID(ctx context.Context, id uint) (*core.Event, error) {
	db := gs.db.WithContext(ctx)

	var event core.Event
	result := db.Preload("Tags").First(&event, id)

	return &event, result.Error
}

func (gs *GormStorage) GetEventsByDate(ctx context.Context, user *core.User, timeOfStart time.Time, timeOfFinish time.Time) ([]core.Event, error) {
	db := gs.db.WithContext(ctx)

	var events []core.Event
	result := db.Model(&core.Event{}).
		Preload("Tags").
		Where(&core.Event{UserEmail: user.Email}).
		Where("time_of_start > ?", timeOfStart).
		Where("time_of_finish < ?", timeOfFinish).
		Find(&events)

	return events, result.Error
}

func (gs *GormStorage) DeleteEvent(ctx context.Context, id uint) error {
	db := gs.db.WithContext(ctx)

	result := db.Delete(&core.Event{ID: id})

	return result.Error
}

func (gs *GormStorage) UpdateEvent(ctx context.Context, event *core.Event) error {
	db := gs.db.WithContext(ctx)

	if event == nil {
		return fmt.Errorf("empty event")
	}

	result := db.Model(event).Updates(*event)

	return result.Error
}
