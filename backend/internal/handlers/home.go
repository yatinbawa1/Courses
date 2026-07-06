package handlers

import (
	"log"
	"net/http"
)

type Home struct {
	l *log.Logger
}

func NewHomeHandler(l *log.Logger) *Home {
	return &Home{l}
}

func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello! Server is Running"))
}
