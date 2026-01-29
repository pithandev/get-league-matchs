package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/pithandev/get-league-matchs/internal/domain"
)

func FetchMatchIDs(puuid string) ([]string, error) {
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
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"riot api error: status %d, body: %s", resp.StatusCode, string(body),
		)
	}

	var matchIDs []string

	err = json.NewDecoder(resp.Body).Decode(&matchIDs)

	return matchIDs, nil

}

func FetchMatchDetails(id string) (domain.MatchResponse, error) {
	apiKey := os.Getenv("RIOT_API_KEY")
	escapedId := url.PathEscape(id)
	url := "https://americas.api.riotgames.com/lol/match/v5/matches/" + escapedId

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return domain.MatchResponse{}, err
	}
	req.Header.Set("X-Riot-Token", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.MatchResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return domain.MatchResponse{}, fmt.Errorf("riot api error: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var matchInfo domain.MatchResponse

	err = json.NewDecoder(resp.Body).Decode(&matchInfo)

	return matchInfo, nil
}
