package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Home struct {
	l *log.Logger
}

func NewHome(l *log.Logger) *Home {
	return &Home{l}
}

func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Home Page")
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "oops", http.StatusBadRequest)
		return
	}

	h.l.Printf("data-> %s", data)
	w.Write(data)
}
