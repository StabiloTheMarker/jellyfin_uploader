package handlers

import (
	"errors"
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
	id := r.PathValue("id")
	uploadProcess, err := repositories.GetUploadProcess(id)
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
		filename := part.FileName()
		filepath := filepath.Join(uploadProcess.DirPath, filename)
		log.Println("Uploading file:", filename)
		dst, err := os.Create(filepath)
		if err != nil {
			return err
		}
		file := models.File{
			Name:            dst.Name(),
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
		err = dst.Close()
		if err != nil {
			return err
		}
		log.Println("Uploaded file " + filename)
		err = part.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
