package message

import (
	"github.com/devLucian93/thesis-go/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type neo4jClient struct {
	session neo4j.Session
}

func newNeo4jClient(config domain.DatabaseConfig) (Client, error) {
	neo4jClient := &neo4jClient{}

	targetUri := "neo4j://" + config.Ip + ":" + strconv.Itoa(config.Port)
	driver, err := neo4j.NewDriver(targetUri, neo4j.BasicAuth(config.Username, config.Password, ""))

	if err != nil {
		return nil, err
	}

	neo4jClient.session = driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	// NOTE: Indexes are created when setting up neo4j

	return neo4jClient, nil
}

/**
 * From here, the API of the Neo4j message client starts
 */
func (neo *neo4jClient) CountReceiversForSuperstep(superstep int64) int64 {
	panic("CountReceiversForSuperstep NOT IMPLEMENTED YET")
}

func (neo *neo4jClient) GetMessageRecipients(superstep int64) []int64 {
	log.Println("Getting recipients for superstep:", superstep)

	var recipients []int64
	_, err := neo.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (r: Recipient {s: $s}) RETURN r.i",
			map[string]interface{}{"s": superstep})

		if err != nil {
			return nil, err
		}
		if result.Err() != nil {
			return nil, result.Err()
		}

		for result.Next() {
			recipient := result.Record().Values[0]
			recipients = append(recipients, recipient.(int64))
		}

		return nil, err
	})

	if err != nil {
		panic(err)
	}

	return recipients
}

func (neo *neo4jClient) GetMessages(vertexId int64, superstep int64) []interface{} {
	var messages []interface{}
	_, err := neo.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (m: Message {s: $s, r: $r}) RETURN m.v",
			map[string]interface{}{"s": superstep, "r": vertexId})

		for result.Next() {
			messages = append(messages, result.Record().Values[0])
		}
		return result, err
	})

	if err != nil {
		panic(err)
	}
	return messages
}

func (neo *neo4jClient) PutMessageForAllEdges(recipients []domain.Edge, message interface{}, superstep int64) {
	recipientIds := mapMessageRecipients(recipients)
	_, err := neo.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"UNWIND $recipients AS recipient "+
				"MERGE (r: Recipient {s: $s, i: recipient.r})",
			map[string]interface{}{"s": superstep, "recipients": recipientIds})
		if err != nil {
			panic(err)
		}
		_, err = transaction.Run(
			"UNWIND $recipients AS recipient "+
				"CREATE (m: Message {s: $s, r: recipient.r, v: $message})",
			map[string]interface{}{"s": superstep, "message": message, "recipients": recipientIds})
		return nil, err
	})
	if err != nil {
		panic(err)
	}
}

func mapMessageRecipients(recipients []domain.Edge) []map[string]interface{} {
	var result = make([]map[string]interface{}, len(recipients))
	for i, recipient := range recipients {
		m := make(map[string]interface{})
		m["r"] = recipient.TargetVertexId
		result[i] = m
	}
	return result
}

func (neo *neo4jClient) PutMessages(recipients []domain.Edge, messages []interface{}, superstep int64) {
	_, err := neo.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"UNWIND $recipients AS recipient "+
				"MERGE (r: Recipient {s: $s, i: recipient.r})",
			map[string]interface{}{"s": superstep, "recipients": mapMessageRecipients(recipients)})
		_, err = transaction.Run(
			"UNWIND $messages AS message "+
				"CREATE (m: Message {s: $s, r: message.r, v: message.v})",
			map[string]interface{}{"s": superstep, "messages": mapMessages(recipients, messages)})
		return nil, err
	})

	if err != nil {
		panic(err)
	}
}

func mapMessages(recipients []domain.Edge, messages []interface{}) []map[string]interface{} {
	var result = make([]map[string]interface{}, len(recipients))
	for i, recipient := range recipients {
		m := make(map[string]interface{})
		m["r"] = recipient.TargetVertexId
		m["v"] = messages[i]
		result[i] = m
	}
	return result
}

func (neo *neo4jClient) PutMessage(recipientId int64, message interface{}, superstep int64) {
	_, err := neo.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MERGE (r: Recipient {s: $s, i: $i})",
			map[string]interface{}{"s": superstep, "i": recipientId})
		if err != nil {
			return nil, err
		}
		if result.Err() != nil {
			return nil, result.Err()
		}

		result, err = transaction.Run(
			"CREATE (m: Message {s: $s, r: $r, v: $v})",
			map[string]interface{}{"s": superstep, "r": recipientId, "v": message})
		if err != nil {
			return nil, err
		}
		return nil, result.Err()
	})
	if err != nil {
		panic(err)
	}
}

func (neo *neo4jClient) Clear() {
	//_, err := neo.session.Run(
	//	"MATCH (n) CALL { WITH n DETACH DELETE n} IN TRANSACTIONS OF 100000 ROWS;",
	//	map[string]interface{}{})
	//if err != nil {
	//	panic(err)
	//}
}
