package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Cache *Cache
}

type Cache struct {
	Redis struct {
		Address  string
		Password string
		DB       int
	}
}

func Read(path string) (*Config, error) {
	log.Printf("Reading config from %s", path)
	config := &Config{}
	if err := getFromJSONPath(path, config); err != nil {
		return nil, err
	}
	return config, nil
}

func getFromJSONPath(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	return json.NewDecoder(file).Decode(v)
}
