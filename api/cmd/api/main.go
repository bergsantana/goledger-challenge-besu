package main

import (
	"log"
	"net/http"
	//"os"

	//
	"github.com/bergsantana/goledger-challenge-besu/api"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	db := InitDB()
	defer db.Close()

	contract, err := LoadContract()
	if err != nil {
		log.Fatal(err)
	}

	SetupRoutes(db, contract)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
