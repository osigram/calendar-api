package khnure

import (
	"encoding/csv"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const URL = "https://cist.nure.ua/ias/app/tt/WEB_IAS_TT_GNR_RASP.GEN_GROUP_POTOK_RASP?ATypeDoc=3&Aid_group=%v&Aid_potok=0&ADateStart=%v&ADateEnd=%v&AMultiWorkSheet=0"

type Subject struct {
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Cabinet       string    `json:"cabinet"`
	DatetimeStart time.Time `json:"datetime_start"`
	DatetimeEnd   time.Time `json:"datetime_end"`
	Teacher       string    `json:"teacher"`
}

func getSubjects(groupID string, timeOfStart time.Time, timeOfFinish time.Time) ([]Subject, error) {
	resp, err := http.Get(fmt.Sprintf(URL, groupID, timeOfStart.Format("02.01.2006"), timeOfFinish.Format("02.01.2006")))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	utfReader := charmap.Windows1251.NewDecoder().Reader(resp.Body)
	utfBytes, err := io.ReadAll(utfReader)
	if err != nil {
		return nil, err
	}
	utfString := strings.ReplaceAll(string(utfBytes), "\r", "\n")
	csvReader := csv.NewReader(strings.NewReader(utfString))
	csvReader.LazyQuotes = true
	rawTimeTable, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	timeTable := parse(rawTimeTable)

	return timeTable, nil
}

func parse(rawTimetable [][]string) []Subject {
	if len(rawTimetable) == 0 {
		panic("rawTimetable is nil, PANIC!!!!")
	}

	result := make([]Subject, 0, len(rawTimetable)-1)

	for num, subj := range rawTimetable {
		if num == 0 || // skip timetable heading
			subj == nil || // for panic-safety
			len(subj) < 5 || // must have name, date of start, time of start, date of end, time of end...
			len(subj[0]) < 1 { // must not be empty name-type-cabinet-group
			continue
		}

		subjFirst := subj[0] // name-type-cabinet-group
		subj = subj[1:]      // subjFirst now has its own var

		// For several subjects at the same time
		if strings.Contains(subjFirst, "; ") {
			subjFirsts := strings.Split(subjFirst, "; ")
			for _, subjFirstElem := range subjFirsts {
				if len(subjFirstElem) < 1 {
					continue
				}

				exportSubject, err := parseOne(subjFirstElem, subj)

				if err != nil {
					log.Println(err)
					continue
				}

				result = append(result, exportSubject)
			}

			continue
		}

		exportSubject, err := parseOne(subjFirst, subj)

		if err != nil {
			log.Println(err)
			continue
		}

		result = append(result, exportSubject)
	}

	return result
}

func parseOne(subjFirst string, subj []string) (Subject, error) {
	if strings.Contains(subjFirst, "\"") {
		subjFirst = subjFirst[1:]
	}

	subjFirstSlice := strings.Split(subjFirst, " ")

	if len(subjFirstSlice) < 3 {
		return Subject{}, errors.New("first element of subj is less than 3")
	}

	subjName := subjFirstSlice[0]
	subjType := subjFirstSlice[1]
	subjCabinet := subjFirstSlice[2]

	if strings.Contains(subjCabinet, ",") {
		subjCabinet = subjCabinet[:len(subjCabinet)-1]
	}

	datetimeStart, err := time.Parse("02.01.2006 15:04:05", fmt.Sprintf("%v %v", subj[0], subj[1]))
	if err != nil {
		return Subject{}, fmt.Errorf("date parsing error: %v", err)
	}

	datetimeEnd, err := time.Parse("02.01.2006 15:04:05", fmt.Sprintf("%v %v", subj[2], subj[3]))
	if err != nil {
		return Subject{}, fmt.Errorf("date parsing error: %v", err)
	}

	exportSubject := Subject{
		Name:          subjName,
		Type:          subjType,
		Cabinet:       subjCabinet,
		DatetimeStart: datetimeStart,
		DatetimeEnd:   datetimeEnd,
		Teacher:       "",
	}

	return exportSubject, nil
}
