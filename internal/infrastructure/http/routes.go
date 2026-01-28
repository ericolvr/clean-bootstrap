package http

import (
	"github.com/gin-gonic/gin"
)

type RouteSetup interface {
	SetupRoutes(v1 *gin.RouterGroup)
}

type Router struct {
	userRoutes RouteSetup
	// Add other route groups here as needed
	// authRoutes RouteSetup
	// propertyRoutes RouteSetup
}

func NewRouter(userRoutes RouteSetup) *Router {
	return &Router{
		userRoutes: userRoutes,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	engine.Static("/uploads", "./uploads")

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 group
	v1 := engine.Group("/api/v1")
	{
		r.userRoutes.SetupRoutes(v1)

		// Add other route groups here
		// r.authRoutes.SetupRoutes(v1)
		// r.propertyRoutes.SetupRoutes(v1)
	}
}
