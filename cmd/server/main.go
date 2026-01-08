package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/shreeram-hegde/go-url-shortener/internal/handler"
	"github.com/shreeram-hegde/go-url-shortener/internal/service"
	"github.com/shreeram-hegde/go-url-shortener/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	//getting port number from .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	st, cleanup, err := createStore()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	svc := service.NewShortenerService(st)

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + port
	}

	h := handler.NewHandler(svc, baseURL)

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

	//health for hosting
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Starting server on :", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))

}

func createStore() (store.Store, func(), error) {
	storeType := os.Getenv("STORE")
	if storeType == "" {
		storeType = "memory" //the default is memory version
	}

	switch storeType {
	case "memory":
		st := store.NewMemoryStore()
		fmt.Println("Starting in memory DB")
		return st, func() {}, nil

	case "sqlite":
		path := os.Getenv("SQLITE_PATH")
		if path == "" {
			path = "data.db"
		}

		st, err := store.NewSQLiteStore(path)
		if err != nil {
			return nil, nil, err
		}

		//cleanup
		fmt.Println(`Starting in SQLite DB with the path`, path)
		return st, func() {
			st.Close()
		}, nil

	default:
		return nil, nil, fmt.Errorf("unknown STORE type: &s", storeType)

	}
}
