package service

import (
	"fmt"

	"github.com/s-404/go-auth-example/pkg/entity"
	"github.com/s-404/go-auth-example/pkg/exception"
	"github.com/s-404/go-auth-example/pkg/repository"
)

type ITokenService interface {
	Create(token entity.Token) (*entity.Token, *exception.ApiError)
	DestroyByToken(token string)
}

type TokenService struct {
	tokenRepository repository.ITokenRepository
}

func NewTokenService(tokenRepository repository.ITokenRepository) *TokenService {
	return &TokenService{tokenRepository}
}

func (s *TokenService) Create(token entity.Token) (*entity.Token, *exception.ApiError) {
	createdToken, err := s.tokenRepository.Create(token)
	if err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed token creating with userGuid '%s'. error: %s", token.UserGuid, err.Error()))
	}

	return createdToken, nil
}

func (s *TokenService) DestroyByToken(token string) {
	s.tokenRepository.DestroyByToken(token)
}
