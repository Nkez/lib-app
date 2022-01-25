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

type Genres struct {
	Genre string
}

type BookPostgres struct {
	db *sqlx.DB
}

func NewBookPostgres(db *sqlx.DB) *BookPostgres {
	return &BookPostgres{db: db}
}

//// Book

func (r *BookPostgres) GetAllBooks(page, limit string) ([]models.ReturnBook, error) {
	logrus.Info("find all books")
	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)
	pageForSql := (p - 1) * 5

	query := fmt.Sprintf(`SELECT b.id, b.book_title, b.year_of_published,b.inventory_count
								FROM books b
								ORDER BY inventory_count LIMIT %d OFFSET %d`, l, pageForSql)

	var books []models.ReturnBook
	row, err := r.db.Query(query)
	if err != nil {
		return books, err
	}
	defer row.Close()
	for row.Next() {
		b := models.ReturnBook{}
		if err := row.Scan(&b.Id, &b.BookTitle, &b.YearOfPublished, &b.InventoryCount); err != nil {
			return books, err
		}
		b.Genre, _ = r.GetGenres(b.Id)
		books = append(books, b)
	}

	return books, nil
}
func (r *BookPostgres) InsertBookInfo(yearOfPublished, registrationDate time.Time, book models.Book) (int, error) {
	logrus.Info(fmt.Sprintf("insert book info, id %v", book.Id))
	var id int
	query := fmt.Sprintf(`
				INSERT INTO %s(book_title,book_title_native,book_price,
				inventory_count,books_in_lib,one_day_price,year_of_published,registration_date,
				number_of_pages,rating
				) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`, bookTable)
	row := r.db.QueryRow(query, book.BookTitle, book.BookTitleNative, book.BookPrice,
		book.InventoryCount, book.TotalCount, book.OneDayPrice, yearOfPublished, registrationDate,
		book.NumberOfPages, book.Rating)
	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

//Genre

func (r *BookPostgres) InsertGenre(genres models.CreateGenre) (int, error) {
	logrus.Info(fmt.Sprintf("insert genre, id %s", genres.Genre))
	var id int
	query := fmt.Sprintf(`INSERT INTO %s(genre) VALUES ($1) RETURNING id`, genreTable)
	row := r.db.QueryRow(query, genres.Genre)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error whith genre table db:%s, %s", genres.Genre, id)
	}
	return id, nil
}
func (r *BookPostgres) GetAllGenres() ([]string, error) {
	logrus.Info("find genres ")
	var genres []string
	query := fmt.Sprintf(`select g.genre from genres g`)
	r.db.Select(&genres, query)
	return genres, nil
}

func (r *BookPostgres) JoinBookGenres(id int, book models.Book) error {
	logrus.Info(fmt.Sprintf("join book : %v, genres : %v", book.Id, id))
	query := fmt.Sprintf(`insert into %s values($1,$2)`, bookGenreTable)
	for i := 0; i < len(book.Genres); i++ {
		r.db.QueryRow(query, id, book.Genres[i].Genre)
	}
	return nil
}
func (r *BookPostgres) GetGenres(id int) ([]string, error) {
	logrus.Info(fmt.Sprintf("get genres %v", id))
	var genres []string
	query := fmt.Sprintf(`SELECT g.genre
								FROM books b
								left join book_genre bg on bg.book_id  = b.id 
								join genres g on bg.genre_id  = g.id 
								where b.id = $1`)
	if row := r.db.Select(&genres, query, id); row != nil {
		return nil, fmt.Errorf("error scaning genres, book id %s", id)
	}
	return genres, nil
}
func (r *BookPostgres) DeleteGenre(id int) error {
	logrus.Info(fmt.Sprintf("delete genre, id %v", id))
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
	logrus.Info(fmt.Sprintf("find author, id %v", id))
	fmt.Println(id)
	var aut models.CreateAuthor
	query := fmt.Sprintf(`select a.id ,a.author_firstname ,a.author_lastname,a.foto_name ,a.author_photo from authors a  where id = $1 `)
	if err := r.db.QueryRow(query, id).Scan(&aut.Id, &aut.FirstName, &aut.LastName, &aut.PhotoName, &aut.AuthorPhoto); err != nil {
		return aut, err
	}
	return aut, nil

}
func (r *BookPostgres) ChangeAutPhoto(author models.CreateAuthor) {
	logrus.Info(fmt.Sprintf("change aut photo, id %v", author.Id))
	query := fmt.Sprintf(`UPDATE authors set author_photo = $2, foto_name = $3 where id = $1`)
	r.db.QueryRow(query, author.Id, author.AuthorPhoto, author.PhotoName)
}
func (r *BookPostgres) InsertAuthorInfo(author models.CreateAuthor) (int, error) {
	logrus.Info(fmt.Sprintf("insert aut info :%v", author.Id))
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
	logrus.Info(fmt.Sprintf("get aut photo :%v", id))
	query := fmt.Sprintf(`    select a.author_photo, a.foto_name from  authors a  where a.id = $1`)
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&path, &photo); err != nil {
		return "", "", fmt.Errorf("error whith scaning path")
	}
	return path, photo, nil
}
func (r *BookPostgres) JoinBookAuthor(id int, book models.Book) error {
	logrus.Info(fmt.Sprintf("join book :%v, author :%v", book.Id, id))
	query := fmt.Sprintf(`insert into %s values($1,$2)`, bookAuthorsTable)
	for i := 0; i < len(book.Authors); i++ {
		r.db.QueryRow(query, id, book.Authors[i].Id)
	}
	return nil
}
func (r *BookPostgres) ChangeAuthorPhoto(id int, paths string) error {
	logrus.Info(fmt.Sprintf("change aut photo id :%v", id))
	query := fmt.Sprintf(`update  authors set author_photo = $1 where id = %2`)
	if err := r.db.QueryRow(query, paths, id); err != nil {
		return fmt.Errorf("error change aut photo db")
	}
	return nil
}

