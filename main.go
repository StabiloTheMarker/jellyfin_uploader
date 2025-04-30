package main

import (
	"jellyfin_uploader/handlers"
	"jellyfin_uploader/util"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	env := readEnvFile()

	fs := http.FileServer(http.Dir(env["WEBAPP_DIR"]))
	http.Handle("/", fs)

	http.HandleFunc("/api/upload", handlers.UploadHandler)
	http.HandleFunc("/api/upload_process", handlers.HandleProcess)

	port := env["PORT"]
	util.InitDB()
	log.Fatal(http.ListenAndServe(":"+port, nil))
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
