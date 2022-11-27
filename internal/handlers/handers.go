package handlers

import (
	"context"
	"encoding/json"
	"github.com/beloslav13/servernotes/internal/models"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var httpContext = context.Background()

type handler struct {
	logger logger.Logger
}

type note struct {
	handler
}

type person struct {
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
