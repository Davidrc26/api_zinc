package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("environment/environment.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		fmt.Println(err.Error())
	}
	mux := Routes()
	server := NewServer(mux)
	server.Run()
}
