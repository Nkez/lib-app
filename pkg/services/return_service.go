package services

import (
	"errors"
	"fmt"
	"github.com/Nkez/library-app.git/models"
	"github.com/Nkez/library-app.git/pkg/repository"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"math"
	"net/smtp"
	"reflect"
	"sort"
	"time"
)

type ReturnService struct {
	repositoryReturn repository.Return
}
type DeborsService struct {
	repositoryDebors repository.Debors
}

func NewDeborsService(repositoryDebors repository.Debors) *DeborsService {
	return &DeborsService{repositoryDebors: repositoryDebors}
}

func NewReturnService(repositoryReturn repository.Return) *ReturnService {
	return &ReturnService{repositoryReturn: repositoryReturn}
}

func (s *ReturnService) ReturnCart(input models.ReturnInput) (returnCart models.DbrInfo, err error) {
	logrus.Info("start return cart")
	logrus.Info("get book info")
	idUser, idBooks, orderDate, err := s.GetBooksIdAndDate(input)
	if err != nil {
		return returnCart, err
	}
	logrus.Info("check returning book")
	if err := s.CheckReturningBook(idUser, idBooks); err != nil {
		return returnCart, fmt.Errorf("u cannot return this book")
	}
	fmt.Println("RETURNING BOOOk")
	rating, _ := s.GetRating(idBooks, input)
	if err != nil {
		return returnCart, err
	}
	logrus.Info("Return Books in Lib")
	returnCart = s.ReturnBook(idUser, idBooks, orderDate)
	logrus.Info("return books")
	s.ReturnBooksInLibUpdateRating(idBooks, rating)
	logrus.Info("update price")
	s.UpdatePrice(idUser, returnCart.Price)

	return returnCart, err

}

func (s *ReturnService) CheckReturningBook(id int, idBooks []int) error {
	chekBookId, _ := s.repositoryReturn.CheckReturningBook(id)
	sort.Ints(chekBookId)
	sort.Ints(idBooks)
	if reflect.DeepEqual(chekBookId, idBooks) == false {
		return errors.New("u cannot return this book")
	}
	return nil
}

func (s *ReturnService) GetBooksIdAndDate(input models.ReturnInput) (idUser int, idBooks []int, orderDate time.Time, err error) {
	idUser = input.IdUser
	fmt.Println(input.ReturnDay)
	if input.ReturnDay == "" {
		orderDate = time.Now().UTC()
	} else {
		orderDate, _ = time.Parse("2006-01-02", input.ReturnDay)
	}

	for _, st := range input.ReturnCart {
		idBooks = append(idBooks, st.IdBook)
	}

	return idUser, idBooks, orderDate, nil
}

func (s *ReturnService) UpdatePrice(id int, newPrice float64) {
	s.repositoryReturn.UpdatePrice(id, newPrice)
}

func (s *DeborsService) GetAllDebors() ([]models.InfoOrdDept, error) {
	return s.repositoryDebors.GetAllDebors()
}

func (s *ReturnService) ReturnBook(idUser int, idBook []int, returnDate time.Time) models.DbrInfo {

	dbrInfo, _ := s.repositoryReturn.ReturnBook(idUser, idBook, returnDate)

	if dbrInfo.OrderDate.After(dbrInfo.DateToReturn) {
		hh := dbrInfo.OrderDate.Sub(dbrInfo.DateToReturn).Hours()
		days := hh / 24
		priceNew := days * 0.01
		dbrInfo.Price = math.Round((dbrInfo.Price+dbrInfo.Price*priceNew)*100) / 100
	}

	return dbrInfo
}

func (s *ReturnService) ReturnBooksInLibUpdateRating(booksId []int, rating []decimal.Decimal) error {
	return s.repositoryReturn.ReturnBooksInLibUpdateRating(booksId, rating)
}

func (s *ReturnService) GetRating(idBooks []int, input models.ReturnInput) ([]decimal.Decimal, error) {
	rat, _ := s.repositoryReturn.GetRating(idBooks)
	newRat, err := s.NewRating(rat, input)
	if err != nil {
		return nil, err
	}
	return newRat, err
}

func (s *ReturnService) NewRating(ratings []float64, input models.ReturnInput) (rating []decimal.Decimal, err error) {
	minRating := decimal.NewFromFloat(-0.1)
	maxRating := decimal.NewFromFloat(5.1)
	for i := 0; i < len(ratings); i++ {
		nRating := input.ReturnCart[i].Rating
		if nRating.GreaterThan(maxRating) || nRating.LessThan(minRating) {
			return nil, errors.New("input rating from 0.0 to 5.0")
		}

		result := (decimal.Decimal.InexactFloat64(nRating) + ratings[i]) / float64(2)
		rating = append(rating, decimal.NewFromFloat(result))
	}

	fmt.Println(rating)
	return rating, err
}

///Email

func SendEmail(email string, message string) {

	// Sender data.
	from := "libtestgolang@gmail.com"
	password := "libgo123456"
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	to := []string{
		email,
	}

	m := []byte(message)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, m)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *DeborsService) CheckForSend() {
	var mes1 string
	var mes2 string

	debors, _ := s.repositoryDebors.GetAllDebors()

	for _, i := range debors {
		t1 := i.DateToReturn
		dt1, _ := time.Parse("2006-01-02", t1)
		t2 := time.Now()
		mes1 = "need return books " + i.Books
		fmt.Println(i.Books)
		if (t2.Sub(dt1).Hours())/24 <= 6 {
			SendEmail(i.EmailAddress, mes1)
		} else if (t2.Sub(dt1).Hours())/24 > 6 {
			days := (t2.Sub(dt1).Hours()) / 24
			p := i.Price.InexactFloat64() * 0.01 * days
			fmt.Println(p)
			mes2 = i.LastName + " " + i.FirstName + " " + mes1 + "new price " + "is" + fmt.Sprintf("%.2f", p)
			SendEmail(i.EmailAddress, mes2)
		}

	}
}

func (s *DeborsService) FirstCheck() {
	dbr, _ := s.repositoryDebors.GetAllDebors()
	for _, i := range dbr {
		t1 := i.DateToReturn
		dt1, _ := time.Parse("2006-01-02", t1)
		t2 := time.Now()
		if t2.After(dt1) {
			s.CheckForSend()
		}
	}
}

func (s *DeborsService) CheckDebt() bool {
	debor, _ := s.repositoryDebors.GetAllDebors()
	if len(debor) == 0 {
		return false
	}
	for _, i := range debor {
		t1 := i.DateToReturn
		dt1, _ := time.Parse("2006-01-02", t1)
		t2 := time.Now()
		if t2.After(dt1) {
			return true
		}
	}
	return false
}

func (s *DeborsService) WaitAndEmailAgain() {
	waitTime := time.NewTicker(5 * time.Minute)
	go func() {
		if true {
			for range waitTime.C {
				if s.CheckDebt() == true {
					s.CheckForSend()
					fmt.Println("Email sent...")
				} else {
					fmt.Println("No debt")
				}
			}
		}
	}()
}
