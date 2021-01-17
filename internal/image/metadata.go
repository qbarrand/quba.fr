//go:generate go run github.com/golang/mock/mockgen -source metadata.go -destination mock_image/metadata.go MetaDB

package image

import (
	"errors"
	"fmt"
	"time"
)

var ErrNoSuchMetadata = errors.New("no such metadata")

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

type MetaDB interface {
	AllNames() ([]string, error)
	GetMetadata(string) (*Metadata, error)
}

type StaticMetaDB struct {
	images map[string]*Metadata
}

func NewStaticMetaDB() *StaticMetaDB {
	return &StaticMetaDB{
		images: map[string]*Metadata{
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
		},
	}
}

func (s *StaticMetaDB) AllNames() ([]string, error) {
	names := make([]string, 0, len(s.images))

	for n, _ := range s.images {
		names = append(names, n)
	}

	return names, nil
}

func (s *StaticMetaDB) GetMetadata(name string) (*Metadata, error) {
	meta := s.images[name]

	if meta == nil {
		return nil, fmt.Errorf("%s: %w", name, ErrNoSuchMetadata)
	}

	return meta, nil
}
