package config

import (
	"MentorIT-Backend/models"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase()  {
	if err:= godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}



	dsn := fmt.Sprintf("%s:@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, host, port, name)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	database.AutoMigrate(
		&models.User{},
		&models.Token{},
		&models.ClassCategory{},
		&models.Class{},
		&models.Transaction{},
		&models.Class{},
		&models.Module{},
		&models.ModuleItem{},
		&models.SubModule{},
		&models.Quiz{},
		&models.QuizQuestion{},
		&models.FinalProject{},
		&models.ProjectPage{},
		&models.ResetToken{},
	)

	DB = database
}