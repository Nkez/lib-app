package book_handlers

import (
	"encoding/json"
	"github.com/Nkez/lib-app.git/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {

	logrus.Info("Create Book handler")
	var input models.Book
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&input); err != nil {

		logrus.WithError(err).Error("error with decode")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bookId, err := h.service.CreateBook(input)
	if err != nil {

		logrus.WithError(err).Error("error with creating book")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	header := w.Header()
	header.Add("id", strconv.Itoa(bookId))
}

//func (h Handler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
//	logrus.Info("Find Books handler")
//	var listBooks []models.ReturnBook
//
//	listBooks, err := h.service.GetAllBooks()
//	if err != nil {
//		logrus.WithError(err).Error("error with getting books")
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	output, err := json.Marshal(&listBooks)
//	if err != nil {
//		logrus.WithError(err).Error("error marshaling books")
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	_, err = w.Write(output)
//	if err != nil {
//
//		logrus.WithError(err).Error("error writing output")
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//}
//
//func (h Handler) ChekRegisterUser(w http.ResponseWriter, r *http.Request) {
//	logrus.Info("Find Register User handler")
//
//	v := r.URL.Query()
//	email := v.Get("email")
//	var regUser models.User
//	regUser, err := h.service.ChekRegisterUser(email)
//	if err != nil {
//		logrus.WithError(err).Error("error with checking user")
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	output, err := json.Marshal(&regUser)
//	if err != nil {
//		logrus.WithError(err).Error("error  marshaling register user")
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//	w.Header().Set("Content-Type", "application/json")
//	_, err = w.Write(output)
//	if err != nil {
//		logrus.WithError(err).Error("error  writing output register user")
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//}
