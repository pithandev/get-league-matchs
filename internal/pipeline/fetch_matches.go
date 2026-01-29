package pipeline

import (
	"sync"
	"time"

	"github.com/pithandev/get-league-matchs/internal/domain"
	"github.com/pithandev/get-league-matchs/internal/riot"
)

func FetchMatches(matchIDs <-chan string, workerCount int) <-chan domain.MatchResponse {
	results := make(chan domain.MatchResponse)
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go fetchMatchDetailsWorker(&wg, matchIDs, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func fetchMatchDetailsWorker(wg *sync.WaitGroup, jobs <-chan string, results chan<- domain.MatchResponse) {
	defer wg.Done()

	for matchID := range jobs {

		match, err := riot.FetchMatchDetails(matchID)
		if err != nil {
			continue
		}

		match.Info.GameDuration *= time.Second
		results <- match

	}
}
