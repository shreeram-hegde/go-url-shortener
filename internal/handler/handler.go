package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shreeram-hegde/go-url-shortener/internal/service"
)

type Handler struct {
	svc *service.ShortenerService
}

func NewHandler(svc *service.ShortenerService) *Handler {
	return &Handler{svc: svc}
}

type shortenRequest struct { //I guess this is the request format, a post request should match this
	URL           string `json:"url"`
	ExpiryMinutes int    `json:"expiry_minutes"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	//If Not Post and somehow this handler or path has been called catch it
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Check if the request is proper
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	//if Expiry in mins is too short make it 60 mins
	if req.ExpiryMinutes <= 0 {
		req.ExpiryMinutes = 60
	}

	//Create a new shorturl
	u, err := h.svc.Create(req.URL, time.Duration(req.ExpiryMinutes)*time.Minute)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Create the response structure
	resp := shortenResponse{
		ShortURL: "http://localhost:8080/" + u.Code,
	}

	//Add header stuff and encode the response to json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:] // remove "/"

	if code == "" {
		http.NotFound(w, r)
		return
	}

	u, err := h.svc.Resolve(code)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, u.LongURL, http.StatusFound)
}
