package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	env := readEnvFile()
	// Handle file uploads
	http.HandleFunc("/api/upload", uploadHandler)

	// Serve static files (Vue frontend)
	fs := http.FileServer(http.Dir(env["WEBAPP_DIR"]))
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
		formName := part.FormName()
		if formName == "path" {
			value, err := io.ReadAll(part)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			path = string(value)
			err = os.MkdirAll(path, 0777)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			filename := part.FileName()
			if lastFileName != filename {
				log.Println("Uploading file:", filename)
				lastFileName = filename
			}
			dst, err := os.Create(path + "/" + filename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = io.Copy(dst, part)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Println("Uploaded file " + filename)
		}

	}
}

func readEnvFile() map[string]string {
	bytes, err := os.ReadFile(".env")
	if err != nil {
		log.Fatal(err)
	}
	content := string(bytes)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")
	env := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			continue
		}
		splitLine := strings.Split(line, "=")
		if len(splitLine) != 2 {
			log.Fatalf("Invalid line in .env file: %s", line)
		}
		env[splitLine[0]] = splitLine[1]
	}
	return env
}
