package crawler

import (
	"fmt"
	"time"

	"github.com/siva2204/web-crawler/config"
	"github.com/siva2204/web-crawler/db"
	redis_crawler "github.com/siva2204/web-crawler/redis"
)

type Seeder struct {
	Foo chan int
}

var SeederInstance *Seeder

func InitSeeder(crawler *Crawler) {
	Foo := make(chan int)
	crawler.SeederListener = Foo
	SeederInstance = &Seeder{Foo: Foo}
}

func (s *Seeder) Run() {
	// sleep for 10 secs
	for {

		crawlDuration := config.Config.CrawlerDuration

		time.Sleep(time.Second * time.Duration(crawlDuration))

		fmt.Printf("crawled for %d seconds, now persisting it\n", crawlDuration)

		s.Foo <- 1 // pausing crawling

		keys, err := redis_crawler.Client.GetAll()
		if err != nil {
			fmt.Errorf("unable to get keys : ", err)
		}

		values, err := redis_crawler.Client.GetMany(keys)
		if err != nil {
			fmt.Errorf("unable to get values : ", err)
		}

		db.PersistIndex(keys, values)

		// flushing redis after persisting in mysql db
		redis_crawler.Client.RDB.FlushAll(redis_crawler.Client.RDB.Context())

		fmt.Println("Successfully persisted all the data, starting to crawl in 3s...")

		time.Sleep(time.Second * 3)

		s.Foo <- 0
	}
}
