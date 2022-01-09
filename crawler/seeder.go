package crawler

import (
	"fmt"
	"time"
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

	fmt.Println("\n\n\nSlept for 5 seconds\n\n\n")

	s.Foo <- 1

	time.Sleep(time.Second * 5)

	fmt.Println("Starting again")

	s.Foo <- 0
	time.Sleep(time.Second * 5)

}
