package routes

import (
	admincontroller "MentorIT-Backend/controllers/adminController"
	authcontroller "MentorIT-Backend/controllers/authController"
	classcategorycontroller "MentorIT-Backend/controllers/classCategoryController"
	classcontroller "MentorIT-Backend/controllers/classController"
	modulecontroller "MentorIT-Backend/controllers/moduleController"
	paymentcontroller "MentorIT-Backend/controllers/paymentController"
	submodulecontroller "MentorIT-Backend/controllers/subModuleController"
	"MentorIT-Backend/helper"
	"MentorIT-Backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	student := string(helper.Student)
	teacher := string(helper.Teacher)
	admin := string(helper.Admin)

	apiRoutes := r.Group("/api")

	// Auth routes
	authRoutes := apiRoutes.Group("/auth")

	authRoutes.POST("/register", authcontroller.Register)
	authRoutes.POST("/login", authcontroller.Login)
	authRoutes.POST("/refresh-token", authcontroller.RefreshToken)
	authRoutes.GET("/logout", authcontroller.Logout)

	//========================================================================================================

	// Class Routes
	classRoutes := apiRoutes.Group("/classes")

	classRoutes.GET("/class", middleware.AuthMiddleware(student, teacher, admin), classcontroller.Index)
	classRoutes.GET("/class/:id", middleware.AuthMiddleware(student, teacher, admin), classcontroller.Show)
	classRoutes.GET("/category/:id/class/", middleware.AuthMiddleware(student, teacher, admin), classcontroller.IndexByCategory)

	// Admin & teachers only ( Manage Class )
	classRoutes.POST("/create", middleware.AuthMiddleware(admin, teacher), classcontroller.Create)
	classRoutes.PUT("/update/:id", middleware.AuthMiddleware(admin, teacher), classcontroller.Update)
	classRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(admin, teacher), classcontroller.Delete)

	// Payment
	classRoutes.POST("/buy-class/:id", middleware.AuthMiddleware(student), paymentcontroller.BuyClass)
	classRoutes.POST("/buy-class/notification", middleware.AuthMiddleware(student), paymentcontroller.PaymentNotification)

	// Module
	classRoutes.GET("/modules/:id", middleware.AuthMiddleware(student, teacher), modulecontroller.GetModules)
	classRoutes.GET("/module/:id", middleware.AuthMiddleware(student, teacher), modulecontroller.GetModuleDetail)

	// Admin & teachers only ( Manage Module )
	classRoutes.POST("/create/module", middleware.AuthMiddleware(teacher, admin), modulecontroller.CreateModule)
	classRoutes.PUT("/edit/module/:id", middleware.AuthMiddleware(teacher, admin), modulecontroller.UpdateModule)
	classRoutes.DELETE("/delete/module/:id", middleware.AuthMiddleware(teacher, admin), modulecontroller.DeleteModule)

	// Sub-Module
	classRoutes.GET("/sub-modules/:id", middleware.AuthMiddleware(student, teacher), submodulecontroller.GetSubModules)
	classRoutes.GET("/sub-module/:id", middleware.AuthMiddleware(student, teacher), submodulecontroller.GetSubModuleDetail)

	// Admin & teachers only ( Manage Sub-Module )
	classRoutes.POST("/create/sub-module", middleware.AuthMiddleware(teacher, admin), submodulecontroller.CreateSubModule)
	classRoutes.PUT("/edit/sub-module/:id", middleware.AuthMiddleware(teacher, admin), submodulecontroller.UpdateSubModule)
	classRoutes.DELETE("/delete/sub-module/:id", middleware.AuthMiddleware(teacher, admin), submodulecontroller.DeleteSubModule)

	//========================================================================================================

	// Class category routes
	classRoutes.GET("/category", classcategorycontroller.Index)
	classRoutes.GET("/category/:id", classcategorycontroller.Show)

	// Admin only
	classRoutes.POST("/category", middleware.AuthMiddleware(admin), classcategorycontroller.Create)
	classRoutes.PUT("/category/:id", middleware.AuthMiddleware(admin), classcategorycontroller.Update)
	classRoutes.DELETE("/category/:id", middleware.AuthMiddleware(admin), classcategorycontroller.Delete)

	//========================================================================================================

	// Admin Routes
	adminRoutes := apiRoutes.Group("/admin")
	adminRoutes.GET("/list-teachers", middleware.AuthMiddleware(admin), admincontroller.ListTeachers)
	adminRoutes.GET("/teacher/:id", middleware.AuthMiddleware(admin), admincontroller.GetTeacherByID)
	adminRoutes.POST("/create-teacher", middleware.AuthMiddleware(admin), admincontroller.CreateTeacher)
	adminRoutes.DELETE("/delete-teacher/:id", middleware.AuthMiddleware(admin), admincontroller.DeleteTeacher)

	//========================================================================================================
}

// func WebhookRoutes(r *gin.Engine) {
// 	student := string(helper.Student)

// 	r.POST("/api/classes/buy-class/notification", middleware.AuthMiddleware(student), paymentcontroller.PaymentNotification)

// }
