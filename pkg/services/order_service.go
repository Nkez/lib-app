package services

import (
	"errors"
	"fmt"
	"github.com/Nkez/library-app.git/models"
	"github.com/Nkez/library-app.git/pkg/repository"
	"github.com/asaskevich/govalidator"
	"github.com/elliotchance/pie/pie"
	"github.com/sirupsen/logrus"
	"math"
	"strings"
	"time"
)

type OrderService struct {
	repositoryOrder repository.Order
}

func NewOrderService(repositoryOrder repository.Order) *OrderService {
	return &OrderService{repositoryOrder: repositoryOrder}
}

func (s *OrderService) CreateOrder(input models.OrderInput) (models.Order, error) {
	logrus.Info("Creating Order")
	var order models.Order
	res, err := govalidator.ValidateStruct(input)
	if res == false {
		return order, err
	}
	logrus.Info("Check input params")
	email, books, err := s.CheckInputParams(input)
	if err != nil {
		return order, errors.New("check input params or u take more than 5 books(can only 5)")
	}
	user, err := s.FindUser(email)
	if err != nil {
		return order, errors.New(fmt.Sprintf("can not find user check email %s or create new", email))
	}
	logrus.Info("Check similar book")
	if err = s.CheckSimilaBook(books); err != nil {
		return order, err
	}

	logrus.Info("Check is return")
	_, err = s.CheckIsReturn(user.Id)
	if err != nil {
		return order, err
	}
	logrus.Info("get books id")
	idBooks, err := s.GetBooksId(books)
	if err != nil {
		return order, err
	}
	fmt.Println(idBooks)
	logrus.Info("total price")
	totalPrice, _ := s.GetPrice(idBooks)

	logrus.Info("Get  order and return day")
	orderDate, returnDate := s.Date()

	if err = s.JoinBookUser(user.Id, idBooks, totalPrice, orderDate, returnDate); err != nil {
		return order, fmt.Errorf("problem Join book user")
	}
	fmt.Println(idBooks)
	logrus.Info("Minus inventory count")

	if err = s.MinusInventoryCount(idBooks); err != nil {
		return order, err
	}
	logrus.Info("Return order cart")
	order, _ = s.ReturnOrder(user.Id)

	return order, err
}

func (s *OrderService) GetBooksId(books []string) (id []int, err error) {
	return s.repositoryOrder.GetBooksId(books)
}

func (s *OrderService) CheckInputParams(input models.OrderInput) (email string, books []string, err error) {
	email = input.EmailAddress

	if email == "" {
		return "", nil, err
	}
	for i, _ := range input.Books {
		books = append(books, input.Books[i].Book)
	}
	if len(books) > 5 {
		return "", nil, errors.New("u can take only 5 books")
	}
	return email, books, nil
}

func (s *OrderService) FindUser(email string) (models.User, error) {
	logrus.Info("CheckBooks")
	return s.repositoryOrder.FindUser(email)
}

func (s *OrderService) CheckSimilaBook(books []string) error {
	if pie.Strings.AreUnique(books) == false {
		return errors.New("cant take same books")
	}
	return nil
}

func (s *OrderService) CheckIsReturn(id int) (books []string, err error) {
	rtBooks, err := s.repositoryOrder.CheckIsReturn(id)
	if len(rtBooks) > 0 {
		return rtBooks, errors.New(fmt.Sprintf("need return this book: %s", strings.Join(rtBooks, ", ")))
	}
	return rtBooks, err
}

func (s *OrderService) GetPrice(books []int) (price float64, err error) {
	prices, err := s.repositoryOrder.GetPrice(books)
	if err != nil {
		return 0.0, err
	}
	price = s.GetTotalPrice(prices)
	return price, nil
}

func (s *OrderService) Date() (orderDate, returnDate time.Time) {
	orderDate = time.Now().UTC()
	returnDate = time.Now().AddDate(0, 3, 0)
	return orderDate, returnDate
}

func (s *OrderService) GetTotalPrice(price []float64) float64 {
	var totalPrice float64
	i := 0
	for _, value := range price {
		if value != 0.0 {
			totalPrice += value
			i++
		}
	}
	if i > 2 && i < 4 {
		totalPrice = totalPrice - (totalPrice * 0.1)
		return math.Round(totalPrice*100) / 100
	}
	if i > 4 {
		totalPrice = totalPrice - (totalPrice * 0.15)
		return math.Round(totalPrice*100) / 100
	}
	return totalPrice
}

func (s *OrderService) JoinBookUser(idUser int, idBooks []int, price float64, orderDate, returnDate time.Time) error {
	fmt.Println(idBooks)
	err := s.repositoryOrder.JoinBookUser(idUser, idBooks, price, orderDate, returnDate)
	return err
}

func (s *OrderService) MinusInventoryCount(idBooks []int) error {
	return s.repositoryOrder.MinusInventoryCount(idBooks)
}

func (s *OrderService) ReturnOrder(id int) (models.Order, error) {
	logrus.Info("Find All Books service")
	books, order, _ := s.repositoryOrder.ReturnOrder(id)
	order.Books = strings.Join(books, ", ")

	return order, nil
}

func (s *OrderService) GetAllOrder(page, limit string) ([]models.InfoOrdDept, error) {
	return s.repositoryOrder.GetAllOrder(page, limit)
}
