// Package classification of Product API
//
// Documentation for Product API
//
// Schemas: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/(text|json)
//
// swagger:meta

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	prohandlers "github.com/cupcake08/grps-golang/handlers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "grps-go ", log.LstdFlags)
	//make a handler
	// hh := handlers.NewHome(l)
	//gh := handlers.NewGoodBye(l)
	ph := prohandlers.NewProducts(l)

	//register that handler with server
	r := mux.NewRouter()
	router := r.PathPrefix("/products").Subrouter()
	router.Headers("Content-Type", "application/(text|json)")

	//Get methods Router
	getRoute := router.Methods(http.MethodGet).Subrouter()
	//Post methods router
	postRoute := router.Methods(http.MethodPost).Subrouter()
	postRoute.Use(ph.Middleware)
	//Put methods router
	putRoute := router.Methods(http.MethodPut).Subrouter()
	putRoute.Use(ph.Middleware)
	//Delete methods router
	deleteRoute := router.Methods(http.MethodDelete).Subrouter()

	//Get Methods
	getRoute.HandleFunc("", ph.GetProducts)

	//post Methods
	postRoute.HandleFunc("", ph.AddProduct)

	//put methods
	putRoute.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)

	//delete Methods
	deleteRoute.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	//creating cors handler
	ch := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))
	//creating our custom server
	s := &http.Server{
		Addr:         ":8080",
		Handler:      ch(r),
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
	schan := make(chan os.Signal, 1)
	signal.Notify(schan, os.Interrupt)

	<-schan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Shutdown(ctx)
	log.Println("Shuting down")
	os.Exit(0)
}
