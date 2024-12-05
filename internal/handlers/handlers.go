package handlers

import (
	"auto-backend-trainee-assignment/internal/models"
	"auto-backend-trainee-assignment/internal/service"
	"encoding/json"
	"net/http"
)

type Handler interface {
	ShortenHandler(w http.ResponseWriter, r *http.Request)
	RedirectHandler(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service service.Service
}


func NewHandler(service service.Service) (Handler) {
	return &handler{service: service}
}


func (h *handler)ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"method not allowed",http.StatusMethodNotAllowed)
		return
	}
	var LongUrl models.ResponseUrl
	err := json.NewDecoder(r.Body).Decode(&LongUrl)
	if err != nil {

	}

	shortUrl,err := h.service.GenerateShortUrl(LongUrl)

	if err := json.NewEncoder(w).Encode(shortUrl); err != nil {
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
}

func (h *handler)RedirectHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w,"method not allowed",http.StatusMethodNotAllowed)
		return
	}
}