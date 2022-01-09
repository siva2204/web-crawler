package crawler

import (
	"fmt"
	"sync"

	"github.com/siva2204/web-crawler/queue"
)

// crawler bot type
type Crawler struct {
	Wg      sync.WaitGroup
	Threads int
	Queue   *queue.Queue
}

// run method starts crawling
func (c *Crawler) Run() {
	// check if the url is already crawled

	if c.Queue.Len() == 0 {
		fmt.Println("queue is empty add some seed url to crawl")
		return
	}

	ch := make(chan string, 10)

	for i := 0; i < c.Threads; i++ {
		c.Wg.Add(1)

		go func(i int) {
			for {
				fmt.Println("receiving")
				url, ok := <-ch
				fmt.Println("enqueued", url)
				fmt.Printf("crawling the %s url, now..", url)
				if !ok {
					c.Wg.Done()
					return
				}

				// crawl with the url
				urls, err := uRLScrape(url)

				if err != nil {
					fmt.Printf("Error crawling url %+v", err)
					c.Wg.Done()
					return
				}

				for _, url := range urls {
					c.Queue.Enqueue(url)
				}
				// enqueue the all the related url
			}
		}(i)
	}

	// traversing the queue
	// BFS
	for {
		if c.Queue.Len() != 0 {
			fmt.Println("dequed", c.Queue.FrontQueue())

			ch <- c.Queue.Dequeue()
		}

		// TODO
		// implementing something to stop the crawling
		// may be with select and one more stop channel
	}

	close(ch)
	c.Wg.Wait()
}
