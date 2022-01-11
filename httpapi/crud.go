package httpapi

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/siva2204/web-crawler/config"
	"github.com/siva2204/web-crawler/crawler"
	"github.com/siva2204/web-crawler/pagerank"
	redis_crawler "github.com/siva2204/web-crawler/redis"
	"github.com/siva2204/web-crawler/trie"
)

type response struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

func HttpServer(rootNode *trie.Node, graph *pagerank.PageRank) {
	app := fiber.New()
	port := config.Config.Port

	app.Static("/", "./frontend")

	app.Get("/search/:param", func(c *fiber.Ctx) error {
		search := c.Params("param")

		fmt.Println(search)

		if search == "" {
			return c.Status(400).JSON(response{
				Status: false,
				Data:   "empty param",
			})
		}

		urls, err := redis_crawler.Client.GetSetValues(search)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(204).JSON(
				response{
					Status: false,
					Data:   err.Error(),
				})
		}

		type urldata struct {
			Url         string `json:"url"`
			Title       string `json:"title"`
			Description string `json:"description"`
		}

		var data []urldata

		fmt.Println(urls)

		for _, k := range urls {

			newUrldata := urldata{
				Url:         k,
				Title:       "",
				Description: "",
			}

			title, descp, err := crawler.MetaScrape(k)

			if err != nil {
				log.Println(fmt.Errorf("error scraping from %s url", k))
			}

			newUrldata.Title = title
			newUrldata.Description = descp

			data = append(data, newUrldata)

		}

		fmt.Println(data)

		return c.Status(200).JSON(
			response{
				Status: true,
				Data:   data,
			})
	})

	app.Post("/words/:param", func(c *fiber.Ctx) error {
		search := c.Params("param")

		if search == "" {
			return c.Status(400).JSON(response{
				Status: false,
				Data:   "empty param",
			})
		}

		words, err := redis_crawler.Client.AutoComplete(search)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(204).JSON(
				response{
					Status: false,
					Data:   err.Error(),
				})
		}

		return c.Status(200).JSON(
			response{
				Status: true,
				Data:   words,
			})
	})

	// pagerank
	app.Get("/runPageRank", func(c *fiber.Ctx) error {
		probability_of_following_a_link := 0.85 // The bigger the number, less probability we have to teleport to some random link
		tolerance := 0.0001                     // the smaller the number, the more exact the result will be but more CPU cycles will be needed
		ranks := map[string]float64{}
		graph.Rank(probability_of_following_a_link, tolerance, func(url string, rank float64) {
			ranks[url] = rank
		})
		return c.Status(200).JSON(
			response{
				Status: true,
				Data:   ranks,
			})
	})

	// staring the http server
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))

}
