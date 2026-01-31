package riot

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pithandev/get-league-matchs/internal/domain"
)

func (c *Client) FetchMatchIDs(puuid string) ([]string, error) {
	url := fmt.Sprintf(
		"https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids?count=20",
		puuid,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", c.APIKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ids []string
	err = json.NewDecoder(resp.Body).Decode(&ids)
	return ids, err
}

func (c *Client) FetchMatchDetails(id string) (domain.MatchResponse, error) {
	url := "https://americas.api.riotgames.com/lol/match/v5/matches/" + id

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", c.APIKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return domain.MatchResponse{}, err
	}
	defer resp.Body.Close()

	var match domain.MatchResponse
	err = json.NewDecoder(resp.Body).Decode(&match)
	return match, err
}
