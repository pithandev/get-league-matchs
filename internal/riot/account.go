package riot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) FetchPUUID(summoner string) (string, error) {
	parts := strings.Split(summoner, "#")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid summoner format, use name#tag")
	}

	game := url.PathEscape(parts[0])
	tag := url.PathEscape(parts[1])

	url := fmt.Sprintf(
		"https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s",
		game,
		tag,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", c.APIKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		PUUID string `json:"puuid"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.PUUID, err
}
