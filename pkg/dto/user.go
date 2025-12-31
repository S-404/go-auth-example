package dto

import "github.com/s-404/go-auth-example/pkg/entity"

type UserDtoType struct {
	Guid     string `json:"guid" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func UserDto(user entity.User) *UserDtoType {
	return &UserDtoType{
		Guid:     user.Guid,
		Username: user.Username,
	}
}
