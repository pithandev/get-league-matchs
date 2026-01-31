package pipeline

import (
	"time"

	"github.com/pithandev/get-league-matchs/internal/domain"
)

func ExtractPlayer(
	in <-chan domain.MatchResponse,
	puuid string,
) <-chan domain.PlayerMatchStats {

	out := make(chan domain.PlayerMatchStats)

	go func() {
		defer close(out)

		for match := range in {
			for _, p := range match.Info.Participants {
				if p.PUUID == puuid {
					out <- domain.PlayerMatchStats{
						Champion: p.ChampionName,
						Kills:    p.Kills,
						Deaths:   p.Deaths,
						Assists:  p.Assists,
						Duration: time.Duration(match.Info.GameDuration) * time.Second,
					}
				}
			}
		}
	}()

	return out
}
