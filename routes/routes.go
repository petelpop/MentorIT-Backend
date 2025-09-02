package routes

import (
	"MentorIT-Backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	apiRoutes := r.Group("/api")
	
	// Auth routes
	authRoutes := apiRoutes.Group("/auth")

	authRoutes.POST("/register", controllers.Register)
	authRoutes.POST("/login", controllers.Login)
	authRoutes.POST("/refresh-token", controllers.RefreshToken)
	authRoutes.GET("/logout", controllers.Logout)

	
}