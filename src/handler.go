package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

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

func HandleBinary(w http.ResponseWriter, r *http.Request, repoLocation, path string) {
	hash := r.URL.Query().Get("hash")

	// Check if hash is a valid hash
	match, _ := regexp.MatchString("^[a-f0-9]{40}$", hash)
	if !match {
		http.Error(w, "Invalid commit hash. It should contain only numbers and characters from a to f and be exactly 40 characters long.", http.StatusBadRequest)
		log.Printf("Error: Invalid commit hash %s. It should contain only numbers and characters from a to f and be exactly 40 characters long.", hash)
		return
	}

	content, err := GetFileAtCommit(repoLocation, hash, path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	// Read the blob content
	reader, err := content.Reader()
	if err != nil {
		http.Error(w, "Error reading blob", http.StatusInternalServerError)
		log.Printf("Error: reading blob %v", err)
		return
	}
	log.Printf("info: using hash %s", hash)
	defer reader.Close()

	// Copy the blob content to the response body
	if _, err := io.Copy(w, reader); err != nil {
		http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
	}
}
