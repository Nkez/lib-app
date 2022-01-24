package handler

import (
	"encoding/json"
	"github.com/Nkez/library-app.git/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) OrderBook(w http.ResponseWriter, r *http.Request) {

	logrus.Info("Create order handler")
	var input models.OrderInput
	var result models.Order
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := h.service.CreateOrder(input)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	output, err := json.Marshal(&result)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (h Handler) GetAllOrder(w http.ResponseWriter, r *http.Request) {
	logrus.Info("get all order")
	var ord []models.InfoOrdDept
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	ord, err := h.service.GetAllOrder(page, limit)
	if err != nil {
		logrus.WithError(err).Error("error with getting books")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(&ord)
	if err != nil {
		logrus.WithError(err).Error("error marshaling books")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {

		logrus.WithError(err).Error("error writing output")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
