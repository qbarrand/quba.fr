package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ImageMetadata struct {
	Date     string   `json:"date"`
	Formats  []string `json:"formats"`
	Heights  []int    `json:"heights"`
	Location string   `json:"location"`
}

type Metadata map[string]*ImageMetadata

func Read(r io.Reader) (Metadata, error) {
	m := make(Metadata)

	return m, json.NewDecoder(r).Decode(&m)
}

func ReadFromFile(name string) (m Metadata, err error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := fd.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("could not close the file: %v", cerr)
		}
	}()

	return Read(fd)
}
