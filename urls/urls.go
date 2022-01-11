package urls

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jRepository struct {
	Driver neo4j.Driver
}

type URL struct {
	URL  string
	RANK float64
}

func (r *Neo4jRepository) Init() error {
	return nil
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
			return u.createUrl(tx, &URL{
				URL:  url,
				RANK: 0,
			})
		}); err != nil {
		return err
	}
	return nil
}

func (u *Neo4jRepository) createUrl(tx neo4j.Transaction, url *URL) (interface{}, error) {
	result, err := tx.Run("CREATE (u:URL {url: $url.URL, rank: $url.RANK}) RETURN u", map[string]interface{}{
		"url": url,
	})
	if err != nil {
		return nil, err
	}
	record, err := result.Single()
	if err != nil {
		return nil, err
	}
	fmt.Println(record, "sds")
	return result.Next(), nil
}

func (u *Neo4jRepository) AddPageRank(url *URL) (err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close()
	}()
	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.addPageRank(tx, url)
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
	return nil, err
}

func (u *Neo4jRepository) GetPageRank(url string) (rank float64, err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close()
	}()
	if _, err := session.
		ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.getPageRank(tx, url)
		}); err != nil {
		return 0, err
	}
	return rank, nil
}

func (u *Neo4jRepository) getPageRank(tx neo4j.Transaction, url string) (interface{}, error) {
	result, err := tx.Run("MATCH (n:URL {url: $url}) RETURN n.rank", map[string]interface{}{
		"url": url,
	})
	if err != nil {
		return nil, err
	}
	// fmt.Println(result.Single())
	return result.Next(), nil
}
