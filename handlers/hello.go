package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func (h *Hello) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	body, _ := io.ReadAll(request.Body)

	result, err := fmt.Fprintf(writer, "Hello %s", string(body))

	if err != nil {
		http.Error(writer, "Oooops!", http.StatusBadRequest)
		return
	}

	h.l.Println("Hello", result)
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}
