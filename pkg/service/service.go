package service

import "github.com/s-404/go-auth-example/pkg/repository"

type Service struct {
	Auth  IAuthService
	Token ITokenService
	User  IUserService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Auth:  NewAuthService(repos.User, repos.Token),
		Token: NewTokenService(repos.Token),
		User:  NewUserService(repos.User),
	}
}
