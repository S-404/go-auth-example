package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/s-404/go-auth-example/pkg/dto"
	"github.com/s-404/go-auth-example/pkg/entity"
	"github.com/s-404/go-auth-example/pkg/exception"
)

// Register @Summary Register
// @Tags AUTH
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param credentials body dto.AuthRequestDto true "credentials"
// @Success 200 {object} dto.EmptyResponse
// @Failure 409 "`ERR_CLIENT_ENTITY_ALREADY_EXIST`: username already used"
// @Failure 400 "`ERR_CLIENT_REQUEST_VALIDATION`: validation error"
// @Failure default "`ERR_APP`: internal error"
// @Router /api/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var input dto.AuthRequestDto
	err := c.BindJSON(&input)
	if err != nil {
		exception.ResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	_, ApiError := h.services.User.Create(input)
	if ApiError != nil {
		exception.ResponseException(c, ApiError)
		return
	}

	c.JSON(http.StatusOK, dto.EmptyResponse{})
}

// Login @Summary Login
// @Tags AUTH
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param credentials body dto.AuthRequestDto true "credentials"
// @Success 200 {object} dto.AuthResponseDto
// @Failure 400 "`ERR_CLIENT_REQUEST_VALIDATION`: validation error, `ERR_CLIENT_BAD_REQUEST`: wrong credentials"
// @Failure default "`ERR_APP`: internal error"
// @Router /api/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var input dto.AuthRequestDto
	err := c.BindJSON(&input)
	if err != nil {
		exception.ResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	user, ApiError := h.services.User.GetByCredentials(input.Username, input.Password)
	if ApiError != nil {
		exception.ResponseException(c, exception.BadRequest("wrong credentials"))
		return
	}

	generatedToken, ApiError := h.services.Auth.Auth(dto.UserDtoType{
		Guid:     user.Guid,
		Username: user.Username,
	})
	if ApiError != nil {
		exception.ResponseException(c, ApiError)
		return
	}

	_, ApiError = h.services.Token.Create(entity.Token{
		UserGuid: user.Guid,
		Token:    generatedToken.RefreshToken,
	})
	if ApiError != nil {
		exception.ResponseException(c, ApiError)
		return
	}

	cookie := h.services.Auth.AuthCookie(generatedToken.RefreshToken)
	c.SetCookie(
		cookie.Name,
		cookie.Value,
		cookie.MaxAge,
		cookie.Path,
		cookie.Domain,
		cookie.Secure,
		cookie.HttpOnly,
	)

	c.JSON(http.StatusOK, dto.AuthResponseDto{
		AccessToken: generatedToken.AccessToken,
	})
}

// Logout @Summary Logout
// @Tags AUTH
// @Description logout
// @ID logout
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.EmptyResponse
// @Failure 401 "`ERR_CLIENT_AUTH`: unauthed"
// @Router /api/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var token, err = c.Cookie("token")
	if err != nil {
		exception.ResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	h.services.Token.DestroyByToken(token)

	c.JSON(http.StatusOK, dto.EmptyResponse{})
}

// RefreshAccessToken @Summary Refresh access token
// @Tags AUTH
// @Description refresh access token
// @ID refreshAccessToken
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.AuthResponseDto
// @Failure 401 "`ERR_CLIENT_AUTH`: unauthed"
// @Failure default "`ERR_APP`: internal error"
// @Router /api/auth/refresh [post]
func (h *Handler) RefreshAccessToken(c *gin.Context) {
	var token, err = c.Cookie("token")
	if err != nil {
		exception.ResponseException(c, exception.AuthError("Unauthed"))
		return
	}
	generatedToken, ApiError := h.services.Auth.RefreshAccessToken(token)
	if ApiError != nil {
		exception.ResponseException(c, ApiError)
		return
	}

	cookie := h.services.Auth.AuthCookie(generatedToken.RefreshToken)
	c.SetCookie(
		cookie.Name,
		cookie.Value,
		cookie.MaxAge,
		cookie.Path,
		cookie.Domain,
		cookie.Secure,
		cookie.HttpOnly,
	)

	c.JSON(http.StatusOK, dto.AuthResponseDto{
		AccessToken: generatedToken.AccessToken,
	})
}
