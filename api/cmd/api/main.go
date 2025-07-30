package main

import (
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bergsantana/goledger-challenge-besu/api/contract"
	"github.com/bergsantana/goledger-challenge-besu/api/database"
	"github.com/bergsantana/goledger-challenge-besu/api/handler"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	db := database.InitDB()
	defer db.Close()

	contract, err := contract.LoadContract()
	if err != nil {
		log.Fatalf("Error loading contract: %v\n", err)
	}

	handler.SetupRoutes(db, contract)

	printStartupLog(os.Getenv("CONTRACT_ADDRESS"), os.Getenv("PG_CONN"), "8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func printStartupLog(contractAddr string, dbConn string, port string) {
	now := time.Now().Format("2006-01-02 15:04:05")

	log.Println("─────────────────────────────────────────────")
	log.Println("🚀 GO-BESU API STARTED")
	log.Println("─────────────────────────────────────────────")
	log.Printf("📆 Started At      : %s\n", now)
	log.Printf("🏗  Contract Addr  : %s\n", contractAddr)
	log.Printf("🛢  PostgreSQL DSN : %s\n", dbConn)
	log.Printf("🛰  API Running On : http://localhost:%s\n", port)
	log.Println("─────────────────────────────────────────────")
	log.Println("💡 Endpoints:")
	log.Println("   ➤ GET    /get")
	log.Println("   ➤ GET    /set?value=123")
	log.Println("   ➤ GET    /sync")
	log.Println("   ➤ GET    /check")
	log.Println("─────────────────────────────────────────────")
	log.Println("📦 Ready to process blockchain interactions!")
}
