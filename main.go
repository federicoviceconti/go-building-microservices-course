package main

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"services/handlers"
	"time"
)

func main() {
	l := log.New(os.Stdout, "products-api", log.LstdFlags)

	ph := handlers.NewProductHandler(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductsValidation)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductsValidation)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	rOpt := middleware.RedocOpts{
		SpecURL: "/swagger.yaml",
	}
	dh := middleware.Redoc(rOpt, nil)

	getRouter.Handle("/docs", dh)

	// We're adding a new handler for the swagger.yaml file looking on the './' directory
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// Let's create a server, with custom parameters
	s := http.Server{
		Addr:         ":9090",
		Handler:      sm,
		ErrorLog:     l,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()

		if err != nil {
			l.Fatal(err)
		}
	}()

	// Create an instance of Signal, that capture an OS signal
	// Type is a channel, where you can send and receives value
	sigChan := make(chan os.Signal)

	// These method will be called when a sigterm or interrupt is received
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Receive from ch and assign value to v
	sig := <-sigChan
	l.Println("Received terminate gracefully", sig)

	// Waiting 30 seconds to complete all the operations when the server is going to shutting down
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// When all requests and connection are completed, the server shutdown gracefully
	err := s.Shutdown(tc)
	if err != nil {
		return
	}
}
