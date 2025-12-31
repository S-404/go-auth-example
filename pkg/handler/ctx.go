package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/s-404/go-auth-example/pkg/dto"
)

const (
	CtxUser = "user"
)

func GetCtxUser(c *gin.Context) (*dto.UserDtoType, error) {
	user, ok := c.Get(CtxUser)
	if !ok {
		return nil, errors.New("ctx user not found")
	}

	userData, ok := user.(dto.UserDtoType)
	if !ok {
		return nil, errors.New("ctx user is of invalid type")
	}

	return &userData, nil
}
