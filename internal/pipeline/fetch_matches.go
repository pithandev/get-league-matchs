package pipeline

import (
	"sync"

	"github.com/pithandev/get-league-matchs/internal/domain"
	"github.com/pithandev/get-league-matchs/internal/riot"
)

func FetchMatches(
	client *riot.Client,
	in <-chan string,
	workers int,
) <-chan domain.MatchResponse {

	out := make(chan domain.MatchResponse)
	var wg sync.WaitGroup

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for id := range in {
				match, err := client.FetchMatchDetails(id)
				if err == nil {
					out <- match
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
