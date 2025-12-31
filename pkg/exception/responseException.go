package exception

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ResponseException(c *gin.Context, exceptionError *ApiError) {
	logrus.Errorf("ResponseException message:%s, code:%s", exceptionError.Message, exceptionError.Code)
	c.AbortWithStatusJSON(exceptionError.Status, &exceptionError)
}
