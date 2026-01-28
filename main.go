package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type AccountResponse struct {
	PUUID string `json:"puuid"`
}

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

	var account AccountResponse

	err = json.Unmarshal(data, &account)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	puuid := account.PUUID

	matchIds, err := fetchMatchIDs(puuid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(matchIds)
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

func fetchMatchIDs(puuid string) ([]string, error) {
	apiKey := os.Getenv("RIOT_API_KEY")

	escapedPUUID := url.PathEscape(puuid)

	url := "https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/" +
		escapedPUUID + "/ids?start=0&count=20"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("riot api error: status %d", resp.StatusCode)
	}

	var matchIDs []string

	err = json.NewDecoder(resp.Body).Decode(&matchIDs)

	return matchIDs, nil

}
