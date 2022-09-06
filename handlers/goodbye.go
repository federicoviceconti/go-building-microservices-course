package handlers

import (
	"log"
	"net/http"
)

type Goodbye struct {
	l *log.Logger
}

func (h *Goodbye) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte("Goodbye"))
	if err != nil {
		return
	}
}

func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}
