package khnure

import "calendar-api/types"

type TimeTableExtension types.ExtensionInfo

func NewTimeTableExtension() *TimeTableExtension {
	return &TimeTableExtension{
		ID:          1,
		Name:        "Khnure TimeTable",
		Description: "Unofficial extension for KhNURE timetable.",
	}
}

func (t *TimeTableExtension) ExtensionInfo() types.ExtensionInfo {
	return types.ExtensionInfo(*t)
}
