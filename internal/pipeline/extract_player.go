package pipeline

import (
	"fmt"

	"github.com/pithandev/get-league-matchs/internal/domain"
)

func ExtractPlayer(in <-chan domain.MatchResponse, targetPUUID string) <-chan domain.PlayerMatchStats {
	out := make(chan domain.PlayerMatchStats)

	go func() {
		defer close(out)

		for match := range in {
			for _, p := range match.Info.Participants {
				if p.PUUID != targetPUUID {
					continue
				}

				stats := domain.PlayerMatchStats{
					MatchID:      match.Metadata.MatchId,
					Champion:     p.ChampionName,
					Kills:        p.Kills,
					Deaths:       p.Deaths,
					Assists:      p.Assists,
					GameDuration: match.Info.GameDuration,
				}

				fmt.Println(stats.GameDuration)
				out <- stats
				break
			}
		}
	}()

	return out
}
