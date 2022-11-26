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
}

func NewHandler(log logger.Logger) interfaces.Handler {
	return &handler{
		logger: log,
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
	w.Header().Set("Content-Type", "application/json")

	resp := models.Response{}
	var note models.Note
	_ = json.NewDecoder(r.Body).Decode(&note)

	validate := validator.New()
	err := validate.Struct(note)
	if err != nil {
		h.logger.Errorln(err)
		resp = models.Response{Result: false, Message: "failed to validate struct", Err: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := database.SaveNotes(httpContext, &note); err != nil {
		h.logger.Errorf("handler cannot save: %w\ndata: %w", err, note)
		resp = models.Response{Result: false, Message: "handler cannot save...", Err: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp = models.Response{Result: true, Message: "save ok", Err: nil}
	json.NewEncoder(w).Encode(resp)
}
