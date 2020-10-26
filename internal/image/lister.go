package image

import (
	"time"
)

type Details struct {
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
}

func newDetails(month time.Month, year int, location string) *Details {
	return &Details{
		Date:     time.Date(year, month, 0, 0, 0, 0, 0, time.UTC),
		Location: location,
	}
}

type Lister interface {
	Images() (map[string]*Details, error)
}

var staticDetails = map[string]*Details{
	"dents_du_midi_1.jpg":  newDetails(time.January, 2020, "Dents du Midi, Switzerland"),
	"dubai_1.jpg":          newDetails(time.June, 2017, "Dubai, UAE"),
	"fuji_1.jpg":           newDetails(time.October, 2017, "Mount Fuji, Japan"),
	"geneva_1.jpg":         newDetails(time.June, 2016, "Geneva, Switzerland"),
	"kyoto_1.jpg":          newDetails(time.October, 2017, "Kyoto, Japan"),
	"lhc_1.jpg":            newDetails(time.August, 2019, "LHC, France / Switzerland"),
	"malibu_1.jpg":         newDetails(time.March, 2019, "Malibu, USA"),
	"montreux_1.jpg":       newDetails(time.October, 2016, "Montreux, Switzerland"),
	"new_delhi_1.jpg":      newDetails(time.June, 2017, "New Delhi, India"),
	"newyork_2.jpg":        newDetails(time.August, 2015, "New York, USA"),
	"nuggets_point_1.jpg":  newDetails(time.January, 2019, "Nuggets Point, New Zealand"),
	"shenzhen_1.jpg":       newDetails(time.August, 2014, "Shenzhen, China"),
	"singapore_1.jpg":      newDetails(time.January, 2019, "Singapore"),
	"thun_1.jpg":           newDetails(time.May, 2016, "Thun, Switzerland"),
	"whaikiti_beach_1.jpg": newDetails(time.January, 2019, "Whaikiti Beach, New Zealand"),
}

type StaticLister struct{}

func (sl *StaticLister) Images() (map[string]*Details, error) {
	return staticDetails, nil
}
