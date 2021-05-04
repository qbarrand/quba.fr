package images

import (
	"embed"
	"fmt"
	"io/fs"
	"time"
)

type Metadata struct {
	Date      time.Time
	Location  string
	MainColor string
}

func newMetadata(month time.Month, year int, location, mainColor string) *Metadata {
	return &Metadata{
		Date:      time.Date(year, month, 1, 0, 0, 0, 0, time.UTC),
		Location:  location,
		MainColor: mainColor,
	}
}

var db = map[string]*Metadata{
	"dents_du_midi_1.jpg":  newMetadata(time.January, 2020, "Dents du Midi, Switzerland", "#4279AC"),
	"dubai_1.jpg":          newMetadata(time.June, 2017, "Dubai, UAE", "#4F4C42"),
	"fuji_1.jpg":           newMetadata(time.October, 2017, "Mount Fuji, Japan", "#68809E"),
	"geneva_1.jpg":         newMetadata(time.June, 2016, "Geneva, Switzerland", "#717D85"),
	"kyoto_1.jpg":          newMetadata(time.October, 2017, "Kyoto, Japan", "#B34B2C"),
	"lhc_1.jpg":            newMetadata(time.August, 2019, "LHC, France / Switzerland", "#827366"),
	"malibu_1.jpg":         newMetadata(time.March, 2019, "Malibu, USA", "#744930"),
	"montreux_1.jpg":       newMetadata(time.October, 2016, "Montreux, Switzerland", "#768692"),
	"new_delhi_1.jpg":      newMetadata(time.June, 2017, "New Delhi, India", "#807B66"),
	"newyork_2.jpg":        newMetadata(time.August, 2015, "New York, USA", "#7F8B8E"),
	"nuggets_point_1.jpg":  newMetadata(time.January, 2019, "Nuggets Point, New Zealand", "#526F7F"),
	"shenzhen_1.jpg":       newMetadata(time.August, 2014, "Shenzhen, China", "#5C0C1A"),
	"singapore_1.jpg":      newMetadata(time.January, 2019, "Singapore", "#483D39"),
	"thun_1.jpg":           newMetadata(time.May, 2016, "Thun, Switzerland", "#587FA4"),
	"whaikiti_beach_1.jpg": newMetadata(time.January, 2019, "Whaikiti Beach, New Zealand", "#6C8296"),
	"zermatt_1.jpg":        newMetadata(time.February, 2021, "Zermatt, Switzerland", "#3B6796"),
}

type MetadataFS interface {
	fs.FS

	OpenWithMetadata(string) (fs.File, *Metadata, error)
}

//go:embed *.jpg
var local embed.FS

type embedded struct {
	fs fs.FS
}

func LocalImagesWithMetadata() MetadataFS {
	return &embedded{fs: local}
}

// Open implements fs.FS.
func (e *embedded) Open(filename string) (fs.File, error) {
	return e.fs.Open(filename)
}

// OpenWithMetadata is like Open, except it also returns the file's metadata.
func (e *embedded) OpenWithMetadata(filename string) (fs.File, *Metadata, error) {
	fd, err := e.fs.Open(filename)
	if err != nil {
		return fd, nil, err
	}

	m, ok := db[filename]
	if !ok {
		return fd, nil, fmt.Errorf("%s: no metadata found", filename)
	}

	return fd, m, nil
}
