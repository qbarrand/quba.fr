package metadata

import (
	"encoding/json"
	"io"
	"os"
	"sort"
)

type ImageMetadata struct {
	Date      string   `json:"date"`
	Formats   []string `json:"formats"`
	Heights   []int    `json:"heights"`
	Location  string   `json:"location"`
	MainColor string   `json:"main_color"`
	Widths    []int    `json:widths`
}

func NewMetadata(date, location, mainColor string) *ImageMetadata {
	return &ImageMetadata{
		Date:      date,
		Location:  location,
		MainColor: mainColor,
	}
}

func (m *ImageMetadata) AddFormat(f string) {
	formats := append(m.Formats, f)
	sort.Strings(formats)

	m.Formats = formats
}

func (m *ImageMetadata) AddHeight(h int) {
	heights := append(m.Heights, h)
	sort.Ints(heights)

	m.Heights = heights
}

func (m *ImageMetadata) AddWidths(h int) {
	widths := append(m.Widths, h)
	sort.Ints(widths)

	m.Widths = widths
}

type Metadata map[string]*ImageMetadata

func Read(r io.Reader) (Metadata, error) {
	m := make(Metadata)

	return m, json.NewDecoder(r).Decode(&m)
}

func ReadFromFile(name string) (Metadata, error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	return Read(fd)
}

func (m Metadata) Write(w io.Writer) error {
	return json.NewEncoder(w).Encode(m)
}

func (m Metadata) WriteToFile(name string) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return m.Write(fd)
}
