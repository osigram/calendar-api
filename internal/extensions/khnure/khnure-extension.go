package khnure

import (
	"calendar-api/internal/core"
)

type TimeTableExtension core.ExtensionInfo

func NewTimeTableExtension() *TimeTableExtension {
	return &TimeTableExtension{
		ID:          1,
		Name:        "Khnure TimeTable",
		Description: "Unofficial extension for KhNURE timetable.",
	}
}

func (t *TimeTableExtension) ExtensionInfo() core.ExtensionInfo {
	return core.ExtensionInfo(*t)
}
