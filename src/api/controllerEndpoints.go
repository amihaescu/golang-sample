package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"sample-golang-project/model"
	"sample-golang-project/types"
)

type DeviceController struct {
	router     *mux.Router
	repository types.DeviceRepository
	output     chan *model.Device
}

func NewDeviceController(router *mux.Router, repository types.DeviceRepository, output chan *model.Device) *DeviceController {
	c := &DeviceController{
		router:     router,
		repository: repository,
		output:     output,
	}
	c.router.HandleFunc("/api/send", c.sendMessage).Methods("POST")
	return c
}

func (c *DeviceController) sendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var message *model.Device
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	save, err := c.repository.Save(ctx, message)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}
	c.output <- message

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(save); err != nil {
		return
	}
}
