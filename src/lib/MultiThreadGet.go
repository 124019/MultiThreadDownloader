package lib

import (
	"fmt"
	"sync"
	"time"
	"errors"
	"MultiThreadDownloader/src/core"
)

func downloadpcs(c string) {
	time.Sleep(1 * time.Second)
	fmt.Println("Done: ", c)
}

type MultiGetConfig struct {
	WriteTo            FixedSizeFile
	RetryTimeInterval  time.Duration
	MaxRetry           int
	MaxThread 		   int
	Headers            map[string]string
	SingleChunkTimeout int
}

type Writer struct {
	chunk core.CHUNK_Range
	data  []byte
}

func MultiTGet(url string, cnk []core.CHUNK_Range, config MultiGetConfig) error { // cnk is the chunk of the file
	start := time.Now()
	var wg sync.WaitGroup
	var writerWG sync.WaitGroup
	c := make(chan core.CHUNK_Range, 5)
	wr := make(chan Writer)
	errChan := make(chan error, config.MaxThread) // Channel to capture errors

	wg.Add(1)
	go func() {
		for _, chunk := range cnk {
			c <- chunk
		}
		wg.Done()
		close(c)
	}()

	writerWG.Add(1)
	go func() {
		defer writerWG.Done()
		for {
			writer, open := <- wr
			if !open {
				break
			}
			fmt.Printf("writer header: %v\n", writer.chunk)
			if err := config.WriteTo.WritePlaceholder(writer.chunk, writer.data); err != nil {
				errChan <- err
				fmt.Println("Failed in writing the chunk in the Storange.")
				return
			}
			continue
		}
		fmt.Println("Done: MultiThreadGet")
	}()

	for i := 0; i < config.MaxThread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reqheader := make(map[string]string, len(config.Headers))
			for k, v := range config.Headers {
				reqheader[k] = v
			}
			for chunk := range c {
				retries := 0
				for {
					reqheader["Range"] = chunk.Header
					resp, err := core.NetRequest(url, "GET", reqheader, nil, config.SingleChunkTimeout)
					if len(resp.RespBody) != chunk.Ed - chunk.St + 1 {
						fmt.Printf("Not Correct byte number: chunk %v\n, Expected %v, Got %v\n", chunk.Header, chunk.Ed - chunk.St + 1, len(resp.RespBody))
						fmt.Printf("resp body: %v\n resp end", resp.RespBody)
					} else {
						fmt.Printf("Correct byte number: chunk %v\n", chunk.Header)
						// config.WriteTo.WritePlaceholder(chunk, resp.RespBody)
						wr <- Writer{chunk: chunk, data: resp.RespBody}
					}
					retries++
					if retries < config.MaxRetry + 1 {
						if errors.Is(err, core.ErrorRequestTimeout) {
							fmt.Printf("Retry %v : Request Timeout(chunk %v), Retry after 1s\n", retries, chunk.Header)
							time.Sleep(config.RetryTimeInterval)
							continue
						}else if err != nil {
							errChan <- fmt.Errorf("download error: %w", err)
							return
						}
					} else {
						errChan <- fmt.Errorf("Too many retries: %v", config.MaxRetry)
					}
					break
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(wr)
	}()

	wg.Wait()
	writerWG.Wait()
	close(errChan)
	fmt.Println("MultiThreaD elapsed time:", time.Since(start))
	for err := range errChan {
		return err
	}
	return nil
}
