package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../environment/environment.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	mux := Routes()
	server := NewServer(mux)
	server.Run()
}
