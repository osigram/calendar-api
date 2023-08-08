package postgres

import (
	"calendar-api/types"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type TempTag struct {
	EventId int64 `db:"event_id"`
	types.Tag
}

func (s *Storage) AddEvent(event *types.Event) error {
	_, err := s.db.NamedExec(
		`INSERT INTO events (color, name, description, time_of_start, time_of_finish, email) 
				VALUES (:color, :name, :description, :time_of_start, :time_of_finish, :user.email)`,
		event,
	)

	return err
}

func (s *Storage) DeleteEvent(id int64) error {
	_, err := s.db.Exec(
		`DELETE FROM events WHERE id = $1`,
		id,
	)

	return err
}

func (s *Storage) UpdateEvent(event *types.Event) error {
	_, err := s.db.NamedExec(
		`UPDATE events set color = :color, name = :name, description = :description, 
                  time_of_start = :time_of_start, time_of_finish = :time_of_finish
                  where id = :id`,
		event,
	)

	return err
}

func (s *Storage) GetEventById(id int64) (*types.Event, error) {
	var event types.Event
	err := s.db.Get(
		&event,
		`SELECT id, color, name, description, time_of_start, time_of_finish, email as "user.email"
		from events
        where id = $1`,
		id,
	)
	if err != nil {
		return nil, errors.New("event not found: " + err.Error())
	}

	err = s.db.Select(
		&(event.Tags),
		`SELECT id, tag_text from tags where event_id = $1`,
		id,
	)
	if err != nil {
		return nil, errors.New("event not found: " + err.Error())
	}

	return &event, nil
}

func (s *Storage) GetEventsByDate(user *types.User, timeOfStart, timeOfFinish time.Time) ([]types.Event, error) {
	var events []types.Event
	err := s.db.Select(
		&events,
		`SELECT id, color, name, description, time_of_start, time_of_finish, email as "user.email"
		from events
		where time_of_start between $1 AND $2 AND email = $3 order by events.id`,
		timeOfStart,
		timeOfFinish,
		user.Email,
	)
	if err != nil {
		return nil, errors.New("event not found: " + err.Error())
	} else if len(events) == 0 {
		return []types.Event{}, nil
	}

	ids := make([]int64, 0, len(events))
	for _, event := range events {
		ids = append(ids, event.Id)
	}

	query, args, err := sqlx.In(`SELECT id, tag_text, event_id from tags where event_id in (?) order by event_id`, ids)
	if err != nil {
		return nil, errors.New("err to create query for tags: " + err.Error())
	}
	var rawTags []TempTag
	err = s.db.Select(
		&rawTags,
		s.db.Rebind(query),
		args...,
	)
	if err != nil {
		return nil, errors.New("err to get tags: " + err.Error())
	}

	eventNum := 0
	for _, eventTag := range rawTags {
		if eventNum >= len(events) {
			break
		}
		for eventTag.EventId > events[eventNum].Id {
			eventNum++
			if eventNum >= len(events) {
				break
			}
		}
		events[eventNum].Tags = append(events[eventNum].Tags, eventTag.Tag)
	}

	return events, nil
}
