package khnure

import (
	"calendar-api/types"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const DefaultColor = "#98dce5"
const DescriptionFormat = "Full Name: %v\nType: %v\nFor: %v\nAuditory: %v\nTeachers: %v"

var Styles = map[string]string{
	"Пз":   "#8cdb5f",
	"Лк":   "#fdcd84",
	"Лб":   "#d9a6f1",
	"Екз":  "#f27674",
	"Зал":  "#d8e147",
	"Конс": "#98dce5",
}

func newEventFromAPIEvent(subject APIEvent) (*types.Event, error) {
	var color string
	var ok bool
	if color, ok = Styles[subject.Type]; !ok {
		color = DefaultColor
	}

	groups := make([]string, 0, len(subject.Groups))
	for _, group := range subject.Groups {
		groups = append(groups, group.Name)
	}
	teachers := make([]string, 0, len(subject.Teachers))
	for _, teacher := range subject.Teachers {
		teachers = append(teachers, teacher.FullName)
	}
	if len(subject.Groups) == 0 {
		return nil, fmt.Errorf("internal KhNURE API error: no groups in event, skipping")
	}

	rawGroupID, err := strconv.Atoi(subject.Groups[0].ID)
	if err != nil {
		return nil, err
	}
	groupID := uint(rawGroupID)
	rawSubjectID, err := strconv.Atoi(subject.ID)
	if err != nil {
		return nil, err
	}
	subjectID := uint(rawSubjectID)

	startTimeUnix, err := strconv.Atoi(subject.StartTime)
	if err != nil {
		return nil, err
	}
	startTime := time.Unix(int64(startTimeUnix), 0)

	endTimeUnix, err := strconv.Atoi(subject.EndTime)
	if err != nil {
		return nil, err
	}
	endTime := time.Unix(int64(endTimeUnix), 0)

	return &types.Event{
		ID:       groupID*uint(math.Pow10(10)) + subjectID,
		SourceID: 1,
		Color:    color,
		Name:     subject.Subject.Brief,
		Description: fmt.Sprintf(DescriptionFormat,
			subject.Subject.Title,
			subject.Type,
			strings.Join(groups, ", "),
			subject.Auditory,
			strings.Join(teachers, ", "),
		),
		Tags: []types.Tag{
			{TagText: subject.Type},
			{TagText: subject.Subject.Brief},
		},
		TimeOfStart:  startTime,
		TimeOfFinish: endTime,
	}, nil
}

func (t *TimeTableExtension) GetEventsByDate(additionalData string, timeOfStart time.Time, timeOfFinish time.Time) ([]types.Event, error) {
	subjects, err := getSubjects(additionalData, timeOfStart, timeOfFinish)
	if err != nil {
		return nil, fmt.Errorf("error in getting data from KhNURE API: %v", err.Error())
	}

	result := make([]types.Event, 0, len(subjects))
	for _, subject := range subjects {
		event, err := newEventFromAPIEvent(subject)
		if err != nil {
			continue
		}
		result = append(result, *event)
	}

	return result, nil
}

func (t *TimeTableExtension) GetEventByID(id uint) (*types.Event, error) {
	groupID := id / uint(math.Pow10(10))
	eventID := id % uint(math.Pow10(10))

	year, month, _ := time.Now().Date()
	var startTime time.Time
	var endTime time.Time
	switch month {
	case 8, 9, 10, 11, 12:
		startTime = time.Date(year, 8, 1, 0, 0, 0, 0, time.Local)
		endTime = time.Date(year+1, 2, 1, 0, 0, 0, 0, time.Local)
	case 1:
		startTime = time.Date(year-1, 9, 1, 0, 0, 0, 0, time.Local)
		endTime = time.Date(year, 2, 1, 0, 0, 0, 0, time.Local)
	case 2, 3, 4, 5, 6, 7:
		startTime = time.Date(year, 2, 1, 0, 0, 0, 0, time.Local)
		endTime = time.Date(year, 7, 31, 0, 0, 0, 0, time.Local)
	}

	subjects, err := getSubjects(fmt.Sprint(groupID), startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error in getting data from KhNURE API: %v", err.Error())
	}
	for _, subject := range subjects {
		rawSubjectID, err := strconv.Atoi(subject.ID)
		if err != nil {
			continue
		}
		subjectID := uint(rawSubjectID)
		if subjectID == eventID {
			event, err := newEventFromAPIEvent(subject)

			return event, err
		}
	}

	return nil, fmt.Errorf("not found")
}
