package main

import (
	"jellyfin_uploader/handlers"
	"jellyfin_uploader/util"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	env := util.ReadEnvFile()
	fs := http.FileServer(http.Dir(env["WEBAPP_DIR"]))
	webDir := env["WEBAPP_DIR"]
	router := mux.NewRouter()
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/", fs))

	// --- Serve index.html for all non-API, non-static routes (SPA fallback)

	router.HandleFunc("/api/upload/{id}", util.MakeApiFunc(handlers.HandleUpload)).Methods("POST")
	router.HandleFunc("/api/upload_process", util.MakeApiFunc(handlers.HandleProcess)).Methods("GET", "POST", "DELETE")
	router.HandleFunc("/api/upload_process/{id}", util.MakeApiFunc(handlers.HandleProcessItem)).Methods("GET", "DELTE")
	router.HandleFunc("/api/file", util.MakeApiFunc(handlers.HandleFile)).Methods("GET", "POST", "DELETE")

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, webDir+"/index.html")
	})
	port := env["PORT"]
	util.InitDB()
	log.Fatal(http.ListenAndServe(":"+port, router))
}
