package main

import "MentorIT-Backend/config"

func main() {
	config.InitEnv()
	config.ConnectDatabase()

}