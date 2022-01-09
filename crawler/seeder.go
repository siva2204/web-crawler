package crawler

import (
	"fmt"
	"time"

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

	time.Sleep(time.Second * 5)

	fmt.Println("slept for 5 seconds")

	s.Foo <- 1

	keys, err := redis_crawler.Client.GetAll()
	if err != nil {
		fmt.Errorf("unable to get keys : ", err)
	}

	values, err := redis_crawler.Client.GetMany(keys)
	if err != nil {
		fmt.Errorf("unable to get values : ", err)
	}

	fmt.Println(values)

	// data, _ := redis_crawler.Client.Get("maintain")

	// fmt.Println(data)

	time.Sleep(time.Second * 5)

	fmt.Println("Starting again")

	s.Foo <- 0
	time.Sleep(time.Second * 5)

}
