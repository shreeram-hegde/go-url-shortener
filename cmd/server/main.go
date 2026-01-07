package main

import (
	"log"
	"net/http"
	"time"

	"github.com/shreeram-hegde/go-url-shortener/internal/handler"
	"github.com/shreeram-hegde/go-url-shortener/internal/service"
	"github.com/shreeram-hegde/go-url-shortener/internal/store"
)

func main() {

	st := store.NewMemoryStore()
	svc := service.NewShortenerService(st)
	h := handler.NewHandler(svc)

	_ = svc

	//Calling cleanup routine

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			err := st.DeleteExpired(time.Now())
			if err != nil {
				log.Println("clean up error:", err)
			}
		}
	}() //running an anonyamous function

	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", h.CreateShortURL)
	mux.HandleFunc("/", h.Redirect)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
