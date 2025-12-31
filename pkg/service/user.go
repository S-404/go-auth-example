package service

import (
	"fmt"

	"github.com/s-404/go-auth-example/pkg/dto"
	"github.com/s-404/go-auth-example/pkg/entity"
	"github.com/s-404/go-auth-example/pkg/exception"
	"github.com/s-404/go-auth-example/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Create(user dto.AuthRequestDto) (*entity.User, *exception.ApiError)
	GetByUsername(username string) (*entity.User, *exception.ApiError)
	GetByCredentials(username string, password string) (*entity.User, *exception.ApiError)
}

type UserService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) *UserService {
	return &UserService{userRepository}
}

func (s *UserService) Create(userData dto.AuthRequestDto) (*entity.User, *exception.ApiError) {
	candidate, err := s.userRepository.FindByUsername(userData.Username)
	if err == nil && candidate.Username == userData.Username {
		return nil, exception.EntityAlreadyExistsError("user", fmt.Sprintf("username:%s", candidate.Username))
	}

	hashedPassword, err := HashPassword(userData.Password)
	if err != nil {
		return nil, exception.InternalError("failed password hashing: create user")
	}

	createdUser, err := s.userRepository.Create(entity.User{
		Username: userData.Username,
		Password: hashedPassword,
	})
	if err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed user creating with username '%s'", userData.Username))
	}

	return createdUser, nil
}

func (s *UserService) GetByUsername(username string) (*entity.User, *exception.ApiError) {
	user, err := s.userRepository.FindByUsername(username)
	if err != nil {
		return nil, exception.EntityNotFoundError("user", fmt.Sprintf("username:%s", username))
	}

	return user, nil
}

func (s *UserService) GetByCredentials(username string, password string) (*entity.User, *exception.ApiError) {
	user, ApiErr := s.GetByUsername(username)
	if ApiErr != nil {
		return nil, ApiErr
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if compareErr != nil {
		return nil, exception.BadRequest("bad credentials")
	}

	return user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
