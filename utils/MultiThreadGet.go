package utils

import (
	"fmt"
	"sync"
	"time"
	"errors"
)

func downloadpcs(c string) {
	time.Sleep(1 * time.Second)
	fmt.Println("Done: ", c)
}

func MultiTGet(url string, headers map[string]string, maxThread int, maxRetry int, cnk []string) error {
	single_timeout := 2
	// var cnk []string = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"}
	// maxThread := 5
	// example

	start := time.Now()
	var wg sync.WaitGroup
	c := make(chan string, 5)
	errChan := make(chan error, maxThread) // Channel to capture errors

	wg.Add(1)
	go func() {
		for _, chunk := range cnk {
			c <- chunk
		}
		wg.Done()
		close(c)
	}()

	for i := 0; i < maxThread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reqheader := make(map[string]string, len(headers))
			for k, v := range headers {
				reqheader[k] = v
			}
			for chunk := range c {
				retries := 0
				for {
					reqheader["Range"] = chunk
					resp := NetRequest(url, "GET", reqheader, nil, single_timeout)
					err := resp.ReqErr
					retries++
					if retries < maxRetry + 1 {
						if errors.Is(err, ErrorRequestTimeout) {
							fmt.Printf("Retry %v : Request Timeout, Retry after 1s\n", retries)
							time.Sleep(751 * time.Millisecond) // Just choose this number randomly.
							continue
						}else if err != nil {
							errChan <- fmt.Errorf("download error: %w", err)
							return
						}
					} else {
						errChan <- fmt.Errorf("Too many retries: %v", maxRetry)
					}
					break
				}
			} // For testing, it will return nothing
		}()
	}

	wg.Wait()
	close(errChan)
	fmt.Println("MultiThreaD elapsed time:", time.Since(start))
	for err := range errChan {
		return err
	}
	return nil
}
