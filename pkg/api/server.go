package api

import (
	"fmt"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/filter", filterHandler)
	http.HandleFunc("/logs", getAllLogsHandler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}