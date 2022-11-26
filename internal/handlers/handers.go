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
	router.HandleFunc("/notes/{id}/", h.GetNote).Methods("GET")
	http.Handle("/", router)
}

func (h *handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Warningln("HOME PAGE!!!!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Gorilla!\n"))
}

func (h *handler) GetNote(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("WOW NOTES!!!")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Notes: %v\n", vars["id"])
}

func (h *handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("Handler create note...")
	w.Header().Set("Content-Type", "application/json")

	var note models.Note
	_ = json.NewDecoder(r.Body).Decode(&note)

	// Если валидация структуры прошла успешно, создаём заметку в БД.
	if validateCreateNote(w, note, h) {
		return
	}

	if err := database.SaveNotes(httpContext, &note); err != nil {
		h.logger.Errorf("handler cannot save: %w\ndata: %w", err, note)
		response(w, "handler cannot save...", http.StatusBadRequest, err.Error())
		return
	}
	// Заметка создана в бд без ошибок, создаём json ответ
	response(w, "save ok", http.StatusCreated, nil)
}

// validateCreateNote validates structure completion
func validateCreateNote(w http.ResponseWriter, note models.Note, h *handler) bool {
	validate := validator.New()
	err := validate.Struct(note)
	if err != nil {
		h.logger.Errorln(err)
		response(w, "failed to validate struct", http.StatusBadRequest, err.Error())
		return true
	}
	return false
}

// response creates json response
func response(w http.ResponseWriter, msg string, status int, err interface{}) {
	resp := models.NewResponse(false, msg, err)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
