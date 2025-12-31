package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/s-404/go-auth-example/pkg/handler/validation"
	"github.com/s-404/go-auth-example/pkg/service"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/s-404/go-auth-example/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation(validation.ValidUsername, validation.UsernameValidation)
		v.RegisterValidation(validation.ValidPassword, validation.PasswordValidation)
	}

	auth := router.Group("api/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
		auth.POST("/refresh", h.RefreshAccessToken)
	}

	api := router.Group("api", h.authMiddleware)
	{
		action := api.Group("action")
		{
			action.POST("/doSmth", h.DoSmth)
		}
	}

	return router
}
