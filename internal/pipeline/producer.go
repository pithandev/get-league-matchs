package pipeline

func ProduceMatchIDs(matchIDs []string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for _, id := range matchIDs {
			out <- id
		}
	}()

	return out

}
