package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/s-404/go-auth-example/pkg/dto"
	"net/http"
)

// DoSmth @Summary DoSmth
// @Tags ACTION
// @Description do smth as authed
// @ID create-account
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.EmptyResponse
// @Failure default "`ERR_APP`: internal error"
// @Router /api/action/doSmth [post]
func (h *Handler) DoSmth(c *gin.Context) {
	c.JSON(http.StatusOK, dto.EmptyResponse{})
}
