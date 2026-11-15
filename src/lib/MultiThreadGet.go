// TODO: Add a 'if success' to CHUNK_Range to check that if all chunk can be downloaded successfully

package lib

import (
	"MultiThreadDownloader/src/core"
	"errors"
	"fmt"
	"sync"
	"time"
)

func downloadpcs(c string) {
	time.Sleep(1 * time.Second)
	fmt.Println("Done: ", c)
}

type MultiGetConfig struct {
	WriteTo            FixedSizeFile
	RetryTimeInterval  time.Duration
	MaxRetry           int
	MaxDownloadThread  int
	MaxWritingThread   int
	Headers            map[string]string
	SingleChunkTimeout int
}

type Writer struct {
	chunk core.CHUNK_Range
	data  []byte
}

func MultiTGet(url string, cnk []core.CHUNK_Range, config MultiGetConfig) (bool, []error) { // cnk is the chunk of the file
	var success_multiget bool
	var errs []error
	start := time.Now()
	var downloadWG sync.WaitGroup
	var initWG sync.WaitGroup
	var writerWG sync.WaitGroup
	var errCollectionWG sync.WaitGroup
	c := make(chan core.CHUNK_Range, config.MaxDownloadThread * 3) // Channel to post the chunk info to the download gorouting
	wr := make(chan Writer, config.MaxWritingThread * 3) // Channel to writting data into disk
	writtingErrChan := make(chan error, config.MaxWritingThread) // Channel to collect if all data can be saved successfully
	errChan := make(chan error, config.MaxDownloadThread * 3) // Channel to capture errors

	errCollectionWG.Add(1)
	go func(){
		defer errCollectionWG.Done()
		for {
			each_err, open := <- errChan
			if !open {
				break
			}
			if each_err != nil {
				errs = append(errs, each_err)
			}
		}
	}()

	initWG.Add(1)
	go func() {
		defer initWG.Done()
		for _, chunk := range cnk {
			c <- chunk
		}
		close(c)
	}()

	writerWG.Add(1) // load the writers
	func() {
		go func() {
			defer writerWG.Done()
			for i := 0; i < config.MaxWritingThread; i++ {	// writing data into local disk
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
							writtingErrChan <- err
							fmt.Println("! Failed in writing the chunk in the Storange.")
						}
						continue	
					}
				}()
			}
		}() // All writerWG end
		
		var WrittingErrors []error
		initWG.Add(1)
		go func() {
			defer initWG.Done()
			for {
				eachWrittingError, open := <- writtingErrChan
				if !open {
					break
				}
				WrittingErrors = append(WrittingErrors, eachWrittingError)
			}
			if len(WrittingErrors) == 0 {
				success_multiget = true // Tempory
				fmt.Println("Done: MultiThreadGet")
				close(errChan)
			} else {
				success_multiget = false // Tempory
				for _, eachErr := range WrittingErrors {
					errChan <- eachErr
				}
				close(errChan)
			}
		}()
	}()


	for i := 0; i < config.MaxDownloadThread; i++ {
		downloadWG.Add(1)
		go func() {
			defer downloadWG.Done()
			reqheader := make(map[string]string, len(config.Headers))
			for k, v := range config.Headers {
				reqheader[k] = v
			}
			for chunk := range c {
				retries := 0
				for {
					var incorrect_bite_number bool
					reqheader["Range"] = chunk.Header
					resp, err := core.NetRequest(url, "GET", reqheader, nil, config.SingleChunkTimeout)
					if err != nil {
						fmt.Println("downloading error: ", err)
						errChan <- err // only for test, it can't do anything
					} else {
						if len(resp.RespBody) != chunk.Ed - chunk.St + 1 {
							fmt.Printf("! Not Correct byte number: chunk %v\n, Expected %v, Got %v\n", chunk.Header, chunk.Ed - chunk.St + 1, len(resp.RespBody))
							fmt.Printf("! ! resp body: %v\n resp end", resp.RespBody)
							incorrect_bite_number = true
						} else {
							fmt.Printf("Correct byte number: chunk %v\n", chunk.Header)
							// config.WriteTo.WritePlaceholder(chunk, resp.RespBody)
							wr <- Writer{chunk: chunk, data: resp.RespBody}
						}
					}	
					retries++
					if retries < config.MaxRetry + 1 {
						if errors.Is(err, core.ErrorRequestTimeout) || incorrect_bite_number {
							fmt.Printf("Retry %v : Request Timeout(chunk %v), Retry after 1s\n", retries, chunk.Header)
							time.Sleep(config.RetryTimeInterval)
							continue
						}else if err != nil {
							fmt.Println("download error:", err)
							errChan <- fmt.Errorf("download error: %w", err)
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
		downloadWG.Wait()
		close(wr)
	}()
	go func() {
		writerWG.Wait()
		close(writtingErrChan)
	}()
	initWG.Wait()
	errCollectionWG.Wait()
	fmt.Println("MultiThreaD elapsed time:", time.Since(start))
	
	return success_multiget, errs
}
