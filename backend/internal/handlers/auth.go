package handlers

import (
	"log"
	"net/http"
)

type Login struct {
	l           *log.Logger
	verifyToken func(string) error
}

func NewLoginHandler(l *log.Logger, verifyToken func(string) error) *Login {
	return &Login{l, verifyToken}
}

func (l *Login) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

}
