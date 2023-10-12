package khnure

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"net/http"
)

const GROUP_URL = "https://cist.nure.ua/ias/app/tt/P_API_GROUP_JSON"

type BaseParseStruct struct {
	ShortName string
	FullName  string
}

type GroupParseData struct {
	University struct {
		BaseParseStruct
		Faculties []struct {
			ID uint
			BaseParseStruct
			Directions []struct {
				ID uint
				BaseParseStruct
				Specialities []struct {
					ID uint
					BaseParseStruct
					Groups []Group
				}
			}
		}
	}
}

type Group struct {
	ID   uint
	Name string
}

func parseGroups() (GroupParseData, error) {
	var data GroupParseData
	resp, err := http.Get(GROUP_URL)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	utfReader := charmap.Windows1251.NewDecoder().Reader(resp.Body)
	utfBytes, err := io.ReadAll(utfReader)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(utfBytes, &data)
	return data, err
}

func (t *TimeTableExtension) AdditionalDataOptions() (map[string]string, error) {
	data, err := parseGroups()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, faculty := range data.University.Faculties {
		for _, directory := range faculty.Directions {
			for _, speciality := range directory.Specialities {
				for _, group := range speciality.Groups {
					result[group.Name] = fmt.Sprint(group.ID)
				}
			}
		}
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
