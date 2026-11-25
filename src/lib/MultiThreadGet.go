// TODO: Add a 'if success' to CHUNK_Range to check that if all chunk can be downloaded successfully

package lib

import (
	"MultiThreadDownloader/src/core"
	"errors"
	"fmt"
	"sync"
	"time"
	"math/rand"
)


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
	var initWG sync.WaitGroup // For errors and load chunks
	var writerWG sync.WaitGroup
	var errCollectionWG sync.WaitGroup
	c := make(chan core.CHUNK_Range, config.MaxDownloadThread * 3) // Channel to post the chunk info to the download goroutine
	wr := make(chan Writer, config.MaxWritingThread * 3) // Channel to writting data into disk
	writtingErrChan := make(chan error, config.MaxWritingThread * 3) // Channel to collect if all data can be saved successfully
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
		defer writerWG.Done()
		
		var WrittingErrors []error // Catch Writting Error
		initWG.Add(1)
		go func() { // Catch Writting Error
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
				fmt.Printf("Done: MultiThreadGet\n")
				close(errChan)
			} else {
				success_multiget = false // Tempory
				for _, eachErr := range WrittingErrors {
					errChan <- eachErr
				}
				close(errChan)
			}
		}()

		for i := 0; i < config.MaxWritingThread; i++ {	// writing data into local disk
			writerWG.Add(1)
			writer_goroutine_number := i
			go func() {
				defer writerWG.Done()
				for {
					writer, open := <- wr
					if !open {
						break
					}
					fmt.Printf("WrGoroutine %v : Writing-> chunk: %v\n", writer_goroutine_number, writer.chunk.No)
					if err := config.WriteTo.WritePlaceholder(writer.chunk, writer.data); err != nil {
						writtingErrChan <- err
						fmt.Printf("! WrGoroutine %v : Failed in writing the chunk in the Storange. (chunk number: %v)\n", writer_goroutine_number, writer.chunk.No)
					} else {
						fmt.Printf("WrGoroutine %v Writing success, chunk: %v\n", writer_goroutine_number, writer.chunk.No)
					}
					continue
				}
			}()
		}
	}()


	for i := 0; i < config.MaxDownloadThread; i++ { // Start Downloading
		downloadWG.Add(1)
		go func() {
			goroutine_number := i
			defer downloadWG.Done()
			reqheader := make(map[string]string, len(config.Headers))
			for k, v := range config.Headers {
				reqheader[k] = v
			}
			fmt.Printf("Start Downloading: %v\n" , time.Since(start))
			for chunk := range c {
				start_chunk := time.Now()
				retries := 0
				for {
					var incorrect_bite_number bool
					reqheader["Range"] = chunk.Header
					resp, err := core.NetRequest(url, "GET", reqheader, nil, config.SingleChunkTimeout)
					if err != nil {
						fmt.Printf("Goroutine %v: downloading error: %v\n", goroutine_number, err)
					} else {
						if len(resp.RespBody) == chunk.Ed - chunk.St + 1 {
							fmt.Printf("Goroutine %v: Correct byte number: chunk %v(After Started: %v, This chunk: %v)\n", goroutine_number, chunk.No, time.Since(start), time.Since(start_chunk))
							// config.WriteTo.WritePlaceholder(chunk, resp.RespBody)
							wr <- Writer{chunk: chunk, data: resp.RespBody}
							// TODO : return that the chunk is finished
							break // Break this loop to get another chunk
						} else {
							fmt.Printf("! Goroutine %v: Not Correct byte number: chunk %v, Expected %v, Got %v\n! ! Detail: resp body: \n%v resp body end (After Started: %v)\n", goroutine_number, chunk.No, chunk.Ed - chunk.St + 1, len(resp.RespBody), resp.RespBody, time.Since(start))
							incorrect_bite_number = true
						}
					}	
					retries++
					if retries < config.MaxRetry + 1 {
						if errors.Is(err, core.ErrorRequestTimeout) || incorrect_bite_number {
							fmt.Printf("! Goroutine %v: (Chunk: %v) Retry %v : Request Timeout(chunk %v), Retry after 1s\n", goroutine_number, chunk.No, retries, chunk.Header)
							time.Sleep(config.RetryTimeInterval)
							continue // To Retry
						}else if err != nil {
							fmt.Printf("! Goroutine %v: (Chunk: %v) (After Started: %v) download error: %v\n! ! Details: Status: %v,Return %v bytes, Header: %v, Elapsed: %v\n", goroutine_number, chunk.No, time.Since(start), err, resp.RespStatusCode, len(resp.RespBody), resp.RespHeader, resp.ReqElapsed)
							errChan <- fmt.Errorf("! Goroutine %v: (Chunk: %v) (After Started: %v) download error: %w", goroutine_number, chunk.No, time.Since(start), err)
							continue // To Retry
						}
					} else {
						errChan <- fmt.Errorf("! Goroutine %v: (Chunk: %v) Too many retries: %v\n", goroutine_number, chunk.No, config.MaxRetry)
					}
					break
				}
			}
		}()
		time.Sleep(time.Duration(rand.Intn(500) + 500) * time.Millisecond)
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
