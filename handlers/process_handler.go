package handlers

import (
	"encoding/json"
	"errors"
	"jellyfin_uploader/models"
	"jellyfin_uploader/repositories"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleProcess(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return ListProcessHandler(w, r)
	}
	if r.Method == http.MethodPost {
		return CreateProcessHandler(w, r)
	} else {
		return errors.New("Cannot handle this http method")
	}
}

func HandleProcessItem(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return ShowProcessHandler(w, r)
	}
	if r.Method == http.MethodDelete {
		return DeleteProcessHandler(w, r)
	} else {
		return errors.New("Can not handle method " + r.Method + " for this URL")
	}
}

func CreateProcessHandler(w http.ResponseWriter, r *http.Request) error {
	var uploadProcess models.UploadProcess
	body := r.Body
	if body == nil {
		return errors.New("Empty Body")
	}
	defer body.Close()
	err := json.NewDecoder(body).Decode(&uploadProcess)
	if err != nil {
		return err
	}
	err = repositories.CreateUploadProcess(&uploadProcess)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(uploadProcess)
	if err != nil {
		return err
	}
	return nil
}
func ShowProcessHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	uploadProcess, err := repositories.GetUploadProcess(id)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(uploadProcess)
	if err != nil {
		return err
	}
	return nil

}

func DeleteProcessHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	err := repositories.DeleteUploadProcess(id)
	if err != nil {
		return err
	}
	return nil
}
func ListProcessHandler(w http.ResponseWriter, r *http.Request) error {
	processes, err := repositories.ListUploadProcesses()
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(processes)
	if err != nil {
		return err
	}
	return nil
}
