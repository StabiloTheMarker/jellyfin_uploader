package main

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	uploadPath = "./uploads"
	staticPath = "./webapp/dist"
	maxMemory  = 1024 * 1024 * 100 // 100 MB
)

func main() {
	// Ensure upload directory exists
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Handle file uploads
	http.HandleFunc("/upload", uploadHandler)

	// Serve static files (Vue frontend)
	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Limit memory usage but allow streaming large files
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		http.Error(w, "Could not parse multipart form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	for _, fileHeader := range files {
		if err := saveUploadedFile(fileHeader); err != nil {
			http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Upload Successful")
}

func saveUploadedFile(fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	dst, err := os.Create(filepath.Join(uploadPath, sanitizeFilename(fileHeader.Filename)))
	if err != nil {
		return err
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(dst)

	_, err = io.Copy(dst, file)
	return err
}

func sanitizeFilename(name string) string {
	// Simple sanitization to avoid path traversal
	return filepath.Base(strings.ReplaceAll(name, "..", ""))
}
