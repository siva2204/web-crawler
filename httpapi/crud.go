package httpapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/siva2204/web-crawler/config"
	neo4j_ "github.com/siva2204/web-crawler/neo4j"
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

func HttpServer(rootNode *trie.Node, graph *pagerank.PageRank, urlsRepository *neo4j_.Neo4jRepository) {
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

		// var key db.Key

		// // checking if key is there in db
		// if err := db.DB.Where("`key` = ?", search).First(&key).Error; err != nil {
		// 	fmt.Errorf("Error fetching key from db %+v", err)
		// 	return c.Status(204).JSON(
		// 		response{
		// 			Status: false,
		// 			Data:   err.Error(),
		// 		})
		// }

		// // fetching all the urls related to the key from db
		// var urls []urls

		// query := "SELECT url FROM IndexRelation LEFT JOIN Url ON IndexRelation.urlId = Url.id WHERE keyId = ? LIMIT 15;"

		// if err := db.DB.Raw(query, key.Id).Scan(&urls).Error; err != nil {
		// 	fmt.Errorf("Error fetching urls from db %+v", err)
		// 	return c.Status(500).JSON(
		// 		response{
		// 			Status: false,
		// 			Data:   err.Error(),
		// 		})
		// }

		// var data []urldata

		urls, _ := urlsRepository.GetUrlsFromToken(search)

		// for _, k := range urls {
		// 	fmt.Println(k)
		// }

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
			// db.AddPageRank(url, rank)
			urlsRepository.AddPageRank(url, rank)
			ranks[url] = rank
		})
		return c.Status(200).JSON(
			response{
				Status: true,
				Data:   ranks,
			})
	})

	app.Get("/data", func(c *fiber.Ctx) error {
		type Site struct {
			Title      string   `json:"title"`
			Body       string   `json:"body"`
			Links      []int    `json:"links"`
			ParentSite string   `json:"parent_site"`
			SiteURL    string   `json:"site_url"`
			LinkURL    []string `json:"link_urls"`
			Slug       string   `json:"slug"`
			Rank       float64
			Inbound    []int
		}
		jsonFile, err := os.Open("./site/data/site_data.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var sites []Site
		json.Unmarshal(byteValue, &sites)
		return c.Status(200).JSON(
			response{
				Status: true,
				Data:   sites,
			})
	})

	// staring the http server
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))

}
