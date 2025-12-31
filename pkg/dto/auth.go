package dto

import "github.com/dgrijalva/jwt-go"

type AuthRequestDto struct {
	Username string `json:"username" binding:"required,valid_username"`
	Password string `json:"password" binding:"required,valid_password"`
}

type AuthResponseDto struct {
	AccessToken string `json:"accessToken" binding:"required"`
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

type TokenClaims struct {
	jwt.StandardClaims
	TokenPayload
}

type TokenPayload struct {
	User UserDtoType
}
