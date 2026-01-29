package pipeline

import (
	"time"

	"github.com/pithandev/get-league-matchs/internal/domain"
)

func Aggregate(in <-chan domain.PlayerMatchStats) <-chan domain.AggregatedStats {
	out := make(chan domain.AggregatedStats)

	go func() {
		defer close(out)

		var (
			totalKills   int
			totalDeaths  int
			totalAssists int
			totalTime    time.Duration
			matchCount   int
			champCount   = make(map[string]int)
		)

		for stat := range in {
			matchCount++
			totalKills += stat.Kills
			totalDeaths += stat.Deaths
			totalAssists += stat.Assists
			totalTime += stat.GameDuration
			champCount[stat.Champion]++
		}

		if matchCount <= 0 {
			return
		}

		mostPlayed := ""
		max := 0
		for champ, count := range champCount {
			if count > max {
				max = count
				mostPlayed = champ
			}
		}

		avgKills := float64(totalKills) / float64(matchCount)
		avgDeaths := float64(totalDeaths) / float64(matchCount)
		avgAssists := float64(totalAssists) / float64(matchCount)

		kda := avgKills + avgAssists
		if avgDeaths > 0 {
			kda /= avgDeaths
		}

		out <- domain.AggregatedStats{
			TotalMatchs:     matchCount,
			AvgKills:        avgKills,
			AvgDeaths:       avgDeaths,
			AvgAssists:      avgAssists,
			AvgKDA:          kda,
			MostPlayedChamp: mostPlayed,
			AvgGameDuration: totalTime / time.Duration(matchCount),
		}
	}()

	return out

}
