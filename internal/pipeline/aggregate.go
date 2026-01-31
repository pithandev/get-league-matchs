package pipeline

import "github.com/pithandev/get-league-matchs/internal/domain"

func Aggregate(
	in <-chan domain.PlayerMatchStats,
) <-chan domain.AggregatedStats {

	out := make(chan domain.AggregatedStats)

	go func() {
		defer close(out)

		var agg domain.AggregatedStats

		for stat := range in {
			agg.TotalMatches++
			agg.TotalKills += stat.Kills
			agg.TotalDeaths += stat.Deaths
			agg.TotalAssists += stat.Assists
		}

		out <- agg
	}()

	return out
}
