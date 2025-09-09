package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JWTKey []byte
var host string
var port string
var user string
var name string
var midtransServerKey string

func InitEnv() {

	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	JWTKey = []byte(os.Getenv("JWT_SECRET"))
	midtransServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	user = os.Getenv("DB_USER")
	name = os.Getenv("DB_NAME")
}
