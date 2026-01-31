package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pithandev/get-league-matchs/internal/pipeline"
	"github.com/pithandev/get-league-matchs/internal/riot"
)

func StatsHandler(client *riot.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		summoner := r.URL.Query().Get("summoner")
		if summoner == "" {
			http.Error(w, "missing summoner", 400)
			return
		}

		puuid, err := client.FetchPUUID(summoner)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		ids, err := client.FetchMatchIDs(puuid)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		idChan := pipeline.ProduceMatchIDs(ids)
		matchChan := pipeline.FetchMatches(client, idChan, 5)
		playerChan := pipeline.ExtractPlayer(matchChan, puuid)
		aggChan := pipeline.Aggregate(playerChan)

		json.NewEncoder(w).Encode(<-aggChan)
	}
}
