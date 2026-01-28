package routes

import (
	"github.com/ericolvr/sec-backend/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	userHandler *api.UserHandler
}

func NewUserRoutes(userHandler *api.UserHandler) *UserRoutes {
	return &UserRoutes{
		userHandler: userHandler,
	}
}

func (ur *UserRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	users := v1.Group("/users")
	{
		users.POST("", ur.userHandler.Create)
		users.GET("", ur.userHandler.List)
		users.GET("/:id", ur.userHandler.GetByID)
		users.GET("/mobile/:mobile", ur.userHandler.GetByMobile)
		users.PUT("/:id", ur.userHandler.Update)
		users.DELETE("/:id", ur.userHandler.Delete)
	}
}
