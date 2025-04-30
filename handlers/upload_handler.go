package handlers

import (
	"io"
	"jellyfin_uploader/models"
	"jellyfin_uploader/repositories"
	"log"
	"net/http"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
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
	uploadProcess := models.UploadProcess{}
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
			err = repositories.CreateUploadProcess(&uploadProcess)
			log.Println("Created upload process")
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
			log.Printf("Created UploadProcess with ID: %d", uploadProcess.ID)
			file := models.File{Filepath: dst.Name(), Uploaded: false, UploadProcessID: uploadProcess.ID}
			err = repositories.CreateFile(&file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_, err = io.Copy(dst, part)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			file.Uploaded = true
			file.UploadedAt = time.Now()
			err = repositories.UpdateFile(&file)
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
