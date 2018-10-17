package handlers

import (
	"go_SkyDrive/handlers/index"
	"net/http"
)

func MyUrls() {
	http.HandleFunc("/", index.Upload_file_v1)
}
