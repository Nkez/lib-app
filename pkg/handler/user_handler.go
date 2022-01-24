package handler

import (
	"encoding/json"
	"github.com/Nkez/library-app.git/models"
	"strconv"

	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Create User handler")

	var input models.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&input); err != nil {

		logrus.WithError(err).Error("error decode user struct")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId, err := h.service.CreateUser(input)
	if err != nil {
		logrus.WithError(err).Error("error from creating user service")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	header := w.Header()
	header.Add("id", strconv.Itoa(userId))
}

func (h Handler) FindByName(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Find user by name")

	v := r.URL.Query()
	name := v.Get("name")
	var user []models.UserName
	user = h.service.GetByName(name)

	output, err := json.Marshal(&user)
	if err != nil {
		logrus.WithError(err).Error("error  marshaling register user")
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		logrus.WithError(err).Error("error  writing output register user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h Handler) FindByWord(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Find book by title")

	v := r.URL.Query()
	name := v.Get("name")
	var user []models.UserName
	user = h.service.FindByWord(name)

	output, err := json.Marshal(&user)
	if err != nil {
		logrus.WithError(err).Error("error  marshaling register user")
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		logrus.WithError(err).Error("error  writing output register user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Find Users handler")
	var listUsers []models.User
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	listUsers, err := h.service.GetAllUsers(page, limit)
	if err != nil {
		logrus.WithError(err).Error("error get all user service")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(&listUsers)
	if err != nil {
		logrus.WithError(err).Error("error marshaling all users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {

		logrus.WithError(err).Error("error writing output all users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
