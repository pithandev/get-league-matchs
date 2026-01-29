package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pithandev/get-league-matchs/internal/handler"
)

func main() {
	godotenv.Load()
	fmt.Println("RIOT_API_KEY =", os.Getenv("RIOT_API_KEY"))

	http.HandleFunc("/stats", handler.StatsHandler)
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)

}
