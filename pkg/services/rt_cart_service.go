package services

import (
	"github.com/Nkez/lib-app.git/pkg/repository"
)

type ReturnCartService struct {
	repository repository.ReturnCart
}

func NewReturnCartService(repository repository.ReturnCart) *ReturnCartService {
	return &ReturnCartService{repository: repository}
}

func (s *ReturnCartService) CreateRtCart() error {
	return s.repository.CreateRtCart()
}
