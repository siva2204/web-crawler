package neo4j

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jRepository struct {
	Driver neo4j.Driver
}

type URL struct {
	URL         string  `json:"url"`
	RANK        float64 `json:"rank"`
	DESCRIPTION string  `json:"description"`
	TITLE       string  `json:"title"`
}

type TOKEN struct {
	TOKEN string
}

func (r *Neo4jRepository) Init() error {
	return nil
}

func (u *Neo4jRepository) CreateToken(token string) error {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		_ = session.Close()
	}()
	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.createToken(tx, &TOKEN{
				TOKEN: token,
			})
		}); err != nil {
		return err
	}
	return nil
}

func (u *Neo4jRepository) createToken(tx neo4j.Transaction, token *TOKEN) (interface{}, error) {
	result, err := tx.Run("MERGE (n:TOKEN {token: $token}) RETURN n", map[string]interface{}{
		"token": token.TOKEN,
	})

	if err != nil {
		return nil, err
	}
	return result.Next(), nil
}

func (u *Neo4jRepository) CreateUrl(url string) error {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		_ = session.Close()
	}()
	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.createUrl(tx, url)
		}); err != nil {
		return err
	}
	return nil
}

func (u *Neo4jRepository) createUrl(tx neo4j.Transaction, url string) (interface{}, error) {
	result, err := tx.Run("MERGE (n:URL {url: $url}) RETURN n", map[string]interface{}{
		"url": url,
	})
	if err != nil {
		return nil, err
	}
	return result.Next(), nil
}

func (u *Neo4jRepository) GetUrlsFromToken(token string) ([]URL, error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		_ = session.Close()
	}()
	urls, err := session.
		ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.getUrlsFromToken(tx, token)
		})
	if err != nil {
		return nil, err
	}

	return urls.([]URL), nil
}

func (u *Neo4jRepository) getUrlsFromToken(tx neo4j.Transaction, token string) ([]URL, error) {
	result, err := tx.Run("MATCH (n:TOKEN {token: $token})-[:FOUNDIN]->(m:URL) RETURN m", map[string]interface{}{
		"token": token,
	})
	if err != nil {
		return nil, err
	}

	var urls []URL
	for result.Next() {
		record := result.Record()
		url, _ := record.Values[0].(neo4j.Node).Props["url"].(string)
		rank, _ := record.Values[0].(neo4j.Node).Props["rank"].(float64)
		description, _ := record.Values[0].(neo4j.Node).Props["description"].(string)
		title, _ := record.Values[0].(neo4j.Node).Props["title"].(string)
		urls = append(urls, URL{
			URL:         url,
			RANK:        rank,
			DESCRIPTION: description,
			TITLE:       title,
		})
	}

	return urls, nil
}

func (u *Neo4jRepository) AddPageRank(url string, rank float64) error {
	fmt.Println("Got the page rank for ", url, " val : ", rank)
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		_ = session.Close()
	}()
	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.addPageRank(tx, &URL{URL: url, RANK: rank})
		}); err != nil {
		return err
	}
	return nil
}

func (u *Neo4jRepository) addPageRank(tx neo4j.Transaction, url *URL) (interface{}, error) {
	_, err := tx.Run("MATCH (n:URL {url: $url}) SET n.rank = $rank", map[string]interface{}{
		"url":  url.URL,
		"rank": url.RANK,
	})
	if err != nil {
		fmt.Println(err)
	}
	return nil, err
}

func (u *Neo4jRepository) GetPageRank(url string) (float64, error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		_ = session.Close()
	}()
	rank, err := session.
		ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.getPageRank(tx, url)
		})
	if err != nil {
		return float64(0), err
	}

	return rank.(float64), nil
}

