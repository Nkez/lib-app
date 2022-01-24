package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Nkez/library-app.git/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type BookPostgres struct {
	db *sqlx.DB
}

func NewBookPostgres(db *sqlx.DB) *BookPostgres {
	return &BookPostgres{db: db}
}

//// Book

func (r *BookPostgres) GetAllBooks(page, limit string) ([]models.ReturnBook, error) {
	logrus.Info("Find All Books DB")
	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)
	pageForSql := (p - 1) * 5
	var books []models.ReturnBook
	query := fmt.Sprintf(`SELECT id, book_title, year_of_published,inventory_count,books_in_lib
								FROM books 
								ORDER BY inventory_count LIMIT %d OFFSET %d`, l, pageForSql)
	if err := r.db.Select(&books, query); err != nil {
		return nil, err
	}
	logrus.Info("Find All Books Genres DB")
	queryGen := fmt.Sprintf(`SELECT   g.genre
								FROM books AS b
								INNER JOIN book_genre AS bg
									on b.id = bg.book_id
								inner join genres AS g
									on g.id = bg.genre_id
								where b.id  = $1
								ORDER BY inventory_count DESC`)
	for i, item := range books {
		var genres []string
		r.db.Select(&genres, queryGen, item.Id)
		books[i].Genre = genres
	}
	return books, nil
}
func (r *BookPostgres) InsertBookInfo(yearOfPublished, registrationDate time.Time, book models.Book) (int, error) {
	logrus.Info("Insert  Book  Info DB")
	var id int
	query := fmt.Sprintf(`
				INSERT INTO %s(book_title,book_title_native,book_price,
				inventory_count,books_in_lib,one_day_price,year_of_published,registration_date,
				number_of_pages,book_state,hide_book,rating
				) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id`, bookTable)
	row := r.db.QueryRow(query, book.BookTitle, book.BookTitleNative, book.BookPrice,
		book.InventoryCount, book.TotalCount, book.OneDayPrice, yearOfPublished, registrationDate,
		book.NumberOfPages, book.BookState, book.HideBook, book.Rating)
	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

//Genre

func (r *BookPostgres) InsertGenre(genres models.CreateGenre) (int, error) {
	logrus.Info("Insert Genre DB")
	var id int
	query := fmt.Sprintf(`INSERT INTO %s(genre) VALUES ($1) RETURNING id`, genreTable)
	row := r.db.QueryRow(query, genres.Genre)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error whith genre table db:%s, %s", genres.Genre, id)
	}
	return id, nil
}
func (r *BookPostgres) GetAllGenres() ([]string, error) {
	logrus.Info("Find genre Db")
	var genres []string
	query := fmt.Sprintf(`select g.genre from genres g`)
	r.db.Select(&genres, query)
	return genres, nil
}
func (r *BookPostgres) JoinBookGenres(id int, book models.Book) error {
	logrus.Info("join book genres")
	query := fmt.Sprintf(`insert into %s values($1,$2)`, bookGenreTable)
	for i := 0; i < len(book.Genres); i++ {
		r.db.QueryRow(query, id, book.Genres[i].Genre)
	}
	return nil
}
func (r *BookPostgres) GetGenres(id int) ([]string, error) {
	logrus.Info("get genres from db")
	var genres []string
	query := fmt.Sprintf(`SELECT g.genre
								FROM %s AS bg
								INNER JOIN %s  AS g
								ON g.id = bg.genre_id`, bookGenreTable, genreTable)
	if row := r.db.Select(&genres, query); row != nil {
		return nil, fmt.Errorf("error scaning genres, book id %s", id)
	}
	return genres, nil
}
func (r *BookPostgres) DeleteGenre(id int) error {
	logrus.Info("delete genre")
	var checkGenre int
	query := fmt.Sprintf(`select g.id  from  genres g 
						join book_genre bg  on bg.genre_id = g.id
						join books_users bu on bu.is_return = false
						where  g.id = $1`)
	r.db.QueryRow(query, id).Scan(&checkGenre)
	fmt.Println(checkGenre)
	if checkGenre == 0 {
		queryG := fmt.Sprintf(`delete from genres where id = $1`)
		r.db.QueryRow(queryG, id)
	}
	return nil
}

