package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/beloslav13/servernotes/internal/database"
	"github.com/beloslav13/servernotes/internal/interfaces"
	"github.com/beloslav13/servernotes/internal/models"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
	router.HandleFunc("/notes/{id:[0-9]+}/", h.GetNote).Methods("GET")
	router.HandleFunc("/notes/{id:[0-9]+}/", h.DeleteNote).Methods("DELETE")
	http.Handle("/", router)
}

func (h *handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Warningln("HOME PAGE!!!!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Gorilla!\n"))
}

func (h *handler) GetNote(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("Handler GetNote")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok {
		h.logger.Errorf("id does not exit in vars: %v", vars)
		return
	}

	id, _ := strconv.Atoi(sid)
	note, err := database.GetNote(httpContext, id)
	if err != nil {
		h.logger.Errorf("handler GetNote error get: %s. id: %d", err.Error(), id)
		response(w, fmt.Sprintf("cannot get note with id: %d", id), http.StatusBadRequest, err.Error(), nil)
		return
	}
	//h.logger.Infof("%+v", note)
	response(w, "get ok", http.StatusOK, nil, note)
}

func (h *handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("Handler CreateNote")
	w.Header().Set("Content-Type", "application/json")

	// TODO: при создании заметки, id создаётся бд, поэтому в ответе json id значение по умолчанию. фикс.
	var note models.Note
	_ = json.NewDecoder(r.Body).Decode(&note)

	// Если валидация структуры прошла успешно, создаём заметку в БД.
	if validateCreateNote(w, note, h) {
		return
	}

	if err := database.SaveNote(httpContext, &note); err != nil {
		h.logger.Errorf("handler cannot save: %w\ndata: %w", err, note)
		response(w, "handler cannot save...", http.StatusBadRequest, err.Error(), nil)
		return
	}
	// Заметка создана в бд без ошибок, создаём json ответ
	response(w, "save ok", http.StatusCreated, nil, note)
}

func (h *handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("Handler DeleteNote")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok {
		h.logger.Errorf("id does not exit in vars: %v", vars)
		return
	}

	id, _ := strconv.Atoi(sid)
	err := database.DeleteNote(httpContext, id)
	if err != nil {
		h.logger.Errorf("handler DeleteNote error delete: %s. id: %d", err.Error(), id)
		response(w, fmt.Sprintf("cannot delete note with id: %d", id), http.StatusBadRequest, err.Error(), nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

// validateCreateNote validates structure completion
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

// response creates json response
func response(w http.ResponseWriter, msg string, status int, err, obj interface{}) {
	var res bool
	if err == nil {
		res = true
	}
	resp := models.NewResponse(res, msg, err, obj)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
