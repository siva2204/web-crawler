package httpapi

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/siva2204/web-crawler/config"
	"github.com/siva2204/web-crawler/crawler"
	"github.com/siva2204/web-crawler/db"
	"github.com/siva2204/web-crawler/pagerank"
	redis_crawler "github.com/siva2204/web-crawler/redis"
	"github.com/siva2204/web-crawler/trie"
)

type response struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

type urldata struct {
	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type urls struct {
	Url string `gorm:"url"`
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

		var key db.Key

		// checking if key is there in db
		if err := db.DB.Where("`key` = ?", search).First(&key).Error; err != nil {
			fmt.Errorf("Error fetching key from db %+v", err)
			return c.Status(204).JSON(
				response{
					Status: false,
					Data:   err.Error(),
				})
		}

		// fetching all the urls related to the key from db
		var urls []urls

		query := "SELECT url FROM IndexRelation LEFT JOIN Url ON IndexRelation.urlId = Url.id WHERE keyId = ? LIMIT 15;"

		if err := db.DB.Raw(query, key.Id).Scan(&urls).Error; err != nil {
			fmt.Errorf("Error fetching urls from db %+v", err)
			return c.Status(500).JSON(
				response{
					Status: false,
					Data:   err.Error(),
				})
		}

		var data []urldata

		for _, k := range urls {

			newUrldata := urldata{
				Url:         k.Url,
				Title:       "",
				Description: "",
			}

			title, descp, err := crawler.MetaScrape(k.Url)

			if err != nil {
				log.Println(fmt.Errorf("error scraping from %s url", k.Url))
			}

			newUrldata.Title = title
			newUrldata.Description = descp

			data = append(data, newUrldata)
		}

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
			db.AddPageRank(url, rank)
		})
		return c.Status(200).JSON(
			response{
				Status: true,
				Data:   ranks,
			})
	})

	app.Get("/data", func(c *fiber.Ctx) error {
		file := "../static/data/site_data.json"
		return c.JSON(response{
			Status: true,
			Data:   file,
		})

	})

	// staring the http server
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))

}
