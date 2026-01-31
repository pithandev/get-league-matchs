package domain

import "time"

type MatchResponse struct {
	Metadata struct {
		MatchID string `json:"matchId"`
	} `json:"metadata"`

	Info struct {
		GameDuration int `json:"gameDuration"`
		Participants []struct {
			PUUID        string `json:"puuid"`
			ChampionName string `json:"championName"`
			Kills        int    `json:"kills"`
			Deaths       int    `json:"deaths"`
			Assists      int    `json:"assists"`
		} `json:"participants"`
	} `json:"info"`
}

type PlayerMatchStats struct {
	Champion string
	Kills    int
	Deaths   int
	Assists  int
	Duration time.Duration
}

type AggregatedStats struct {
	TotalMatches int
	TotalKills   int
	TotalDeaths  int
	TotalAssists int
}
