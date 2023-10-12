package khnure

import (
	"calendar-api/types"
	"fmt"
	"time"
)

const DEFAULT_COLOR = ""
const DESCRIPTION_FORMAT = "Type: %v\nRoom: %v\nTeacher: %v"

var Styles = map[string]string{
	"Пз":  "#34566",
	"Лк":  "",
	"Лб":  "",
	"Екз": "",
	"Зал": "",
}

func (t *TimeTableExtension) GetEventsByDate(additionalData string, timeOfStart time.Time, timeOfFinish time.Time) ([]types.Event, error) {
	subjects, err := getSubjects(additionalData, timeOfStart, timeOfFinish)
	if err != nil {
		return nil, err
	}

	result := make([]types.Event, 0, len(subjects))
	for _, subject := range subjects {
		var color string
		var ok bool
		if color, ok = Styles[subject.Type]; !ok {
			color = DEFAULT_COLOR
		}

		result = append(result, types.Event{
			ID:          0, // TODO: what is ID?
			SourceID:    1,
			Color:       color,
			Name:        subject.Name,
			Description: fmt.Sprintf(DESCRIPTION_FORMAT, subject.Type, subject.Cabinet, subject.Teacher),
			Tags: []types.Tag{
				{TagText: subject.Type},
			},
			TimeOfStart:  subject.DatetimeStart,
			TimeOfFinish: subject.DatetimeEnd,
		})
	}

	return result, nil
}
