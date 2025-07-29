package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	//"os"

	//
	"github.com/bergsantana/goledger-challenge-besu/api"

	"github.com/joho/godotenv"
)

func main() {
	log.SetReportCaller(true)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	db := api.InitDB()
	defer db.Close()

	contract, err := api.LoadContract()
	if err != nil {
		log.Fatal(err)
	}

	api.SetupRoutes(db, contract)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
