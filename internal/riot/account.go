package riot

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func FetchPUUID(summoner string) (string, error) {
	apiKey := os.Getenv("RIOT_API_KEY")

	parts := strings.Split(summoner, "#")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid summoner format, use name#tag")
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
		return "", err
	}

	fmt.Println("SUMMONER RECEBIDO:", summoner)
	defer resp.Body.Close()

	summ, _ := io.ReadAll(resp.Body)
	return string(summ), nil
}
