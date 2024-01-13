package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

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
	repoLocation := flag.String("repolocation", "", "The location of the repository")
	ihmLocation := flag.String("ihmlocation", "builds/IHM/ihm.exe", "The location of the ihm binay")
	startupLocation := flag.String("startuplocation", "builds/sbRIO-9651/home/lvuser/natinst/bin/startup.rtexe", "The location of the startup.rtexe binary")

	flag.Parse()

	// Check if repoLocation and commitHash are provided
	if *repoLocation == "" || !isInteger(*listenPort) {
		fmt.Println("Error: Not enough arguments. Usage: ./getlaserfile --repolocation=<repoLocation> --port=<integer>")
		return
	}

	http.HandleFunc("/ihm.exe", func(w http.ResponseWriter, r *http.Request) {
		handleBinary(w, r, repoLocation, ihmLocation)
	})

	http.HandleFunc("/startup.rtexe", func(w http.ResponseWriter, r *http.Request) {
		handleBinary(w, r, repoLocation, startupLocation)
	})

	http.ListenAndServe(":"+*listenPort, nil)
}
