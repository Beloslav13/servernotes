package handlers

import (
	"context"
	"encoding/json"
	"github.com/beloslav13/servernotes/internal/interfaces"
	"github.com/beloslav13/servernotes/internal/models"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

var httpContext = context.Background()

type handler struct {
	logger logger.Logger
}

func NewHandler(log logger.Logger) interfaces.Handler {
	return &handler{
		logger: log,
	}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc("/", h.HomeHandler)
	router.HandleFunc("/notes/", h.CreateNote).Methods("POST")
	router.HandleFunc("/notes/", h.GetAllNotes).Methods("GET")
	router.HandleFunc("/notes/{id:[0-9]+}/", h.GetNote).Methods("GET")
	router.HandleFunc("/notes/{id:[0-9]+}/", h.DeleteNote).Methods("DELETE")
	http.Handle("/", router)
}

func (h *handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Warningln("HOME PAGE!!!!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Gorilla!\n"))
}

// validateCreateNote validates structure completion.
func validateCreateNote(w http.ResponseWriter, note models.Note, h *handler) bool {
	validate := validator.New()
	err := validate.Struct(note)
	if err != nil {
		h.logger.Errorln(err)
		response(w, "failed to validate struct", http.StatusBadRequest, err.Error(), nil)
		return true
	}
	return false
}

// response creates json response.
func response(w http.ResponseWriter, msg string, status int, err, obj interface{}) {
	var res bool
	if err == nil {
		res = true
	}
	resp := models.NewResponse(res, msg, err, obj)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
