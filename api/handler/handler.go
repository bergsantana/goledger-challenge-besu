package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bergsantana/goledger-challenge-besu/api/contract"

	log "github.com/sirupsen/logrus"
)

// Defines the HTTP API routes
func SetupRoutes(db *sql.DB, contract *contract.ContractClient) {
	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {

		val := r.URL.Query().Get("value")
		num, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			log.Errorf("Invalid value: '%v'\n", val)
			http.Error(w, "Invalid value", http.StatusBadRequest)
			return
		}
		_, err = contract.SetValue(num)
		if err != nil {
			log.Errorf("Could not set value to contract: '%v'\n", err)
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success"})
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {

		val, _, err := contract.GetValue()
		if err != nil {
			log.Errorf("Could get the value of contract: '%v'\n", err)
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"value": val.Int64()})
	})

	http.HandleFunc("/sync", func(w http.ResponseWriter, r *http.Request) {
		val, address, err := contract.GetValue()
		if err != nil {
			log.Errorf("Error getting the value of contract: '%v'\n", err)
			http.Error(w, err.Error(), 500)
			return
		}
		_, err = db.Exec(`
			INSERT INTO storage (value, address)
			VALUES ($1, $2)
			ON CONFLICT (address) DO UPDATE
			SET value = EXCLUDED.value
		`, val.Int64(), address.Hex())

		if err != nil {
			log.Errorf("Error while sync value to database: '%v'\n", err)
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"synced": true, "value": val.Int64()})
	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		blockVal, address, err := contract.GetValue()
		if err != nil {
			log.Errorf("Error getting the value of contract: '%v'\n", err)

			http.Error(w, err.Error(), 500)
			return
		}

		var dbVal int
		err = db.QueryRow("SELECT value FROM storage WHERE address= $1 ORDER BY id DESC LIMIT 1", address.Hex()).Scan(&dbVal)
		if err != nil {
			log.Errorf("Error while querying value from database: '%v'\n", err)
			http.Error(w, err.Error(), 500)
			return
		}

		match := dbVal == int(blockVal.Int64())
		json.NewEncoder(w).Encode(match)
	})
}
