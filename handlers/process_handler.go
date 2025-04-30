package handlers

import (
	"encoding/json"
	"jellyfin_uploader/repositories"
	"net/http"
)

func HandleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ListProcessHandler(w, r)
	} else {
		DeleteProcessHandler(w, r)
	}
}

func DeleteProcessHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("processId")
	err := repositories.DeleteUploadProcess(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func ListProcessHandler(w http.ResponseWriter, r *http.Request) {
	processes, err := repositories.ListUploadProcesses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(processes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
