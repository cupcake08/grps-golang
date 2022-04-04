package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cupcake08/grps-golang/handlers"
)

func main() {
	l := log.New(os.Stdout, "grps-go", log.LstdFlags)
	//make a handler
	hh := handlers.NewHome(l)
	gh := handlers.NewGoodBye(l)
	ph := handlers.NewProducts(l)

	//register that handler with server
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/bye", gh)
	sm.Handle("/products/", ph)

	//creating our custom server
	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  100 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 1 * time.Second,
		ErrorLog:     l,
	}
	go func() {
		l.Println("Starting server at Port 8080")
		if err := s.ListenAndServe(); err != nil {
			l.Fatal(err)
		}
	}()
	//os signal
	schan := make(chan os.Signal)
	signal.Notify(schan, os.Interrupt)
	signal.Notify(schan, os.Kill)

	sgn := <-schan
	l.Println("Received terminate, graceful shutdown", sgn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}
