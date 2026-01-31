package pipeline

func ProduceMatchIDs(ids []string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for _, id := range ids {
			out <- id
		}
	}()

	return out
}
