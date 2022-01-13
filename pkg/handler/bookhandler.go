package handler

import (
	"encoding/json"
	library_app "github.com/Nkez/lib-app.git"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Create Book handler")
	var input library_app.Book
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&input); err != nil {

		http.Error(w, err.Error(), 400)
		return
	}

	bookId, err := h.service.CreateBook(input) //       .CreateUser(&input)
	if err != nil {

		http.Error(w, err.Error(), 500)
		return
	}

	header := w.Header()
	header.Add("id", strconv.Itoa(bookId))
}

func (h Handler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Find Books handler")
	var listBooks []library_app.ReturnBook

	listBooks, err := h.service.GetAllBooks()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(&listBooks)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {

		http.Error(w, err.Error(), 500)
		return
	}
}
