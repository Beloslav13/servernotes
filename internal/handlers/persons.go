package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/beloslav13/servernotes/internal/models"
	pr "github.com/beloslav13/servernotes/internal/person"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func NewPersonHandler(r pr.Repository, log logger.Logger) Person {
	return &person{
		repository: r,
		handler: handler{
			logger: &log,
		},
	}
}

func (p *person) Register(router *mux.Router) {
	router.HandleFunc("/persons/", p.Create).Methods("POST")
	router.HandleFunc("/persons/", p.GetAll).Methods("GET")
	router.HandleFunc("/persons/{id:[0-9]+}/", p.Get).Methods("GET")
	router.HandleFunc("/persons/{id:[0-9]+}/", p.Delete).Methods("DELETE")
}

// Create person in postgres.
func (p *person) Create(w http.ResponseWriter, r *http.Request) {
	p.logger.Infoln("Handler CreatePerson")
	w.Header().Set("Content-Type", "application/json")

	var person models.Person
	_ = json.NewDecoder(r.Body).Decode(&person)

	// Если валидация структуры прошла успешно, создаём заметку в БД.
	if p.Validate(w, person) {
		return
	}
	id, err := p.repository.Create(httpContext, &person)
	if err != nil {
		p.logger.Errorf("handler cannot save: %w\ndata: %w", err, person)
		response(w, "handler cannot save...", http.StatusBadRequest, err.Error(), nil)
		return
	}
	person.Id = uint(id)
	// Заметка создана в бд без ошибок, создаём json ответ
	response(w, "save ok", http.StatusCreated, nil, person)
}

func (p *person) Get(w http.ResponseWriter, r *http.Request) {
	p.logger.Infoln("Handler GetPerson")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok {
		p.logger.Errorf("id does not exit in vars: %v", vars)
		return
	}

	id, _ := strconv.Atoi(sid)
	person, err := p.repository.Get(httpContext, id)
	if err != nil {
		p.logger.Errorf("handler GetPerson error get: %s. id: %d", err.Error(), id)
		response(w, fmt.Sprintf("cannot get person with id: %d", id), http.StatusBadRequest, err.Error(), nil)
		return
	}

	response(w, "get ok", http.StatusOK, nil, person)
}

func (p *person) GetAll(w http.ResponseWriter, r *http.Request) {
	p.logger.Infoln("Handler GetAllPersons")
	w.Header().Set("Content-Type", "application/json")

	persons, err := p.repository.GetAll(httpContext)
	if err != nil {
		p.logger.Errorf("handler GetAllPersons error get all: %s", err.Error())
		response(w, "cannot get all persons:", http.StatusBadRequest, err.Error(), nil)
		return
	}
	response(w, "get all ok", http.StatusOK, nil, persons)
}

func (p *person) Delete(w http.ResponseWriter, r *http.Request) {
	p.logger.Infoln("Handler DeletePerson")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok {
		p.logger.Errorf("id does not exit in vars: %v", vars)
		return
	}

	id, _ := strconv.Atoi(sid)
	err := p.repository.Delete(httpContext, id)
	if err != nil {
		p.logger.Errorf("handler DeletePerson error delete: %s. id: %d", err.Error(), id)
		response(w, fmt.Sprintf("cannot delete person with id: %d", id), http.StatusBadRequest, err.Error(), nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
