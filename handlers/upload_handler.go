package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"io"
	"jellyfin_uploader/models"
	"jellyfin_uploader/repositories"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return errors.New("Invalid Request Method")
	}
	reader, err := r.MultipartReader()
	if err != nil {
		return err
	}
	if reader == nil {
		return errors.New("multipart reader is nil")
	}
	vars := mux.Vars(r)
	id := vars["id"]
	uploadProcess, err := repositories.GetUploadProcess(id)
	os.MkdirAll(uploadProcess.DirPath, 0777)
	if err != nil {
		return err
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		defer part.Close()
		filename := part.FileName()
		filepath := filepath.Join(uploadProcess.DirPath, filename)
		log.Println("Uploading file:", filename)
		dst, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer dst.Close()
		file := models.File{
			Name:            filename,
			Uploaded:        false,
			UploadProcessID: uploadProcess.ID,
		}
		err = repositories.CreateFile(&file)
		if err != nil {
			return err
		}
		_, err = io.Copy(dst, part)
		if err != nil {
			return err
		}
		file.Uploaded = true
		file.UploadedAt = time.Now()
		err = repositories.UpdateFile(&file)
		if err != nil {
			return err
		}
		log.Println("Uploaded file " + filename)
	}
	return nil
}