//Photo

func (r *BookPostgres) BookPhoto(idPhoto int) (path, name string, err error) {
	logrus.Info(fmt.Sprintf("get book photo: %v", idPhoto))
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
	logrus.Info(fmt.Sprintf("join book-photo:id %s", idBook))
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
	logrus.Info(fmt.Sprintf("change book photo: %s , %s", newPath, newName))
	query := fmt.Sprintf(`update books_photo set foto_name = $2, book_photo =$3 where book_photo = $1`)
	r.db.QueryRow(query, paths, newName, newPath)
}

//Defect

func (r *BookPostgres) JoinDefetBookPhoto(idBook int, defect string, paths []string) error {
	logrus.Info(fmt.Sprintf("join defect photo id:%v", idBook))
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
	logrus.Info(fmt.Sprintf("get defect photo id:%v", idPhoto))
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

func (r *BookPostgres) GetByTitle(title string) (models.ReturnBook, error) {
	logrus.Info(fmt.Sprintf("find title:%s", title))

	query := fmt.Sprintf(`SELECT b.id, b.book_title, b.year_of_published,b.inventory_count, g.genre
								FROM books b
								left join book_genre bg on bg.book_id  = b.id 
								join genres g on bg.genre_id  = g.id 
								where b.book_title  = $1`)
	row, _ := r.db.Query(query, title)
	defer row.Close()
	var gen []string
	var books models.ReturnBook
	for row.Next() {
		var g string
		if err := row.Scan(&books.Id, &books.BookTitle, &books.YearOfPublished, &books.InventoryCount, &g); err != nil {
			return books, err
		}
		gen = append(gen, g)

	}
	books.Genre = gen

	return books, nil

}
func (r *BookPostgres) GetByWord(word string) ([]string, error) {
	logrus.Info(fmt.Sprintf("find book word:%s", word))
	var s []string
	query := fmt.Sprintf(`SELECT book_title FROM books WHERE book_title LIKE $1 `)
	r.db.Select(&s, query, word+"%")
	return s, nil
}
func (r *BookPostgres) ChekRegisterUser(email string) (models.User, error) {
	logrus.Info(fmt.Sprintf("find user email: %s", email))
	user := models.User{}
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
	logrus.Info("find top rating book")
	var top []models.TopRating
	query := fmt.Sprintf(`select b.id , b.book_title, b.rating from books b order by rating desc limit  3`)
	r.db.Select(&top, query)
	return top
}
func (r *BookPostgres) GetBookById(bookId int) (models.ReturnBook, error) {
	logrus.Info(fmt.Sprintf("find book, id :%v", bookId))
	var books models.ReturnBook

	query := fmt.Sprintf(`SELECT b.id, b.book_title, b.year_of_published,b.inventory_count, g.genre
								FROM books b
								left join book_genre bg on bg.book_id  = b.id 
								join genres g on bg.genre_id  = g.id 
								where b.id = $1`)
	row, _ := r.db.Query(query, bookId)
	var gen []string
	defer row.Close()
	for row.Next() {
		var g string
		if err := row.Scan(&books.Id, &books.BookTitle, &books.YearOfPublished, &books.InventoryCount, &g); err != nil {
			return books, err
		}
		gen = append(gen, g)
	}
	books.Genre = gen

	return books, nil
}

//Copies

func (r *BookPostgres) GetAllCopies(page, limit string) ([]models.CopiesInfo, error) {
	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)
	pageForSql := (p - 1) * 5

	query := fmt.Sprintf(`select 
								 bc.id, b.book_title, b.id, bc.defect, bc.defect_photo   from book_copies bc
								join books b  on bc.books_id  = b.id 
								order by b.id   
								
								LIMIT %d OFFSET %d`, l, pageForSql)

	var copies []models.CopiesInfo
	row, err := r.db.Query(query)
	if err != nil {
		return copies, err
	}
	defer row.Close()
	for row.Next() {
		b := models.CopiesInfo{}
		if err = row.Scan(&b.CopId, &b.BookTitle, &b.BookId, &b.Defect, &b.DefectPhoto); err != nil {
			return copies, err
		}

		copies = append(copies, b)
	}
	return copies, nil
}

func (r *BookPostgres) GetCopieByBookIdAndId(bookId, copId string) (models.CopiesSolo, error) {
	logrus.Info(fmt.Sprintf("get copies, id %v, book id %v", copId, bookId))
	cop, _ := strconv.Atoi(copId)
	bkId, _ := strconv.Atoi(bookId)
	var copie models.CopiesSolo
	query := fmt.Sprintf(`select *
								from book_copies bc
								where bc.id = $1`)
	r.db.QueryRow(query, cop, bkId).Scan(&copie.CopId, &copie.BookId, &copie.HideBook, &copie.Defect)

	return copie, fmt.Errorf("check get copies")
}
