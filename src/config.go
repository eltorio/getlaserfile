package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Paths []struct {
		RepoLocation string `yaml:"repolocation"`
		Url          string `yaml:"url"`
		Path         string `yaml:"path"`
	} `yaml:"paths"`
}

func ReadConfig(configYaml string) (Config, error) {
	var config Config

	data, err := os.ReadFile(configYaml)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
