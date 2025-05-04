package main

import (
	"jellyfin_uploader/handlers"
	"jellyfin_uploader/util"
	"log"
	"net/http"
)

func main() {

	env := util.ReadEnvFile()

	fs := http.FileServer(http.Dir(env["WEBAPP_DIR"]))
	http.Handle("/", fs)

	http.HandleFunc("/api/upload", handlers.UploadHandler)
	http.HandleFunc("/api/upload_process", handlers.HandleProcess)

	port := env["PORT"]
	util.InitDB()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
