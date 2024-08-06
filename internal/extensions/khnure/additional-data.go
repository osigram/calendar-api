package khnure

import "fmt"

type BaseParseStruct struct {
	ShortName string `json:"short_name"`
	FullName  string `json:"full_name"`
}

func (t *TimeTableExtension) AdditionalDataOptions() (map[string]string, error) {
	groups, err := parseGroups()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, group := range groups {
		result[group.Name] = fmt.Sprint(group.ID)
	}

	return result, nil
}

func (t *TimeTableExtension) ValidateAdditionalData(additionalData string) bool {
	options, err := t.AdditionalDataOptions()
	if err != nil {
		return false
	}

	for _, id := range options {
		if id == additionalData {
			return true
		}
	}

	return false
}
