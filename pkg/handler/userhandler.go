package handler

import (
	"encoding/json"
	library_app "github.com/Nkez/lib-app.git"
	"strconv"

	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Create User handler")

	var input library_app.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&input); err != nil {

		http.Error(w, err.Error(), 400)
		return
	}

	userId, err := h.service.CreateUser(input) //       .CreateUser(&input)
	if err != nil {

		http.Error(w, err.Error(), 500)
		return
	}

	header := w.Header()
	header.Add("id", strconv.Itoa(userId))
	//var input library_app.User
	//decoder := json.NewDecoder(r.Body)
	//
	//if err := decoder.Decode(&input); err != nil {
	//
	//	http.Error(w, err.Error(), 400)
	//	return
	//}
	//
	//userId, err := h.service.User.CreateUser(input)
	//if err != nil{
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//
	//header := w.Header()
	//header.Add("id", strconv.Itoa(userId))
}

func (h Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Find Users handler")
	var listUsers []library_app.User

	listUsers, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(&listUsers)
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
