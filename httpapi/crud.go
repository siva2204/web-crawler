package httpapi

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/siva2204/web-crawler/config"
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

		urls, err := redis_crawler.Client.GetUnEncoded(search)

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
				Data:   urls,
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

		words := rootNode.AutoCompletePrefix(search)

		return c.Status(200).JSON(
			response{
				Status: true,
				Data:   words,
			})
	})

	log.Fatal(app.Listen(":" + port))

}
