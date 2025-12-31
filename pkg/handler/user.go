package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) userCreate(c *gin.Context) {
	println("handler user create")
}

func (h *Handler) userList(c *gin.Context) {
	println("handler user list")
}

func (h *Handler) userDelete(c *gin.Context) {
	println("handler user delete")
}
