package services

import (
	"fmt"
	library_app "github.com/Nkez/lib-app.git"
	"github.com/Nkez/lib-app.git/pkg/repository"
	"github.com/elliotchance/pie/pie"
	"github.com/sirupsen/logrus"
	"net/smtp"
	"time"
)

type CartService struct {
	repository repository.Cart
}

func NewCartService(repository repository.Cart) *CartService {
	return &CartService{repository: repository}
}

func (s *CartService) CreateCart(email string, books []string) (library_app.OrderCart, error) {
	logrus.Info("Create Cart service")

	rtBook, _ := s.repository.CheckOrderBook(email)
	rtBook = RemoveEmptyStrings(rtBook)
	books = RemoveEmptyStrings(books)
	getUser, _ := s.GetUser(email)

	if len(rtBook) > 0 {
		return getUser, fmt.Errorf("user %s,%s email(%s) need to return %q", getUser.FirstName, getUser.LastName, email, rtBook)
	}
	var checkPrice []float64
	var checkBooks []string
	for i := 0; i < len(books); i++ {
		book, err := s.CheckBooks(books[i])
		if err != nil {
			return getUser, err
		} else {
			checkBooks = append(checkBooks, book)
		}
	}

	if !pie.Strings.AreUnique(checkBooks) {
		return getUser, fmt.Errorf("cant take same books")
	}

	for i := 0; i < len(checkBooks); i++ {
		price, err := s.GetPrice(checkBooks[i])
		if err != nil {
			return getUser, err
		}
		checkPrice = append(checkPrice, price)
	}

	checkBooks = CheckArray(checkBooks)

	totalPrice := GetPrice(checkPrice)
	timeToReturn := time.Now().UTC().AddDate(0, 3, 0)
	getUser.DateToReturn = timeToReturn
	getUser.Price = totalPrice
	getUser.Book1 = checkBooks[0]
	getUser.Book2 = checkBooks[1]
	getUser.Book3 = checkBooks[2]
	getUser.Book4 = checkBooks[3]
	getUser.Book5 = checkBooks[4]

	return s.repository.CreateCart(getUser)
}

func (s *CartService) GetPrice(book string) (priceBook float64, err error) {
	logrus.Info("CheckBooks")
	priceBook, err = s.repository.GetPrice(book)
	if err != nil {
		return priceBook, err
	}
	return priceBook, nil
}

func (s *CartService) CheckBooks(book string) (checkBook string, err error) {
	logrus.Info("CheckBooks")
	checkBook, err = s.repository.CheckBooks(book)
	if err != nil {
		return checkBook, err
	}
	return checkBook, err
}

func (s *CartService) GetUser(email string) (library_app.OrderCart, error) {
	logrus.Info("CheckBooks")
	return s.repository.GetUser(email)
}

func (s *CartService) CheckOrderBook(email string) (orderBook []string, err error) {
	logrus.Info("CheckOrderBook")
	orderBook, err = s.repository.CheckOrderBook(email)
	return orderBook, err
}

func (s *CartService) FindRegisterUser(email string) (library_app.User, error) {
	logrus.Info("Find Register User service")
	user, err := s.repository.FindRegisterUser(email)
	if err != nil {
		return user, err
	}
	return user, err
}

func (s *CartService) GetEmailToSend() ([]library_app.OrderCart, []string, error) {
	carts, _ := s.repository.GetEmailToSend()
	var emails []string
	books := ""
	for _, v := range carts {
		emails = append(emails, v.EmailAddress)
	}
	fmt.Println(books)

	return carts, emails, nil
}
func (s *CartService) UpdatePrice() error {
	if err := s.UpdatePrice(); err != nil {
		return fmt.Errorf("dd")
	}
	return s.UpdatePrice()
}

func (s *CartService) SendEmail() {
	//fmt.Println(emails)
	//data := emailParam[0].DateToReturn
	//fmt.Println(data)
	//data2 := data.AddDate(0,0,5)
	//fmt.Println(data)
	//fmt.Println(data.Day() - data2.Day())
	for {
		emailParam, _, _ := s.GetEmailToSend()
		for _, v := range emailParam {
			from := "libtestgolang@gmail.com"
			password := "libgo123456"
			host := "smtp.gmail.com"
			port := "587"
			address := host + ":" + port

			//to := []string{v.EmailAddress}
			to := []string{"kezmikita@gmail.com"}

			message := []byte(fmt.Sprintf("Книги сюда %s\n%s\n%s\n%s\n%s", v.Book1, v.Book2, v.Book3, v.Book4, v.Book5))

			auth := smtp.PlainAuth("", from, password, host)

			err := smtp.SendMail(address, auth, from, to, message)
			if err != nil {
				fmt.Println("err:", err)
				return
			}

		}
		time.Sleep(5600 * time.Second /*time.Hour*/)
	}
}
