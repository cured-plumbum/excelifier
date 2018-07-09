package main

import (
	"encoding/xml"
	"io"
	"strings"
	"time"
)

// TaskElement represents task element in XML
type TaskElement struct {
	UID           string
	Name          string
	WBS           string
	Active        int
	OutlineNumber string // it seems to be the same as WBS
	OutlineLevel  int
	Duration      string
}

// DurationTime returns duration of the task
func (t *TaskElement) DurationTime() (time.Duration, error) {
	return time.ParseDuration(strings.ToLower(strings.TrimLeft(t.Duration, "PT")))
}

// Hours returns duration, hours
// It returns 0.0, if there is an error occurred
func (t *TaskElement) Hours() float64 {
	d, err := t.DurationTime()
	if err != nil {
		return 0.0
	}
	return d.Hours()
}

func extractTasks(r io.Reader) ([]*TaskElement, error) {
	tt := make([]*TaskElement, 0)
	decoder := xml.NewDecoder(r)
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			return tt, nil
		}

		switch se := tok.(type) {
		case xml.StartElement:
			if se.Name.Local == "Task" {
				t := new(TaskElement)
				err = decoder.DecodeElement(t, &se)
				if err != nil {
					return nil, err
				}
				tt = append(tt, t)
			}
		}

	}
}
