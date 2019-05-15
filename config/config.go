package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`

	Mongo struct {
		URI        string `json:"uri"`
		DB         string `json:"db"`
		Collection string `json:"collection"`
	} `json:"mongo"`
}

func GetConfigFromJSON(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var cfg Config

	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, err
}
