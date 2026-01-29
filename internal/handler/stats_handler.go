package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pithandev/get-league-matchs/internal/pipeline"
	"github.com/pithandev/get-league-matchs/internal/riot"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	summoner := r.URL.Query().Get("summoner")
	if summoner == "" {
		http.Error(w, "missing summoner", 400)
		return
	}

	puuid, err := riot.FetchPUUID(summoner)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	matchIDs, _ := riot.FetchMatchIDs(puuid)

	matchIDChan := pipeline.ProduceMatchIDs(matchIDs)
	matchesChan := pipeline.FetchMatches(matchIDChan, 5)
	playerChan := pipeline.ExtractPlayer(matchesChan, puuid)
	aggChan := pipeline.Aggregate(playerChan)

	json.NewEncoder(w).Encode(<-aggChan)
}
