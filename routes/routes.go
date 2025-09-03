package routes

import (
	authcontroller "MentorIT-Backend/controllers/authController"
	classcategorycontroller "MentorIT-Backend/controllers/classCategoryController"
	"MentorIT-Backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	apiRoutes := r.Group("/api")

	// Auth routes
	authRoutes := apiRoutes.Group("/auth")

	authRoutes.POST("/register", authcontroller.Register)
	authRoutes.POST("/login", authcontroller.Login)
	authRoutes.POST("/refresh-token", authcontroller.RefreshToken)
	authRoutes.GET("/logout", authcontroller.Logout)

	// Class Routes
	classRoutes := apiRoutes.Group("/classes")

	// Class category routes
	classRoutes.GET("/category", classcategorycontroller.Index)
	classRoutes.GET("/category/:id", classcategorycontroller.Show)
	
	// Admin only
	classRoutes.POST("/category", middleware.AuthMiddleware("admin"), classcategorycontroller.Create)
	classRoutes.PUT("/category/:id", middleware.AuthMiddleware("admin"), classcategorycontroller.Update)
	classRoutes.DELETE("/category/:id", middleware.AuthMiddleware("admin"), classcategorycontroller.Delete)

	
}
