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
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func main() {
	listenPort := flag.String("port", "80", "The listening port")
	configYaml := flag.String("config", "", "The location of the config file")

	flag.Parse()

	config, err := ReadConfig(*configYaml)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Check if repoLocation and commitHash are provided
	if *configYaml == "" || !isInteger(*listenPort) {
		fmt.Println("Error: Not enough arguments. Usage: ./getlaserfile --config=<config.yaml> --port=<integer>")
		fmt.Println("\t config.yaml sample:")
		fmt.Println(`paths:
  - repolocation: "/repo1"
    url: "/url1"
    path: "/path1"
  - repolocation: "/repo2"
    url: "/url2"
    path: "/path2"`)
		return
	}

	for _, path := range config.Paths {
		if path.Path == "" || path.RepoLocation == "" || path.Url == "" {
			log.Fatalf("error: malformed config")
		} else {
			log.Printf("info: serve %s corresponding to repo %s and file %s", path.Url, path.RepoLocation, path.Path)
			url := path.Url
			repo := path.RepoLocation
			file := path.Path
			http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				log.Printf("info: request %s corresponding to repo %s and file %s", url, repo, file)
				HandleBinary(w, r, repo, file)
			})
		}

	}

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	http.ListenAndServe(":"+*listenPort, nil)
}
