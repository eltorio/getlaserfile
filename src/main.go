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
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"gopkg.in/yaml.v2"
)

type Path struct {
	RepoLocation string `yaml:"repolocation"`
	Url          string `yaml:"url"`
	Path         string `yaml:"path"`
}

type Config struct {
	Paths []Path `yaml:"paths"`
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func GetFileAtCommit(repoPath string, commitHash string, filePath string) (object.Blob, error) {
	// Open the repository
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return object.Blob{}, err
	}

	// Convert the commit hash into a plumbing.Hash
	hash := plumbing.NewHash(commitHash)

	// Retrieve the commit
	c, err := r.CommitObject(hash)
	if err != nil {
		return object.Blob{}, err
	}

	// Get the root tree of the commit
	t, err := c.Tree()
	if err != nil {
		return object.Blob{}, err
	}

	// Find the desired file in the tree
	f, err := t.File(filePath)
	if err != nil {
		return object.Blob{}, err
	}

	return f.Blob, nil
}

func handleBinary(w http.ResponseWriter, r *http.Request, repoLocation *string, path *string) {
	hash := r.URL.Query().Get("hash")

	// Check if hash is a valid hash
	match, _ := regexp.MatchString("^[a-f0-9]{40}$", hash)
	if !match {
		http.Error(w, "Invalid commit hash. It should contain only numbers and characters from a to f and be exactly 40 characters long.", http.StatusBadRequest)
		return
	}

	content, err := GetFileAtCommit(*repoLocation, hash, *path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	// Read the blob content
	reader, err := content.Reader()
	if err != nil {
		http.Error(w, "Error reading blob", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	// Copy the blob content to the response body
	if _, err := io.Copy(w, reader); err != nil {
		http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
	}
}

func main() {
	listenPort := flag.String("port", "80", "The listening port")
	configYaml := flag.String("config", "", "The location of the config file")

	flag.Parse()

	// Read config file
	data, err := os.ReadFile(*configYaml)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unmarshal YAML data into Config struct
	var config Config
	err = yaml.Unmarshal(data, &config)
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
			http.HandleFunc(path.Url, func(w http.ResponseWriter, r *http.Request) {
				handleBinary(w, r, &path.RepoLocation, &path.Path)
			})
		}

	}

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	http.ListenAndServe(":"+*listenPort, nil)
}
