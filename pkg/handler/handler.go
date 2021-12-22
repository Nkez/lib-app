package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"html"
	"net/http"
)

type Handler struct {
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/books", GetAllBooks).Methods("GET")
	r.HandleFunc("/users", GetAllUsers).Methods("GET")
	r.HandleFunc("/create/book", CreateBook).Methods("POST")
	r.HandleFunc("/create/user", CreateUser).Methods("POST")
	return r
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Info("create user")
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	log.Info("create book")
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GetAllBooks, %q", html.EscapeString(r.URL.Path))
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GetAllBooks, %q", html.EscapeString(r.URL.Path))
}
