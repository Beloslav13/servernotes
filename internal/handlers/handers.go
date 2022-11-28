package handlers

import (
	"context"
	"encoding/json"
	"github.com/beloslav13/servernotes/internal/models"
	"github.com/beloslav13/servernotes/internal/notes"
	p "github.com/beloslav13/servernotes/internal/person"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

var httpContext = context.Background()

type Handler interface {
	Register(router *mux.Router)
	Get(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Validate(w http.ResponseWriter, m interface{}) bool
	Logger() *logger.Logger
}

type Note interface {
	Handler
}

type Person interface {
	Handler
}

type handler struct {
	logger *logger.Logger
}

func (h handler) Logger() *logger.Logger {
	return h.logger
}

type note struct {
	repository notes.Repository
	handler
}

type person struct {
	repository p.Repository
	handler
}

// Validate structure completion.
func (h *handler) Validate(w http.ResponseWriter, m interface{}) bool {
	validate := validator.New()
	err := validate.Struct(m)
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
