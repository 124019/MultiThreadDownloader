package utils

import (
	"fmt"
	"time"
	"math/rand/v2"
	"sync"
)


func downloadpcs(c string) {
	ran := time.Duration(rand.IntN(2))
	time.Sleep(ran * time.Second)
	fmt.Println("Done: ", c)
}



func StartThread(wg *sync.WaitGroup, totalThread *chan int) {
	wg.Add(1)
	*totalThread <- 1
}
func StopThread(wg *sync.WaitGroup, totalThread *chan int) {
	wg.Done()
	close(*totalThread)
}

func main() {
	var cnk []string = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"}
	maxThread := 10
	// example
	start := time.Now()
	var wg sync.WaitGroup
	c := make(chan string, 5)

	wg.Add(1)
	go func() {
		defer close(c)
		for _, chunk := range cnk {
			c <- chunk
		}
		wg.Done()
	}()


	for i := 0; i < maxThread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range c {
				downloadpcs(chunk)
			}
		}()
	}

	wg.Wait()
	fmt.Println("Elapsed time:", time.Since(start))
}
