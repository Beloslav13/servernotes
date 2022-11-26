package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/beloslav13/servernotes/internal/database"
	"github.com/beloslav13/servernotes/internal/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// GetNote get note in database.
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

func (h *handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("Handler GetAllNotes")
	w.Header().Set("Content-Type", "application/json")

	notes, err := database.GetAllNotes(httpContext)
	if err != nil {
		h.logger.Errorf("handler GetAllNotes error get all: %s", err.Error())
		response(w, "cannot get all notes:", http.StatusBadRequest, err.Error(), nil)
		return
	}
	response(w, "get all ok", http.StatusOK, nil, notes)
}

// CreateNote create note in database.
func (h *handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("Handler CreateNote")
	w.Header().Set("Content-Type", "application/json")

	var note models.Note
	_ = json.NewDecoder(r.Body).Decode(&note)

	// Если валидация структуры прошла успешно, создаём заметку в БД.
	if validateCreateNote(w, note, h) {
		return
	}
	id, err := database.CreateNote(httpContext, &note)
	if err != nil {
		h.logger.Errorf("handler cannot save: %w\ndata: %w", err, note)
		response(w, "handler cannot save...", http.StatusBadRequest, err.Error(), nil)
		return
	}
	note.Id = uint(id)
	// Заметка создана в бд без ошибок, создаём json ответ
	response(w, "save ok", http.StatusCreated, nil, note)
}

// DeleteNote delete note in database.
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
