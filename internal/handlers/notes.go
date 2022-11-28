package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/beloslav13/servernotes/internal/models"
	"github.com/beloslav13/servernotes/internal/notes"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func NewNoteHandler(r notes.Repository, log logger.Logger) Note {
	return &note{
		repository: r,
		handler: handler{
			logger: &log,
		},
	}
}

func (n *note) Register(router *mux.Router) {
	router.HandleFunc("/notes/", n.Create).Methods("POST")
	router.HandleFunc("/notes/", n.GetAll).Methods("GET")
	router.HandleFunc("/notes/{id:[0-9]+}/", n.Get).Methods("GET")
	router.HandleFunc("/notes/{id:[0-9]+}/", n.Delete).Methods("DELETE")
	http.Handle("/", router)
}

// Get note in postgres.
func (n *note) Get(w http.ResponseWriter, r *http.Request) {
	n.logger.Infoln("Handler GetNote")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok {
		n.logger.Errorf("id does not exit in vars: %v", vars)
		return
	}

	id, _ := strconv.Atoi(sid)
	note, err := n.repository.Get(httpContext, id)
	if err != nil {
		n.logger.Errorf("handler GetNote error get: %s. id: %d", err.Error(), id)
		response(w, fmt.Sprintf("cannot get note with id: %d", id), http.StatusBadRequest, err.Error(), nil)
		return
	}

	response(w, "get ok", http.StatusOK, nil, note)
}

func (n *note) GetAll(w http.ResponseWriter, r *http.Request) {
	n.logger.Infoln("Handler GetAllNotes")
	w.Header().Set("Content-Type", "application/json")

	notes, err := n.repository.GetAll(httpContext)
	if err != nil {
		n.logger.Errorf("handler GetAllNotes error get all: %s", err.Error())
		response(w, "cannot get all notes:", http.StatusBadRequest, err.Error(), nil)
		return
	}
	response(w, "get all ok", http.StatusOK, nil, notes)
}

// Create note in postgres.
func (n *note) Create(w http.ResponseWriter, r *http.Request) {
	n.logger.Infoln("Handler CreateNote")
	w.Header().Set("Content-Type", "application/json")

	var note models.Note
	_ = json.NewDecoder(r.Body).Decode(&note)

	// Если валидация структуры прошла успешно, создаём заметку в БД.
	if n.Validate(w, note) {
		return
	}
	id, err := n.repository.Create(httpContext, &note)
	if err != nil {
		n.logger.Errorf("handler cannot save: %w\ndata: %w", err, note)
		response(w, "handler cannot save...", http.StatusBadRequest, err.Error(), nil)
		return
	}
	note.Id = uint(id)
	// Заметка создана в бд без ошибок, создаём json ответ
	response(w, "save ok", http.StatusCreated, nil, note)
}

// Delete note in postgres.
func (n *note) Delete(w http.ResponseWriter, r *http.Request) {
	n.logger.Infoln("Handler DeleteNote")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok {
		n.logger.Errorf("id does not exit in vars: %v", vars)
		return
	}

	id, _ := strconv.Atoi(sid)
	err := n.repository.Delete(httpContext, id)
	if err != nil {
		n.logger.Errorf("handler DeleteNote error delete: %s. id: %d", err.Error(), id)
		response(w, fmt.Sprintf("cannot delete note with id: %d", id), http.StatusBadRequest, err.Error(), nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
