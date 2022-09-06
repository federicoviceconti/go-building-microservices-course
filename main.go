package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"services/handlers"
	"time"
)

func main() {
	l := log.New(os.Stdout, "products-api", log.LstdFlags)

	sm := http.NewServeMux()
	//sm.Handle("/", handlers.NewHello(l))
	//sm.Handle("/goodbye", handlers.NewGoodbye(l))
	sm.Handle("/", handlers.NewProductHandler(l))

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

/**
// Bind the address for hosting the server
// ----
// :9090 		  means: all the ip on the 9090 port
// 127.0.0.1:9090 means: specific ip on the 9090 port
// ----
// The second parameter indicate a ServeMux, in case is nil will be used
// a DefaultServeMux
err := http.ListenAndServe(":9090", sm)

// Register the path for the pattern provided. In this case we have the '/' path.
http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	handlers.NewHello(log).ServeHttp(writer, request)
})

// HandleFunc under the hood just map a request to a specific path, giving a response to
// client.
http.HandleFunc("/goodbye", func(writer http.ResponseWriter, request *http.Request) {
	println("Goodbye world!")
})
*/
