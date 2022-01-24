package services

import (
	library_app "github.com/Nkez/lib-app.git/models"
	"github.com/Nkez/lib-app.git/pkg/repository"
	"github.com/sirupsen/logrus"
	"time"
)

type UserService struct {
	repository repository.User
}

func NewUserService(repository repository.User) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) CreateUser(user library_app.User) (int, error) {
	logrus.Info("Create User service")
	time.Parse("2006-01-02 :", user.Birthday)
	return s.repository.CreateUser(user)
}

func (s *UserService) GetAllUsers() ([]library_app.User, error) {
	logrus.Info("Find Users service")
	return s.repository.GetAllUsers()
}
