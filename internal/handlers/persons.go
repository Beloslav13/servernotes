package handlers

import (
	"encoding/json"
	"github.com/beloslav13/servernotes/internal/models"
	pr "github.com/beloslav13/servernotes/internal/person"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
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

func (p *person) GetAll(w http.ResponseWriter, r *http.Request) {

}

func (p *person) Get(w http.ResponseWriter, r *http.Request) {

}

func (p *person) Delete(w http.ResponseWriter, r *http.Request) {

}
