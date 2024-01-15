/*
 * getlaserfile
 * Copyright (C) 2024 Ronan LE MEILLAT
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 * Developed for ISMO Group (https://www.ismo-group.co.uk)
 */

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
