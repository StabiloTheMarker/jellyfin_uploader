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

	http.HandleFunc("/api/upload/{id}", util.MakeApiFunc(handlers.HandleUpload))
	http.HandleFunc("/api/upload_process", util.MakeApiFunc(handlers.HandleProcess))
	http.HandleFunc("/api/upload_process/{id}", util.MakeApiFunc(handlers.HandleProcessItem))
	http.HandleFunc("/api/file", util.MakeApiFunc(handlers.HandleFile))

	port := env["PORT"]
	util.InitDB()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
