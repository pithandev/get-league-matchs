package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/pithandev/get-league-matchs/internal/handler"
	"github.com/pithandev/get-league-matchs/internal/riot"
)

func main() {
	_ = godotenv.Load()

	riotClient := riot.NewClient()

	http.HandleFunc("/stats", handler.StatsHandler(riotClient))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
