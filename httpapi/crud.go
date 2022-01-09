package httpapi

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/siva2204/web-crawler/config"
	"github.com/siva2204/web-crawler/crawler"
	redis_crawler "github.com/siva2204/web-crawler/redis"
	"github.com/siva2204/web-crawler/trie"
)

type response struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

func HttpServer(rootNode *trie.Node) {
	app := fiber.New()
	port := config.Getenv("PORT")

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
			Title       string `json:"title`
			Description string `json:"description`
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
				fmt.Errorf("error scraping from %s url", k)
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

	// staring the http server
	log.Fatal(app.Listen(":" + port))

}
