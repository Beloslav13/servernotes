package interfaces

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	Register(router *mux.Router)
	HomeHandler(w http.ResponseWriter, r *http.Request)
	GetNotes(w http.ResponseWriter, r *http.Request)
}