//Author

func (r *BookPostgres) GetAutInfo(id int) (models.CreateAuthor, error) {
	logrus.Info("find author")
	fmt.Println(id)
	var aut models.CreateAuthor
	query := fmt.Sprintf(`select a.id ,a.author_firstname ,a.author_lastname,a.foto_name ,a.author_photo from authors a  where id = $1 `)
	if err := r.db.QueryRow(query, id).Scan(&aut.Id, &aut.FirstName, &aut.LastName, &aut.PhotoName, &aut.AuthorPhoto); err != nil {
		return aut, err
	}
	return aut, nil

}
func (r *BookPostgres) ChangeAutPhoto(author models.CreateAuthor) {
	logrus.Info("change aut photo")
	quert := fmt.Sprintf(`UPDATE authors set author_photo = $2, foto_name = $3 where id = $1`)
	r.db.QueryRow(quert, author.Id, author.AuthorPhoto, author.PhotoName)
}
func (r *BookPostgres) InsertAuthorInfo(author models.CreateAuthor) (int, error) {
	logrus.Info("Insert Author Info DB")
	var id int
	query := fmt.Sprintf(`INSERT INTO %s(author_firstname,author_lastname,foto_name,author_photo)
								 VALUES ($1,$2,$3,$4) RETURNING id`, authorsTable)
	row := r.db.QueryRow(query, author.FirstName, author.LastName, author.PhotoName, author.AuthorPhoto)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error whith author info table db:%s, %s", author.LastName, id)
	}
	return id, nil
}
func (r *BookPostgres) GetAutPhoto(id int) (path, photo string, err error) {

	query := fmt.Sprintf(`    select a.author_photo, a.foto_name from  authors a  where a.id = $1`)
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&path, &photo); err != nil {
		return "", "", fmt.Errorf("error whith scaning path")
	}
	return path, photo, nil
}
func (r *BookPostgres) JoinBookAuthor(id int, book models.Book) error {
	logrus.Info("join book author")
	query := fmt.Sprintf(`insert into %s values($1,$2)`, bookAuthorsTable)
	for i := 0; i < len(book.Authors); i++ {
		r.db.QueryRow(query, id, book.Authors[i].Id)
	}

	return nil
}
func (r *BookPostgres) ChangeAuthorPhoto(id int, paths string) error {
	logrus.Info("change aut photo db")
	query := fmt.Sprintf(`update  authors set author_photo = $1 where id = %2`)
	if err := r.db.QueryRow(query, paths, id); err != nil {
		return fmt.Errorf("error change aut photo db")
	}
	return nil
}

//Photo

func (r *BookPostgres) BookPhoto(idPhoto int) (path, name string, err error) {

	query := fmt.Sprintf(`    select bp.book_photo, bp.foto_name from books_photo bp 
   				inner join books_books_photo bbp 
				on  bp.id = bbp.books_photo_id 
				inner join books b 
   				on b.id =bbp.books_id 
				where b.id = $1`)
	row := r.db.QueryRow(query, idPhoto)
	if err := row.Scan(&path, &name); err != nil {
		return "", "", fmt.Errorf("error whith scaning path")
	}
	return path, name, nil
}
func (r *BookPostgres) JoinBookPhoto(idBook int, p, n []string) error {
	logrus.Info("join book photo")
	var ids []int
	var id int
	query := fmt.Sprintf(`insert into books_photo(foto_name,book_photo) values($1,$2) RETURNING id`)
	for i := 0; i < len(p); i++ {
		err := r.db.QueryRow(query, n[i], p[i]).Scan(&id)
		if err != nil {
			return fmt.Errorf("insert in book_photo")
		}
		ids = append(ids, id)
	}

	fmt.Println(ids)

	queryId := fmt.Sprintf(`insert into books_books_photo values($1,$2)`)
	for _, v := range ids {
		r.db.QueryRow(queryId, idBook, v)
	}
	return nil
}
func (r *BookPostgres) ChangeBookPhoto(newPath, newName, paths string) {
	logrus.Info("change book foto")
	query := fmt.Sprintf(`update books_photo set foto_name = $2, book_photo =$3 where book_photo = $1`)
	r.db.QueryRow(query, paths, newName, newPath)
}

