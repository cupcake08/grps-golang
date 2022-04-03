package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cupcake08/grps-golang/handlers"
)

func main() {
	l := log.New(os.Stdout, "grps-go", log.LstdFlags)
	//make a handler
	hh := handlers.NewHome(l)
	gh := handlers.NewGoodBye(l)

	//register that handler with server
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/bye", gh)
	http.ListenAndServe(":8080", sm)
}
