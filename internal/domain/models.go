package domain

import "time"

type AccountResponse struct {
	PUUID string `json:"puuid"`
}

type Metadata struct {
	MatchId string `json:"matchId"`
}

type InfoParticipants struct {
	ChampionName string `json:"championName"`
	Kills        int    `json:"kills"`
	Deaths       int    `json:"deaths"`
	Assists      int    `json:"assists"`
	PUUID        string `json:"puuid"`
}

type Info struct {
	GameDuration time.Duration
	Participants []InfoParticipants
}

type MatchResponse struct {
	Metadata Metadata `json:"metadata"`
	Info     Info     `json:"info"`
}

type PlayerMatchStats struct {
	MatchID      string
	Champion     string
	Kills        int
	Deaths       int
	Assists      int
	GameDuration time.Duration
}

type AggregatedStats struct {
	TotalMatchs     int
	AvgKills        float64
	AvgDeaths       float64
	AvgAssists      float64
	AvgKDA          float64
	MostPlayedChamp string
	AvgGameDuration time.Duration
}
