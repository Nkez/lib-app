package handler

import (
	"github.com/Nkez/lib-app.git/pkg/services"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	http.Handle("/", r)
	r.HandleFunc("/book", h.GetAllBooks).Methods("GET")
	r.HandleFunc("/book/create", h.CreateBook).Methods("POST")

	r.HandleFunc("/user", h.GetAllUsers).Methods("GET")
	r.HandleFunc("/find", h.FindRegisterUser).Methods("GET")
	r.HandleFunc("/user/create", h.CreateUser).Methods("POST")

	r.HandleFunc("/cart/create", h.CreateCard).Methods("GET")

	//r.HandleFunc("/return", h.ReturnBook).Methods("DELETE")
	//r.HandleFunc("/writeoff", h.WriteOffBook).Methods("DELETE")
	//r.HandleFunc("/profitability", h.LibraryProfitabilaty).Methods("GET")

	return r
}
