package rest

import (
	"github.com/gorilla/mux"
)

func Route(r *mux.Router) {
	r.HandleFunc("/content/get/{id}", GetContent)
}
