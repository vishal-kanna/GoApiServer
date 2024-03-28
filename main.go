package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Favourite struct {
	NftId   string `json:"nftid"`
	Address string `json:"address"`
	Type    string `json:"type"`
}
type Rating struct {
	NftId   string `json:"nftid"`
	Address string `json:"address"`
	Score   int32  `json:"score"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "user"
	dbname   = "postgres"
)

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL")
}
func main() {
	defer db.Close()
	http.HandleFunc("/favourite", favourite)
	http.HandleFunc("/unfavourite", unfavourite)
	http.HandleFunc("/rating", rating)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
func favourite(w http.ResponseWriter, r *http.Request) {
	// 	Create the Favourite table if it doesn't exist
	_, err := db.Exec(`
					CREATE TABLE IF NOT EXISTS Favourite6 (
						nftid TEXT NOT NULL,
						address TEXT NOT NULL,
						type TEXT NOT NULL,
						PRIMARY KEY(nftid,address)
					)
				`)
	if err != nil {
		log.Fatal(err)
	}
	var Data Favourite
	err = json.NewDecoder(r.Body).Decode(&Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM Favourite6 WHERE nftid = $1 AND address = $2", Data.NftId, Data.Address).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "This data already present in the DB", http.StatusInternalServerError)
		return
	}

	// Insert the data if it doesn't already exist
	sqlStatement := "INSERT INTO Favourite6 (nftid, address,type) VALUES ($1, $2,$3)"
	_, err = db.Exec(sqlStatement, Data.NftId, Data.Address, Data.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "User added successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func unfavourite(w http.ResponseWriter, r *http.Request) {
	var requestData Favourite
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Check if the data exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM Favourite6 WHERE nftid = $1 AND address = $2", requestData.NftId, requestData.Address).Scan(&count)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	// Data exists, perform the deletion
	sqlStatement := "DELETE FROM Favourite6 WHERE nftid = $1 AND address = $2"
	_, err = db.Exec(sqlStatement, requestData.NftId, requestData.Address)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Data deleted successfully",
		"nftid":   requestData.NftId,
		"address": requestData.Address,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func rating(w http.ResponseWriter, r *http.Request) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Rating (
		nftid TEXT NOT NULL,
		address TEXT NOT NULL,
		score TEXT NOT NULL,
		PRIMARY KEY(nftid,address)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	var Data Rating
	err = json.NewDecoder(r.Body).Decode(&Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM Rating WHERE nftid = $1 AND address = $2", Data.NftId, Data.Address).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "This rating already present in the DB", http.StatusBadRequest)
		return
	}
	sqlStatement := "INSERT INTO Rating (nftid, address,score) VALUES ($1,$2,$3)"
	_, err = db.Exec(sqlStatement, Data.NftId, Data.Address, Data.Score)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "Rating added successfully for",
		"nftid":   Data.NftId,
		"address": Data.Address,
		"score":   Data.Score,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
