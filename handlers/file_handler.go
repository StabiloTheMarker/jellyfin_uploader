package handlers

import (
	"encoding/json"
	"errors"
	"jellyfin_uploader/repositories"
	"net/http"
)

func HandleFile(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return handleFileGet(w, r)
	} else {
		return errors.New("Cannot handle method " + r.Method + "for this URL")
	}
}

func handleFileGet(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	file, err := repositories.GetFile(id)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(file)
	if err != nil {
		return err
	}
	return nil
}
