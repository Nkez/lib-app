package handler

import (
	"encoding/json"
	library_app "github.com/Nkez/lib-app.git"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var cart library_app.OrderCart
	v := r.URL.Query()
	email := v.Get("email")
	book1 := v.Get("book1")
	book2 := v.Get("book2")
	book3 := v.Get("book3")
	book4 := v.Get("book4")
	book5 := v.Get("book5")
	books := []string{book1, book2, book3, book4, book5}
	cart, err := h.service.CreateCart(email, books)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&cart)
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

func (h Handler) FindRegisterUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Find Register User handler")

	v := r.URL.Query()
	email := v.Get("email")
	var regUser library_app.User
	regUser, err := h.service.FindRegisterUser(email)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(&regUser)
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
