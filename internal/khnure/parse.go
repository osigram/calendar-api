package khnure

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const GroupURL = "https://nure-dev.pp.ua/api/groups"
const URL = "https://nure-dev.pp.ua/api/schedule?type=group&id=%v&start_time=%v&end_time=%v"

type Teacher struct {
	ID        string
	FullName  string `json:"full_name"`
	ShortName string `json:"short_name"`
}

type Group struct {
	ID   string
	Name string
}

type APIEvent struct {
	ID         string
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Auditory   string
	NumberPair uint `json:"number_pair"`
	Type       string
	UpdatedAt  time.Time `json:"updatedAt"`
	Groups     []Group
	Teachers   []Teacher
	Subject    struct {
		ID    string
		Brief string // Short name of the subject
		Title string // Full name
	}
}

func parseGroups() ([]Group, error) {
	resp, err := http.Get(GroupURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var groups []Group
	err = json.Unmarshal(body, &groups)

	return groups, err
}

func getSubjects(additionalData string, startTime time.Time, endTime time.Time) ([]APIEvent, error) {
	resp, err := http.Get(fmt.Sprintf(URL, additionalData, startTime.Unix(), endTime.Unix()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiEvents []APIEvent
	err = json.Unmarshal(body, &apiEvents)

	return apiEvents, err
}
