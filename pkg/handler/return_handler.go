package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Nkez/library-app.git/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h Handler) GetAllDebors(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Find Books handler")
	var debors []models.InfoOrdDept

	debors, err := h.service.GetAllDebors()
	for _, debor := range debors {
		fmt.Println(debor.LastName)
	}
	if err != nil {
		logrus.WithError(err).Error("error with getting books")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(&debors)
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
func (h *Handler) ReturnBook(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Create return handler")
	var input models.ReturnInput
	//var result models.Return
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&input); err != nil {

		logrus.WithError(err).Error("error decoding in return books")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := h.service.ReturnCart(input)
	if err != nil {

		logrus.WithError(err).Error("error from return books service")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(&result)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		logrus.WithError(err).Error("error writing output return books")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
