package interfaces

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	Register(router *mux.Router)
	Get(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Validate(w http.ResponseWriter, m interface{}) bool
}

type Note interface {
	Handler
}

type Person interface {
	Handler
}
