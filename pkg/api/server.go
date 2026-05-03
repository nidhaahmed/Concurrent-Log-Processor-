package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func StartServer() {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/filter", filterHandler)
	mux.HandleFunc("/logs", getAllLogsHandler)
	mux.HandleFunc("/metrics", metricsHandler)
	mux.HandleFunc("/compare", compareHandler)

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	handler := LoggingMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		fmt.Println("Server running on " + ":" + port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}
