package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type AccountResponse struct {
	PUUID string `json:"puuid"`
}

type Metadata struct {
	MatchId      string
	Participants []string
}

type InfoParticipants struct {
	ChampionName string
	Kills        int
	Deaths       int
	Assists      int
}

type Info struct {
	GameDuration time.Duration
	Participants []InfoParticipants
}

type MatchResponse struct {
	//Metadata Metadata
	Info Info
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

	matchIDs, err := fetchMatchIDs(puuid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	start := time.Now()

	matchIdChan := produceMatchIDs(matchIDs)
	matchesChan := fetchMatchesPipeline(matchIdChan, 5)

	matches := make([]MatchResponse, 0)

	for match := range matchesChan {
		matches = append(matches, match)
	}

	fmt.Println("tempo: ", time.Since(start))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)

}

func fetchFromRiot(summoner string) ([]byte, error) {
	apiKey := os.Getenv("RIOT_API_KEY")

	parts := strings.Split(summoner, "#")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid summoner format, use name#tag")
	}

	gameName := url.PathEscape(parts[0])
	tagLine := url.PathEscape(parts[1])

	url := fmt.Sprintf(
		"https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s",
		gameName,
		tagLine,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println("SUMMONER RECEBIDO:", summoner)
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

func fetchMatchDetails(id string) (MatchResponse, error) {
	apiKey := os.Getenv("RIOT_API_KEY")
	escapedId := url.PathEscape(id)
	url := "https://americas.api.riotgames.com/lol/match/v5/matches/" + escapedId

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return MatchResponse{}, err
	}
	req.Header.Set("X-Riot-Token", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return MatchResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return MatchResponse{}, fmt.Errorf("riot api error: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var matchInfo MatchResponse

	err = json.NewDecoder(resp.Body).Decode(&matchInfo)

	return matchInfo, nil
}

func matchWorker(wg *sync.WaitGroup, jobs <-chan string, results chan<- MatchResponse) {
	defer wg.Done()

	for matchId := range jobs {
		match, err := fetchMatchDetails(matchId)
		if err != nil {
			continue
		}

		match.Info.GameDuration *= time.Second

		results <- match
	}
}

func produceMatchIDs(matchIDs []string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for _, id := range matchIDs {
			out <- id
		}
	}()

	return out

}

func fetchMatchDetailsWorker(wg *sync.WaitGroup, jobs <-chan string, results chan<- MatchResponse) {
	defer wg.Done()

	for matchID := range jobs {

		match, err := fetchMatchDetails(matchID)
		if err != nil {
			continue
		}

		match.Info.GameDuration *= time.Second
		results <- match

	}
}

func fetchMatchesPipeline(matchIDs <-chan string, workerCount int) <-chan MatchResponse {
	results := make(chan MatchResponse)
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go fetchMatchDetailsWorker(&wg, matchIDs, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
