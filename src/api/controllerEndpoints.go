package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"sample-golang-project/kafka"
	"sample-golang-project/model"
	"sample-golang-project/types"
)

type ControllerEndpoint struct {
	router     *mux.Router
	repository types.ControllerRepository
}

func NewControllerEndpoints(router *mux.Router, repository types.ControllerRepository) *ControllerEndpoint {

	c := &ControllerEndpoint{
		router:     router,
		repository: repository,
	}
	c.router.HandleFunc("/api/send", c.sendMessage).Methods("POST")
	return c
}

func (c *ControllerEndpoint) sendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var message *model.Controller
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	kafka.PublishMessage(messageBytes)
	save, err := c.repository.Save(ctx, message)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(save); err != nil {
		return
	}
}
