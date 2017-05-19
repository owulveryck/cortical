package cortical

import "sync"

func fanOut(ch <-chan []byte, size, lag int) []chan []byte {
	cs := make([]chan []byte, size)
	for i := range cs {
		// The size of the channels buffer controls how far behind the recievers
		// of the fanOut channels can lag the other channels.
		cs[i] = make(chan []byte, lag)
	}
	go func() {
		for msg := range ch {
			for _, c := range cs {
				c <- msg
			}
		}
		for _, c := range cs {
			// close all our fanOut channels when the input channel is exhausted.
			close(c)
		}
	}()
	return cs
}

func merge(done <-chan struct{}, cs ...<-chan []byte) <-chan []byte {
	var wg sync.WaitGroup
	out := make(chan []byte)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c or done is closed, then calls
	// wg.Done.
	output := func(c <-chan []byte) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
