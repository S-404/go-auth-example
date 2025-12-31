package service

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/s-404/go-auth-example/pkg/dto"
	"github.com/s-404/go-auth-example/pkg/entity"
	"github.com/s-404/go-auth-example/pkg/exception"
	"github.com/s-404/go-auth-example/pkg/repository"
	"github.com/spf13/viper"
)

type IAuthService interface {
	Auth(user dto.UserDtoType) (*dto.AuthTokens, *exception.ApiError)
	RefreshAccessToken(refreshToken string) (*dto.AuthTokens, *exception.ApiError)
	AuthCookie(refreshToken string) http.Cookie
	ParseToken(accessToken string) (*dto.TokenPayload, error)
}

type AuthService struct {
	userRepository  repository.IUserRepository
	tokenRepository repository.ITokenRepository
}

func NewAuthService(
	userRepository repository.IUserRepository,
	tokenRepository repository.ITokenRepository,
) *AuthService {
	return &AuthService{userRepository, tokenRepository}
}

func (s *AuthService) Auth(user dto.UserDtoType) (*dto.AuthTokens, *exception.ApiError) {
	return GenerateAuthTokens(user)
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (*dto.AuthTokens, *exception.ApiError) {
	_, err := s.tokenRepository.FindByToken(refreshToken)
	if err != nil {
		return nil, exception.AuthError("Invalid Token")
	}

	parsedTokenPayload, err := s.ParseToken(refreshToken)
	if err != nil {
		return nil, exception.AuthError("Invalid Token")
	}

	user, err := s.userRepository.FindByUsername(parsedTokenPayload.User.Username)
	if err != nil {
		return nil, exception.AuthError("Invalid Token")
	}

	generatedTokens, ApiError := GenerateAuthTokens(dto.UserDtoType{
		Guid:     user.Guid,
		Username: user.Username,
	})
	if ApiError != nil {
		return nil, ApiError
	}

	_, createErr := s.tokenRepository.Create(entity.Token{
		UserGuid: user.Guid,
		Token:    generatedTokens.RefreshToken,
	})
	if createErr != nil {
		return nil, exception.InternalError("Create token failed")
	}

	return generatedTokens, nil
}

func (s *AuthService) AuthCookie(refreshToken string) http.Cookie {
	return http.Cookie{
		Name:     "token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   viper.GetBool("COOKIE_SECURE"),
		MaxAge:   viper.GetInt("COOKIE_MAX_AGE_SECONDS"),
	}
}

func (s *AuthService) ParseToken(accessToken string) (*dto.TokenPayload, error) {
	token, err := jwt.ParseWithClaims(accessToken, &dto.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(viper.GetString("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*dto.TokenClaims)
	if !ok {
		return nil, errors.New("token claims are not of type *dto.TokenClaims")
	}

	return &dto.TokenPayload{
		User: dto.UserDtoType{
			Guid:     claims.User.Guid,
			Username: claims.User.Username,
		},
	}, nil
}

func GenerateAuthTokens(user dto.UserDtoType) (*dto.AuthTokens, *exception.ApiError) {
	accessToken, err := GenerateToken(user, viper.GetInt("JWT_ACCESS_TOKEN_EXPIRATION_SECONDS"))
	if err != nil {
		return nil, exception.InternalError("failed generate access token")
	}

	refreshToken, err := GenerateToken(user, viper.GetInt("JWT_REFRESH_TOKEN_EXPIRATION_SECONDS"))
	if err != nil {
		return nil, exception.InternalError("failed generate refresh token")
	}

	return &dto.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func GenerateToken(user dto.UserDtoType, duration int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &dto.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   "username",
			ExpiresAt: time.Now().Add(time.Second * time.Duration(duration)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		TokenPayload: dto.TokenPayload{
			User: user,
		},
	})

	return token.SignedString([]byte(viper.GetString("JWT_SECRET_KEY")))
}
