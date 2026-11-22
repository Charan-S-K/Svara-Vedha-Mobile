package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"mobile/database"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

//go:embed index.html
var indexHTML []byte

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
}

func checker(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == "GET" {
		var name, day, period string
		query := r.URL.Query()
		name = query.Get("name")
		day = query.Get("day")
		period = query.Get("period")
		database.IsAccepted = false
		database.Check(name, day, period)
		if database.IsAccepted && database.IsPresent {
			fmt.Println("Is avaible ", database.IsAccepted)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(database.IsAccepted)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(database.IsAccepted)
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("GET method only, invalid methos used")
	}
}

func GetData(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		database.DataProvider()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(database.Names)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("GET method only, invalid methos used")
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		http.NotFound(w, r)
		return
	}
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("GET method only, invalid methos used")
		return
	}
	// Retain the teacher-name JSON functionality at / if client requests JSON
	if strings.Contains(r.Header.Get("Accept"), "application/json") || r.URL.Query().Get("format") == "json" {
		GetData(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if data, err := os.ReadFile("index.html"); err == nil {
		w.Write(data)
		return
	}
	w.Write(indexHTML)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("failed to load env")
		return
	}
	port := os.Getenv("PORT")
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/teachers", GetData)
	http.HandleFunc("/check", checker)
	http.ListenAndServe(port, nil)
}
