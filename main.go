package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/stats", statsHandler)
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)

}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	summoner := r.URL.Query().Get("summoner")
	if summoner == "" {
		http.Error(w, "missing summoner", http.StatusBadRequest)
		return
	}

	data, err := fetchFromRiot(summoner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)

}

func fetchFromRiot(summoner string) ([]byte, error) {
	apiKey := os.Getenv("RIOT_API_KEY")
	url := "https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/LAP Pithanjng/LAPID"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
