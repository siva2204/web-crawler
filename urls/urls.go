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
	result, err := tx.Run("CREATE (n:URL {url: $url, rank: $rank}) RETURN n", map[string]interface{}{
		"url":  url.URL,
		"rank": url.RANK,
	})
	if err != nil {
		return nil, err
	}
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
	if err != nil {
		fmt.Println(err)
	}
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
	for result.Next() {
		record := result.Record()
		fmt.Println(record.Values[0])
	}

	return result, nil
}

func (u *Neo4jRepository) ConnectTwoUrls(url1 string, url2 string) error {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		_ = session.Close()
	}()
	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.connectTwoUrls(tx, url1, url2)
		}); err != nil {
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

// func (u *Neo4jRepository) ConnectTokenToUrl(url string, token string) error {
// 	session := u.Driver.NewSession(neo4j.SessionConfig{
// 		AccessMode: neo4j.AccessModeWrite,
// 	})
// 	defer func() {
// 		_ = session.Close()
// 	}()
// 	if _, err := session.
// 		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
// 			return u.connectTwoUrls(tx, url, token)
// 		}); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (u *Neo4jRepository) connectTokenToUrl(tx neo4j.Transaction, url string, token string) (interface{}, error) {
// 	result, err := tx.Run("MATCH (n:URL {url: $url1}), (m:URL {url: $url2}) MERGE (n)-[:LINK]->(m)", map[string]interface{}{
// 		"url1": url1,
// 		"url2": url2,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }
