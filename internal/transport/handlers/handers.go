package handlers

import (
	"fmt"
	"net/http"

	"github.com/beloslav13/servernotes/internal/transport/interfaces"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/gorilla/mux"
)

type handler struct {
	logger logger.Logger
}

func NewHandler(log logger.Logger) interfaces.Handler {
	return &handler{
		logger: log,
	}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc("/", h.HomeHandler).Methods("GET")
	router.HandleFunc("/notes/{id}/", h.GetNotes)
	http.Handle("/", router)
}

func (h *handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Warningln("HOME PAGE!!!!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Gorilla!\n"))
}

func (h *handler) GetNotes(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("WOW NOTES!!!")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Notes: %v\n", vars["id"])
}