//Defect

func (r *BookPostgres) JoinDefetBookPhoto(idBook int, defect string, paths []string) error {
	logrus.Info("join book defect copies")
	var id int
	var idCopies []int
	que := fmt.Sprintf(`insert into book_copies(defect,defect_photo) values ($1,$2) RETURNING id`)
	for i := 0; i < len(paths); i++ {
		row := r.db.QueryRow(que, defect, paths[i])
		if err := row.Scan(&id); err != nil {
			return fmt.Errorf("error in join defect photo")
		}
		idCopies = append(idCopies, id)
	}
	fmt.Println(idCopies)
	query := fmt.Sprintf(`insert into books_book_copies values($1,$2)`)
	for i := 0; i < len(idCopies); i++ {
		r.db.QueryRow(query, idBook, idCopies[i])
		fmt.Println(idCopies[i])
	}
	return nil
}
func (r *BookPostgres) GetDefectPhoto(idPhoto int) (string, error) {
	var path string
	query := fmt.Sprintf(`    select  bc.defect_photo from  book_copies bc
	inner join books_book_copies bbc 
	on bc.id = bbc.book_copies_id 
	inner join books b 
	on b.id  = bbc.books_id 
	where b.id = $1`)
	row := r.db.QueryRow(query, idPhoto)
	if err := row.Scan(&path); err != nil {
		return "", fmt.Errorf("error whith scaning path")
	}
	return path, nil
}

//Find smt..

func (r *BookPostgres) GetByTitle(title string) ([]models.ReturnBook, error) {
	logrus.Info("Find All Books DB")

	var books []models.ReturnBook
	query := fmt.Sprintf(`SELECT id, book_title, year_of_published,books_in_lib,inventory_count
								FROM books 
								where book_title = $1`)
	r.db.Select(&books, query, title)
	logrus.Info("Find All Books Genres DB")
	queryGen := fmt.Sprintf(`SELECT   g.genre
								FROM books AS b
								INNER JOIN book_genre AS bg
									on b.id = bg.book_id
								inner join genres AS g
									on g.id = bg.genre_id
								where b.book_title  = $1
								`)
	for i, item := range books {
		var genres []string
		r.db.Select(&genres, queryGen, item.Id)
		books[i].Genre = genres
	}
	return books, nil

}
func (r *BookPostgres) GetByWord(word string) ([]string, error) {
	logrus.Info("Find genre Db")

	var s []string
	query := fmt.Sprintf(`SELECT book_title FROM books WHERE book_title LIKE $1 `)
	r.db.Select(&s, query, word+"%")
	return s, nil
}
func (r *BookPostgres) ChekRegisterUser(email string) (models.User, error) {
	user := models.User{}
	logrus.Info("Find Register User DB")
	query := fmt.Sprintf(`SELECT id,email_address,last_name,first_name FROM %s WHERE  email_address = $1`, userTable)
	if err := r.db.QueryRow(query,
		email,
	).Scan(
		&user.Id,
		&user.EmailAddress,
		&user.LastName,
		&user.FirstName,
	); err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("Create new user")
		}
		return user, err
	}
	return user, nil
}
func (r *BookPostgres) GetByTopRating() []models.TopRating {
	var top []models.TopRating
	query := fmt.Sprintf(`select b.id , b.book_title, b.rating from books b order by rating desc limit  3`)
	r.db.Select(&top, query)
	return top
}
