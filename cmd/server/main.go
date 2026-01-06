package main

import (
	"log"
	"net/http"

	"github.com/shreeram-hegde/go-url-shortener/internal/handler"
	"github.com/shreeram-hegde/go-url-shortener/internal/service"
	"github.com/shreeram-hegde/go-url-shortener/internal/store"
)

func main() {

	st := store.NewMemoryStore()
	svc := service.NewShortenerService(st)
	h := handler.NewHandler(svc)

	_ = svc

	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", h.CreateShortURL)
	mux.HandleFunc("/", h.Redirect)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
