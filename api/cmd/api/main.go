package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/bergsantana/goledger-challenge-besu/api/contract"
	"github.com/bergsantana/goledger-challenge-besu/api/database"
	"github.com/bergsantana/goledger-challenge-besu/api/handler"

	"github.com/joho/godotenv"
)

func main() {
	log.SetReportCaller(true)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	db := database.InitDB()
	defer db.Close()

	contract, err := contract.LoadContract()
	if err != nil {
		log.Fatalf("Error ao ler contrato: %v", err)
	}

	handler.SetupRoutes(db, contract)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
