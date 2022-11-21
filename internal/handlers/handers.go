package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/beloslav13/servernotes/internal/database"
	"github.com/beloslav13/servernotes/internal/interfaces"
	"github.com/beloslav13/servernotes/internal/models"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var httpContext = context.Background()

type handler struct {
	logger logger.Logger
	db     *database.Storage
}

func NewHandler(log logger.Logger, db *database.Storage) interfaces.Handler {
	return &handler{
		logger: log,
		db:     db,
	}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc("/", h.HomeHandler)
	router.HandleFunc("/notes/", h.CreateNote).Methods("POST")
	router.HandleFunc("/notes/{id}/", h.GetNotes).Methods("GET")
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

func (h *handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("Handler create note...")
	vars := mux.Vars(r)
	h.logger.Infoln(vars)

	var note models.Note
	_ = json.NewDecoder(r.Body).Decode(&note)

	validate := validator.New()
	err := validate.Struct(note)
	if err != nil {
		h.logger.Errorln(err)
		http.Error(w, "failed to validate struct", http.StatusBadRequest)
		return
	}

	h.db.SaveNotes(httpContext, &note)

}
