package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

type ServableFile struct {
	Path      string
	Extension string
	Base      string
	Name      string
	Pattern   string
	Content   []byte
}

//go:embed public/*
var embedFilesystem embed.FS

func main() {
	router := chi.NewRouter()
	files, err := getFiles(&embedFilesystem)

	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}

	servableFiles, err := getServableFiles(&embedFilesystem, files)

	if err != nil {
		log.Fatalf("failed to read files: %v", err)
	}

	for _, file := range servableFiles {
		log.Printf("add route: GET %s\n", file.Pattern)
		router.Get(file.Pattern, staticFileHandler(file))
	}

	log.Println("start web server on port 1337")

	if err := http.ListenAndServe(":1337", router); err != nil {
		log.Fatalf("failed to start web server: %v", err)
	}
}

func staticFileHandler(file *ServableFile) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if file.Extension == ".html" {
			w.Header().Add("Content-Type", "text/html")
		} else if file.Extension == ".css" {
			w.Header().Add("Content-Type", "text/css")
		} else {
			w.Header().Add("Content-Type", "text/plain")
		}

		w.Header().Add("Content-Length", fmt.Sprint(len(file.Content)))
		w.WriteHeader(http.StatusOK)
		w.Write(file.Content)
	}
}

func getServableFiles(efs *embed.FS, files []string) ([]*ServableFile, error) {
	var servableFiles []*ServableFile

	for _, file := range files {
		content, err := fs.ReadFile(efs, file)

		if err != nil {
			return servableFiles, err
		}

		extension := filepath.Ext(file)
		base := filepath.Base(file)
		name := strings.TrimSuffix(base, extension)
		pattern := fmt.Sprintf("/%s", base)

		if extension == ".html" {
			pattern = fmt.Sprintf("/%s", name)

			if name == "index" {
				pattern = "/"
			}
		}

		servableFiles = append(servableFiles, &ServableFile{
			Path:      file,
			Extension: extension,
			Base:      base,
			Name:      name,
			Pattern:   pattern,
			Content:   content,
		})
	}

	return servableFiles, nil
}

func getFiles(efs *embed.FS) ([]string, error) {
	var files []string

	err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
