package services

import (
	"github.com/Nkez/library-app.git/models"
	"github.com/Nkez/library-app.git/pkg/repository"
	"github.com/asaskevich/govalidator"
	"github.com/sirupsen/logrus"
	"time"
)

type UserService struct {
	repository repository.User
}

func NewUserService(repository repository.User) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) CreateUser(user models.User) (id int, err error) {
	logrus.Info("Create User service")
	res, err := govalidator.ValidateStruct(user)
	if err != nil {
		println("error: " + err.Error())
	}
	if res == false {
		return 0, err
	}
	time.Parse("2006-01-02 :", user.Birthday)
	id, _ = s.repository.CreateUser(user)
	return id, err
}

func (s *UserService) GetByName(name string) []models.UserName {
	return s.repository.GetByName(name)
}

func (s *UserService) FindByWord(name string) []models.UserName {
	return s.repository.FindByWord(name)
}
func (s *UserService) GetAllUsers(page, limit string) ([]models.User, error) {
	logrus.Info("Find Users service")
	return s.repository.GetAllUsers(page, limit)
}
