package main

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
		
	config.InitEnv()
	config.ConnectDatabase()
	config.InitMidtrans()

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run()
}