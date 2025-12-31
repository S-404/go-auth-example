package handler

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/s-404/go-auth-example/pkg/exception"
)

const (
	AUTHORIZATION_HEADER = "Authorization"
)

func (h *Handler) authMiddleware(c *gin.Context) {
	header := c.GetHeader(AUTHORIZATION_HEADER)
	if header == "" {
		exception.ResponseException(c, exception.AuthError("Empty Authorization Header"))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		exception.ResponseException(c, exception.AuthError("Invalid Authorization Header"))
		return
	}

	if len(headerParts[1]) == 0 {
		exception.ResponseException(c, exception.AuthError("Token is empty"))
		return
	}

	parsedTokenPayload, err := h.services.Auth.ParseToken(headerParts[1])
	if err != nil {
		exception.ResponseException(c, exception.AuthError(fmt.Sprintf("Invalid token: %s", err.Error())))
		return
	}

	c.Set(CtxUser, parsedTokenPayload.User)
}
