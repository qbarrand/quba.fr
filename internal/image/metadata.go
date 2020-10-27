//go:generate go run github.com/golang/mock/mockgen -source metadata.go -destination mock_image/metadata.go MetaDB

package image

import (
	"errors"
	"fmt"
	"time"
)

var ErrNoSuchMetadata = errors.New("no such metadata")

type Metadata struct {
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
}

func newMetadata(month time.Month, year int, location string) *Metadata {
	return &Metadata{
		Date:     time.Date(year, month, 0, 0, 0, 0, 0, time.UTC),
		Location: location,
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
			"dents_du_midi_1.jpg":  newMetadata(time.January, 2020, "Dents du Midi, Switzerland"),
			"dubai_1.jpg":          newMetadata(time.June, 2017, "Dubai, UAE"),
			"fuji_1.jpg":           newMetadata(time.October, 2017, "Mount Fuji, Japan"),
			"geneva_1.jpg":         newMetadata(time.June, 2016, "Geneva, Switzerland"),
			"kyoto_1.jpg":          newMetadata(time.October, 2017, "Kyoto, Japan"),
			"lhc_1.jpg":            newMetadata(time.August, 2019, "LHC, France / Switzerland"),
			"malibu_1.jpg":         newMetadata(time.March, 2019, "Malibu, USA"),
			"montreux_1.jpg":       newMetadata(time.October, 2016, "Montreux, Switzerland"),
			"new_delhi_1.jpg":      newMetadata(time.June, 2017, "New Delhi, India"),
			"newyork_2.jpg":        newMetadata(time.August, 2015, "New York, USA"),
			"nuggets_point_1.jpg":  newMetadata(time.January, 2019, "Nuggets Point, New Zealand"),
			"shenzhen_1.jpg":       newMetadata(time.August, 2014, "Shenzhen, China"),
			"singapore_1.jpg":      newMetadata(time.January, 2019, "Singapore"),
			"thun_1.jpg":           newMetadata(time.May, 2016, "Thun, Switzerland"),
			"whaikiti_beach_1.jpg": newMetadata(time.January, 2019, "Whaikiti Beach, New Zealand"),
		},
	}
}

func (smdb *StaticMetaDB) AllNames() ([]string, error) {
	names := make([]string, 0, len(smdb.images))

	for n, _ := range smdb.images {
		names = append(names, n)
	}

	return names, nil
}

func (smdb *StaticMetaDB) GetMetadata(name string) (*Metadata, error) {
	meta := smdb.images[name]

	if meta == nil {
		return nil, fmt.Errorf("%s: %w", name, ErrNoSuchMetadata)
	}

	return meta, nil
}
