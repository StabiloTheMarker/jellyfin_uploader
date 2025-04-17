package main

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const (
	staticPath = "./webapp/dist"
)

func main() {

	env := readEnvFile()
	// Handle file uploads
	http.HandleFunc("/api/upload", uploadHandler)

	// Serve static files (Vue frontend)
	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/", fs)

	port := env["PORT"]
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if reader == nil {
		http.Error(w, "multipart reader is nil", http.StatusInternalServerError)
		return
	}

	var lastFileName string
	var fileParts []*multipart.Part
	var path string
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if part == nil {
			http.Error(w, "part is nil", http.StatusInternalServerError)
			return
		}
		formName := part.FormName()
		if formName == "path" {
			log.Println("im here")
			value, err := io.ReadAll(part)
			log.Printf("value is %s", value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			path = string(value)
			log.Printf("path is %s", path)
			err = os.MkdirAll(path, 0777)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			fileParts = append(fileParts, part)
		}

	}
	for _, filePart := range fileParts {
		filename := filePart.FileName()
		if lastFileName != filename {
			log.Println("Uploading file:", filename)
			lastFileName = filename
		}
		dst, err := os.Create(path + "/" + filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(dst, filePart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func readEnvFile() map[string]string {
	bytes, err := os.ReadFile(".env")
	if err != nil {
		log.Fatal(err)
	}
	content := string(bytes)
	lines := strings.Split(content, "\n")
	env := make(map[string]string)
	for _, line := range lines {
		splitLine := strings.Split(line, "=")
		if len(splitLine) != 2 {
			log.Fatalf("Invalid line in .env file: %s", line)
		}
		env[splitLine[0]] = splitLine[1]
	}
	return env
}
