package main

import (
	"encoding/json"
	"io"
	"os"
)

type ImageMetadata struct {
	Date      string   `json:"date"`
	Formats   []string `json:"formats"`
	Heights   []int    `json:"heights"`
	Location  string   `json:"location"`
	MainColor string   `json:"main_color"`
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
