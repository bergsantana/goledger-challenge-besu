package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func SetupRoutes(db *sql.DB, contract *ContractClient) {
	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/set this contract ")
		fmt.Println(contract.contractAbi)
		val := r.URL.Query().Get("value")
		num, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			http.Error(w, "Invalid value", http.StatusBadRequest)
			return
		}
		_, err = contract.SetValue(num)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success"})
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/get this contract ")
		fmt.Println(contract.contractAbi)
		val, err := contract.GetValue()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"value": val.Int64()})
	})

	http.HandleFunc("/sync", func(w http.ResponseWriter, r *http.Request) {
		val, err := contract.GetValue()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_, err = db.Exec("INSERT INTO storage (value) VALUES ($1)", val.Int64())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"synced": true, "value": val.Int64()})
	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		var dbVal int
		err := db.QueryRow("SELECT value FROM storage ORDER BY id DESC LIMIT 1").Scan(&dbVal)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		blockVal, err := contract.GetValue()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		match := dbVal == int(blockVal.Int64())
		json.NewEncoder(w).Encode(map[string]bool{"match": match})
	})
}
