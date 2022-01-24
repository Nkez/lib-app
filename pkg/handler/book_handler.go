package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Nkez/library-app.git/models"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

///Book

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Create Book handler")
	var book models.Book

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&book); err != nil {

		logrus.WithError(err).Error("error decode user struct")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bookId, err := h.service.CreateBook(book)
	if err != nil {
		logrus.WithError(err).Error("error from creating user service")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	header := w.Header()
	header.Add("id", strconv.Itoa(bookId))
}
func (h Handler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Find Books handler")
	var listBooks []models.ReturnBook
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	listBooks, err := h.service.GetAllBooks(page, limit)
	if err != nil {
		logrus.WithError(err).Error("error with getting books")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(&listBooks)
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

///Author

func (h *Handler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Working createAuthor")
	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		logrus.Errorf("parsing error :%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var input models.CreateAuthor
	input.FirstName = r.PostFormValue("name")
	input.LastName = r.PostFormValue("lastname")
	input.PhotoName = r.PostFormValue("photoname")

	photo, err := GetAuthorPhoto(r, input)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	input.AuthorPhoto = photo
	authorId, err := h.service.CreateAuthor(input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	header := w.Header()
	header.Add("id", strconv.Itoa(authorId))
}

///Genre

func (h *Handler) CreateGenre(w http.ResponseWriter, r *http.Request) {

	logrus.Info("Create Genre handler")
	var input models.CreateGenre
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&input); err != nil {

		fmt.Errorf("error whit decode %s", input)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bookId, err := h.service.CreateGenre(input)
	if err != nil {

		logrus.WithError(err).Error("error with creating book")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	header := w.Header()
	header.Add("id", strconv.Itoa(bookId))
}
func (h *Handler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	logrus.Info("delete genre handler")
	id := r.URL.Query().Get("id")
	idG, _ := strconv.Atoi(id)
	fmt.Println(idG)
	if err := h.service.DeleteGenre(idG); err != nil {
		logrus.Errorf("error delete genre :%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

//Change photo

func (h *Handler) ChangeAuthorPhoto(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Working change aut photo")
	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		logrus.Errorf("parsing error :%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var aut models.CreateAuthor
	id := r.PostFormValue("id")
	name := r.PostFormValue("photoname")
	idAut, _ := strconv.Atoi(id)
	aut, err := h.service.GetAutInfo(idAut)
	if err != nil {
		logrus.Errorf("error find aut info:%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	os.Remove(aut.AuthorPhoto)
	aut.PhotoName = name
	newPath, err := GetAuthorPhoto(r, aut)
	if err != nil {
		logrus.Errorf("error find new path:%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	aut.AuthorPhoto = newPath
	h.service.ChangeAutPhoto(aut)
}
func (h *Handler) ChangeBookPhoto(w http.ResponseWriter, r *http.Request) {
	logrus.Info(" change book photo")
	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		logrus.Errorf("parsing error :%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	idBook := r.PostFormValue("idBook")
	id, err := strconv.Atoi(idBook)

	paths, _, err := h.service.BookPhoto(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := os.Remove(paths); err != nil {
		logrus.Errorf("errro delete book photo :%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(paths)
	var newPath string
	var newName string
	p, n, err := GetBookPhoto(r, idBook)
	for _, s := range p {
		newPath = s
	}
	for _, v := range n {
		newName = v
	}

	h.service.ChangeBookPhoto(newPath, newName, paths)

}

///Photo

func (h *Handler) DefectPhoto(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Add defect phot")
	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		logrus.Errorf("parsing error :%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	idBook := r.PostFormValue("idBook")
	defect := r.PostFormValue("defect")

	paths, err := GetDefectBookPhoto(r, idBook)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	id, err := strconv.Atoi(idBook)
	if err != nil {
		return
	}
	if err = h.service.JoinDefetBookPhoto(id, defect, paths); err != nil {
		return
	}
}
func (h *Handler) AddBookPhotos(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Add Book photo")
	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		logrus.Errorf("parsing error :%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	idBook := r.PostFormValue("idBook")

	paths, name, err := GetBookPhoto(r, idBook)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	id, err := strconv.Atoi(idBook)
	if err != nil {
		return
	}
	if err = h.service.JoinBookPhoto(id, paths, name); err != nil {
		return
	}
}
func (h *Handler) AuthorPhoto(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idAut, _ := strconv.Atoi(id)
	var path string
	path, _, _ = h.service.GetAutPhoto(idAut)
	filename := path
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Can't open file: " + filename)
	} else {
		w.Write(file)
	}
}
func GetAuthorPhoto(r *http.Request, author models.CreateAuthor) (string, error) {
	reqFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("error form:%w", err)
	}
	defer reqFile.Close()
	filePath := fmt.Sprintf("images/authors/%s.%s", author.PhotoName, (strings.Split(fileHeader.Filename, "."))[1])
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {

		return "", fmt.Errorf("error opening file %s:%w", filePath, err)
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(reqFile)
	if err != nil {

		return "", fmt.Errorf("error from request:%w", err)
	}
	_, err = file.Write(fileBytes)
	if err != nil {

		return "", fmt.Errorf("error  writing file:%w", err)
	}

	return filePath, nil
}
func (h *Handler) BookPhoto(w http.ResponseWriter, r *http.Request) {

	bookId := r.URL.Query().Get("id")
	idAut, _ := strconv.Atoi(bookId)
	path, _, _ := h.service.BookPhoto(idAut)
	filename := path
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Can't open file: " + filename)
	} else {
		w.Write(file)
	}
}
func GetBookPhoto(r *http.Request, idBook string) (p, n []string, err error) {

	m := r.MultipartForm
	files := m.File["file"]

	fmt.Println(files)
	for i, headers := range files {
		requestfile, err := files[i].Open()
		if err != nil {
			return nil, nil, fmt.Errorf("error form:%w", err)
		}
		defer requestfile.Close()
		fileBytes, err := ioutil.ReadAll(requestfile)
		if err != nil {

			return nil, nil, fmt.Errorf("error request:%w", err)
		}
		dirPath := fmt.Sprintf("images/books/bookid%v", idBook)
		filePath := fmt.Sprintf("%s/%s", dirPath, headers.Filename)
		err = os.MkdirAll(dirPath, 0777)
		if err != nil {
			return nil, nil, fmt.Errorf("error dir (%s):%w", dirPath, err)
		}
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return nil, nil, fmt.Errorf("error opening file %s:%w", filePath, err)
		}
		defer file.Close()
		_, err = file.Write(fileBytes)
		if err != nil {
			return nil, nil, fmt.Errorf("error writing file:%w", err)
		}
		p = append(p, filePath)
		n = append(n, headers.Filename)

	}

	return p, n, nil
}
func (h *Handler) GetDefectPhoto(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idAut, _ := strconv.Atoi(id)
	var path string
	path, _ = h.service.GetDefectPhoto(idAut)
	filename := path
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Can't open file: " + filename)
	} else {
		w.Write(file)
	}
}
func GetDefectBookPhoto(r *http.Request, idBook string) ([]string, error) {
	var paths []string
	m := r.MultipartForm
	files := m.File["file"]
	fmt.Println(files)
	for i, headers := range files {
		requestfile, err := files[i].Open()
		if err != nil {
			return nil, fmt.Errorf("error form:%w", err)
		}
		defer requestfile.Close()
		fileBytes, err := ioutil.ReadAll(requestfile)
		if err != nil {

			return nil, fmt.Errorf("error request:%w", err)
		}
		dirPath := fmt.Sprintf("images/books/defect/bookid%v", idBook)
		filePath := fmt.Sprintf("%s/%s", dirPath, headers.Filename)
		err = os.MkdirAll(dirPath, 0777)
		if err != nil {
			return nil, fmt.Errorf("error dir (%s):%w", dirPath, err)
		}
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return nil, fmt.Errorf("error opening file %s:%w", filePath, err)
		}
		defer file.Close()
		_, err = file.Write(fileBytes)
		if err != nil {
			return nil, fmt.Errorf("error writing file:%w", err)
		}
		paths = append(paths, filePath)
	}
	return paths, nil
}

//// DownloadFiles

func (h *Handler) DowlondAutPhoto(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idAut, _ := strconv.Atoi(id)
	var photo string
	var path string
	fmt.Println(path)
	path, photo, _ = h.service.GetAutPhoto(idAut)
	openfile, err := os.Open(path)
	defer openfile.Close()

	if err != nil {
		http.Error(w, "File not found.", 404)
		return
	}

	downloadBytes, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
	}

	mime := http.DetectContentType(downloadBytes)

	fileSize := len(string(downloadBytes))

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", "attachment; filename="+photo+".jpg"+"")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(fileSize))
	w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

	http.ServeContent(w, r, photo+".jpg", time.Now(), bytes.NewReader(downloadBytes))

}
func (h *Handler) DownloadBookPhoto(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idAut, _ := strconv.Atoi(id)
	var photo string
	var path string
	fmt.Println(path)
	path, photo, _ = h.service.BookPhoto(idAut)
	fmt.Println(path)
	fmt.Println(photo)
	openfile, err := os.Open(path)
	defer openfile.Close()

	if err != nil {
		http.Error(w, "File not found.", 404)
		return
	}

	downloadBytes, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
	}

	mime := http.DetectContentType(downloadBytes)

	fileSize := len(string(downloadBytes))

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", "attachment; filename="+photo+".jpg"+"")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(fileSize))
	w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

	http.ServeContent(w, r, photo+".jpg", time.Now(), bytes.NewReader(downloadBytes))

}

///Find

func (h Handler) GetByTopRating(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Find book by title")

	var top []models.TopRating
	top = h.service.GetByTopRating()

	output, err := json.Marshal(&top)
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
func (h Handler) GetByWord(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Find book by title")

	v := r.URL.Query()
	word := v.Get("word")
	var book []string
	book, err := h.service.GetByWord(word)
	if err != nil {
		logrus.WithError(err).Error("error with checking user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	output, err := json.Marshal(&book)
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
func (h Handler) GetBookByTitle(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Find book by title")

	v := r.URL.Query()
	title := v.Get("title")
	var book []models.ReturnBook
	book, err := h.service.GetByTitle(title)
	if err != nil {
		logrus.WithError(err).Error("error with checking user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	output, err := json.Marshal(&book)
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
func (h Handler) GetALlGenres(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Find genre handler")
	var genres []string

	genres, err := h.service.GetAllGenres()
	if err != nil {
		logrus.WithError(err).Error("error with getting books")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(&genres)
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

// Delete Photo
