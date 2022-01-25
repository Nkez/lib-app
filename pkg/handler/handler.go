package handler

import (
	"github.com/Nkez/library-app.git/pkg/services"
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
	////Book
	r.HandleFunc("/book/all", h.GetAllBooks).Methods("GET")
	r.HandleFunc("/book", h.CreateBook).Methods("POST")
	r.HandleFunc("/book", h.GetBookById).Methods("GET")
	r.HandleFunc("/book/download", h.DownloadBookPhoto).Methods("GET")
	r.HandleFunc("/book/change", h.ChangeBookPhoto).Methods("POST")

	r.HandleFunc("/book/", h.GetBookByTitle).Methods("GET")
	r.HandleFunc("/book/word", h.GetByWord).Methods("GET")
	r.HandleFunc("/book/image", h.BookPhoto).Methods("GET")
	r.HandleFunc("/book/rating", h.GetByTopRating).Methods("GET")

	r.HandleFunc("/create/genre", h.CreateGenre).Methods("POST")
	r.HandleFunc("/book/add/photos", h.AddBookPhotos).Methods("POST")
	r.HandleFunc("/book/defect", h.DefectPhoto).Methods("POST")
	r.HandleFunc("/book/defect/get", h.GetDefectPhoto).Methods("Get")
	r.HandleFunc("/genres", h.GetALlGenres).Methods("GET")
	r.HandleFunc("/genres/delete", h.DeleteGenre).Methods("GET")
	////Author

	r.HandleFunc("/create/author", h.CreateAuthor).Methods("POST")
	r.HandleFunc("/authors/image", h.AuthorPhoto).Methods("GET")
	r.HandleFunc("/authors/download", h.DowlondAutPhoto).Methods("GET")
	r.HandleFunc("/authors/change", h.ChangeAuthorPhoto).Methods("POST")
	//User
	r.HandleFunc("/user", h.GetAllUsers).Methods("GET")
	r.HandleFunc("/user/create", h.CreateUser).Methods("POST")
	r.HandleFunc("/user/name", h.FindByName).Methods("GET")
	r.HandleFunc("/find/word", h.FindByWord).Methods("GET")

	//Order
	//r.HandleFunc("/find", h.FindUser).Methods("GET")
	r.HandleFunc("/order", h.OrderBook).Methods("POST")
	r.HandleFunc("/order/all", h.GetAllOrder).Methods("GET")
	//Return
	r.HandleFunc("/return", h.ReturnBook).Methods("POST")
	r.HandleFunc("/debors", h.GetAllDebors).Methods("GET")

	// Copies
	r.HandleFunc("/copies", h.GetAllCopies).Methods("GET")
	r.HandleFunc("/copies/id", h.GetCopiesByIdAndBookID).Methods("GET")

	return r
}