func (u *Neo4jRepository) getPageRank(tx neo4j.Transaction, url string) (float64, error) {
	result, err := tx.Run("MATCH (n:URL {url: $url}) RETURN n.rank", map[string]interface{}{
		"url": url,
	})
	if err != nil {
		return float64(0), err
	}

	record, _ := result.Single()
	rank, _ := record.Values[0].(float64)

	return rank, nil
}

func (u *Neo4jRepository) ConnectTwoUrls(url1 string, url2 string) error {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		_ = session.Close()
	}()
	_, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.connectTwoUrls(tx, url1, url2)
		})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (u *Neo4jRepository) connectTwoUrls(tx neo4j.Transaction, url1 string, url2 string) (interface{}, error) {
	result, err := tx.Run("MATCH (n:URL {url: $url1}), (m:URL {url: $url2}) MERGE (n)-[:LINK]->(m)", map[string]interface{}{
		"url1": url1,
		"url2": url2,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Neo4jRepository) GetUrls(url string) (urls []*URL, err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close()
	}()
	if _, err := session.
		ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.getUrls(tx, url)
		}); err != nil {
		return nil, err
	}
	return urls, nil
}

func (u *Neo4jRepository) getUrls(tx neo4j.Transaction, url string) (interface{}, error) {
	result, err := tx.Run("MATCH (n:URL {url: $url}) RETURN n.url", map[string]interface{}{
		"url": url,
	})
	if err != nil {
		return nil, err
	}
	for result.Next() {
		record := result.Record()
		fmt.Println(record.Values[0])
	}

	return result, nil
}

func (u *Neo4jRepository) GetConnectedUrls(url string) ([]*URL, error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		_ = session.Close()
	}()
	urls, _ := session.
		ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.getConnectedUrls(tx, url)
		})
	return urls.([]*URL), nil
}

func (u *Neo4jRepository) getConnectedUrls(tx neo4j.Transaction, url string) ([]*URL, error) {
	result, err := tx.Run("MATCH (n:URL {url: $url})-[:LINK]->(m:URL) RETURN m.url", map[string]interface{}{
		"url": url,
	})
	if err != nil {
		return nil, err
	}
	urls := make([]*URL, 0)
	for result.Next() {
		record := result.Record()
		urls = append(urls, &URL{
			URL: record.Values[0].(string),
		})
	}

	return urls, nil
}

func (u *Neo4jRepository) ConnectTokenAndUrl(token string, url string) error {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		_ = session.Close()
	}()
	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.connectTokenAndUrl(tx, token, url)
		}); err != nil {
		return err
	}
	return nil
}

func (u *Neo4jRepository) connectTokenAndUrl(tx neo4j.Transaction, token string, url string) (interface{}, error) {
	result, err := tx.Run("MATCH (n:TOKEN {token: $token}), (m:URL {url: $url}) MERGE (n)-[:FOUNDIN]->(m)", map[string]interface{}{
		"token": token,
		"url":   url,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Neo4jRepository) AddDescriptionAndTitle(description string, title string, url string) error {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		_ = session.Close()
	}()
	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.addDescriptionAndTitle(tx, description, title, url)
		}); err != nil {
		return err
	}
	return nil
}

func (u *Neo4jRepository) addDescriptionAndTitle(tx neo4j.Transaction, description string, title string, url string) (interface{}, error) {
	result, err := tx.Run("MATCH (n:URL {url: $url}) SET n.description = $description, n.title = $title", map[string]interface{}{
		"url":         url,
		"title":       title,
		"description": description,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Neo4jRepository) GetDescriptionAndTitle(url string) (description string, title string, err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		_ = session.Close()
	}()
	if _, err := session.
		ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.getDescriptionAndTitle(tx, url)
		}); err != nil {
		return "", "", err
	}
	return description, title, nil
}

func (u *Neo4jRepository) getDescriptionAndTitle(tx neo4j.Transaction, url string) (interface{}, error) {
	result, err := tx.Run("MATCH (n:URL {url: $url}) RETURN n.description, n.title", map[string]interface{}{
		"url": url,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
