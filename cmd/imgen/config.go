package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	OverflowPercent uint `json:"overflowPercent"`
}

func ReadConfigFile(path string) (*Config, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %v", path, err)
	}
	defer fd.Close()

	config := Config{}

	return &config, json.NewDecoder(fd).Decode(&config)
}
