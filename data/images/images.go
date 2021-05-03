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
	"dents_du_midi_1.jpg": newMetadata(time.January,
		2020,
		"Dents du Midi, Switzerland",
		"TODO"),
	"dubai_1.jpg": newMetadata(time.June,
		2017,
		"Dubai, UAE",
		"TODO"),
	"fuji_1.jpg": newMetadata(time.October,
		2017,
		"Mount Fuji, Japan",
		"TODO"),
	"geneva_1.jpg": newMetadata(time.June,
		2016,
		"Geneva, Switzerland",
		"TODO"),
	"kyoto_1.jpg": newMetadata(time.October,
		2017,
		"Kyoto, Japan",
		"TODO"),
	"lhc_1.jpg": newMetadata(time.August,
		2019,
		"LHC, France / Switzerland",
		"TODO"),
	"malibu_1.jpg": newMetadata(time.March,
		2019,
		"Malibu, USA",
		"TODO"),
	"montreux_1.jpg": newMetadata(time.October,
		2016,
		"Montreux, Switzerland",
		"TODO"),
	"new_delhi_1.jpg": newMetadata(time.June,
		2017,
		"New Delhi, India",
		"TODO"),
	"newyork_2.jpg": newMetadata(time.August,
		2015,
		"New York, USA",
		"TODO"),
	"nuggets_point_1.jpg": newMetadata(time.January,
		2019,
		"Nuggets Point, New Zealand",
		"TODO"),
	"shenzhen_1.jpg": newMetadata(time.August,
		2014,
		"Shenzhen, China",
		"TODO"),
	"singapore_1.jpg": newMetadata(time.January,
		2019,
		"Singapore",
		"TODO"),
	"thun_1.jpg": newMetadata(time.May,
		2016,
		"Thun, Switzerland",
		"TODO"),
	"whaikiti_beach_1.jpg": newMetadata(time.January,
		2019,
		"Whaikiti Beach, New Zealand",
		"TODO"),
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
