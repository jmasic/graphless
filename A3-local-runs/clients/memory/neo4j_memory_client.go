package memory

import (
	"github.com/devLucian93/thesis-go/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type neo4jClient struct {
	session neo4j.Session
}

func newNeo4jClient() (Client, error) {
	neo4jClient := &neo4jClient{}
	//uri := "neo4j://127.0.0.1:7687"
	uri := "neo4j://3.144.32.80:7687"
	username := "neo4j"
	password := "n"

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))

	if err != nil {
		return nil, err
	}

	neo4jClient.session = driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	return neo4jClient, nil
}

func (memory *neo4jClient) VertexRange(startKey int64, endKey int64) []domain.Vertex {
	panic("VertexRange NOT IMPLEMENTED YET")
}

// This function should work
func (memory *neo4jClient) GetAllVertexIds() []int64 {
	var vertexIds []int64
	_, err := memory.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (v: Vertex) RETURN v.id",
			map[string]interface{}{})

		if err != nil {
			return nil, err
		}
		if result.Err() != nil {
			return nil, result.Err()
		}

		for result.Next() {
			vertexIds = append(vertexIds, result.Record().Values[0].(int64))
		}
		return vertexIds, result.Err()
	})

	if err != nil {
		panic(err)
	}

	return vertexIds
}

func (memory *neo4jClient) GetVertices(intKeys []int64) []domain.Vertex {
	vertices := make([]domain.Vertex, 0, len(intKeys))
	for _, vertexId := range intKeys {
		_, err := memory.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			result, err := transaction.Run(
				"MATCH (v: Vertex {id: $id}) RETURN v.binary",
				map[string]interface{}{"id": vertexId})

			if err != nil {
				return nil, err
			}

			for result.Next() {
				vertex := &domain.Vertex{}
				vertexBytes := result.Record().Values[0].([]byte)
				err = vertex.UnmarshalBinary(vertexBytes)
				vertices = append(vertices, *vertex)
			}
			return vertices, result.Err()
		})
		if err != nil {
			panic(err)
		}
	}
	return vertices
}

func (memory *neo4jClient) PutVertex(vertex *domain.Vertex) {
	jsonBytes, _ := vertex.MarshalBinary()
	_, err := memory.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"MERGE (v: Vertex {id: $id}) ON CREATE SET v.binary = $binary, v.active = 0 ON MATCH SET v.binary = $binary",
			map[string]interface{}{"id": vertex.Id, "binary": jsonBytes})
		return nil, err
	})

	if err != nil {
		panic(err)
	}
}

func (memory *neo4jClient) PutVertices(vertices []domain.Vertex) {
	for _, vertex := range vertices {
		memory.PutVertex(&vertex)
	}
}

func (memory *neo4jClient) DeleteVertex(key string) error {
	log.Println("----------------------")
	log.Println("The key to be deleted is: ", key)
	_, err := memory.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (n:Vertex {id: $key}) DETACH DELETE n",
			map[string]interface{}{"key": key})
		if err != nil {
			return nil, err
		}

		return nil, result.Err()
	})

	if err != nil {
		panic(err)
	}
	return nil
}

func (memory *neo4jClient) GetGlobalParams() (*domain.GlobalParams, error) {
	gp := &domain.GlobalParams{}
	gpBytes, err := memory.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a: GlobalParams) RETURN a.binary",
			map[string]interface{}{})
		if err != nil {
			return nil, err
		}
		if result.Next() {
			return result.Record().Values[0], nil
		}
		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	}
	err = gp.UnmarshalBinary(gpBytes.([]byte))
	return gp, err
}

func (memory *neo4jClient) PutGlobalParams(gp *domain.GlobalParams) error {
	//log.Println("Saving global params")
	gpBytes, err := gp.MarshalBinary()
	//log.Println("Saving global params ", &(gp))
	if err != nil {
		return err
	}

	_, err = memory.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"MERGE (gp: GlobalParams) SET gp.binary = $binary",
			map[string]interface{}{"binary": gpBytes})
		return nil, err
	})
	return err
}

func (memory *neo4jClient) AddActiveVertices(activeVertices []int64) {
	for _, v := range activeVertices {
		_, _ = memory.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			_, err := transaction.Run(
				"MATCH (v: Vertex {id: $id}) SET v.active = 1",
				map[string]interface{}{"id": v})
			return nil, err

		})
	}
}

func (memory *neo4jClient) RemoveHaltedVertices(haltedVertices []int64) {
	log.Info("Removing halted vertices: ", haltedVertices)
	for _, v := range haltedVertices {
		_, _ = memory.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			_, err := transaction.Run(
				"MATCH (v: Vertex {id: $id}) SET v.active = 0",
				map[string]interface{}{"id": v})
			return nil, err

		})
	}
}

func (memory *neo4jClient) GetActiveVertices() []int64 {
	var activeVertices []int64
	_, err := memory.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (n: Vertex {active: 1}) RETURN n.id",
			map[string]interface{}{})

		if err != nil {
			return nil, err
		}

		if result.Err() != nil {
			return nil, result.Err()
		}

		for result.Next() {
			activeVertices = append(activeVertices, result.Record().Values[0].(int64))
		}
		return nil, err
	})

	if err != nil {
		panic(err)
	}

	return activeVertices
}

func (memory *neo4jClient) GetActiveVerticesCount() int64 {
	var countActiveVertices int64
	_, _ = memory.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (n: Vertex {active: 1}) RETURN count(n)",
			map[string]interface{}{})
		if err != nil {
			return nil, err
		}
		if result.Err() != nil {
			return nil, result.Err()
		}

		if result.Next() {
			countActiveVertices = result.Record().Values[0].(int64)
			log.Info("Number of active vertices: ", countActiveVertices)
			return result.Record().Values[0], nil
		}
		return result, err
	})

	return countActiveVertices
}

func (memory *neo4jClient) SetActiveWorkersCount(count int64) {
	_, err := memory.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"MATCH (fw: FinishedWorker) DETACH DELETE fw",
			map[string]interface{}{"count": count})
		return nil, err
	})
	if err != nil {
		panic(err)
	}

	_, err = memory.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"MERGE (aw: ActiveWorkers) SET aw.count = $count",
			map[string]interface{}{"count": count})
		return nil, err
	})
	if err != nil {
		panic(err)
	}
}

// DecrementActiveWorkersCount decrements the active worker count and returns 0 if it's the last worker
func (memory *neo4jClient) DecrementActiveWorkersCount() int64 {
	var count int64 = 1
	var timestamp int64
	for true {
		_, err := memory.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			result, err := transaction.Run(
				"CREATE (fw: FinishedWorker {t: TIMESTAMP()}) RETURN fw.t",
				map[string]interface{}{})
			if err != nil {
				return nil, err
			}
			if result.Next() {
				timestamp = result.Record().Values[0].(int64)
				return nil, nil
			}
			return nil, result.Err()
		})
		if err == nil {
			break
		}
		time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
	}

	_, err := memory.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (fw: FinishedWorker) RETURN count(fw) AS c UNION ALL MATCH (aw: ActiveWorkers) RETURN aw.count AS c",
			map[string]interface{}{})
		if err != nil {
			return nil, err
		}

		counts, _ := result.Collect()
		finishedWorkersCount := counts[0].Values[0].(int64)
		activeWorkersCount := counts[1].Values[0].(int64)
		count = finishedWorkersCount - activeWorkersCount
		if finishedWorkersCount == activeWorkersCount {
			// The following should fix concurrency issues
			concurrencyCheckResult, err := transaction.Run(
				"MATCH (fw: FinishedWorker) RETURN fw.t ORDER BY fw.t DESC LIMIT 1",
				map[string]interface{}{})
			if err != nil {
				return nil, err
			}
			winningIds, _ := concurrencyCheckResult.Collect()
			if winningIds[0].Values[0] == timestamp {
				count = 0
			} else {
				count = -1
			}
		}
		return nil, result.Err()
	})
	if err != nil {
		panic(err)
	}

	return count
}

func (memory *neo4jClient) CountReceiversForSuperstep(superstep int64) int64 {
	panic("CountReceiversForSuperstep NOT IMPLEMENTED YET")
}

func (memory *neo4jClient) GetMessageRecipients(superstep int64) []int64 {
	log.Println("Getting recipients for superstep:", superstep)

	var recipients []int64
	_, err := memory.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
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

func (memory *neo4jClient) GetMessages(vertexId int64, superstep int64) []interface{} {
	var messages []interface{}
	_, err := memory.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
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

func (memory *neo4jClient) PutMessageForAllEdges(recipients []domain.Edge, message interface{}, superstep int64) {
	for i := 0; i < len(recipients); i++ {
		memory.PutMessage(recipients[i].TargetVertexId, message, superstep)
	}
}

func (memory *neo4jClient) PutMessages(recipients []domain.Edge, messages []interface{}, superstep int64) {
	for i := 0; i < len(messages); i++ {
		memory.PutMessage(recipients[i].TargetVertexId, messages[i], superstep)
	}
}

func (memory *neo4jClient) PutMessage(recipientId int64, message interface{}, superstep int64) {
	_, err := memory.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
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

	_, err = memory.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MERGE (r: Recipient {s: $s, i: $i})",
			map[string]interface{}{"s": superstep, "i": recipientId})
		if err != nil {
			return nil, err
		}
		return nil, result.Err()
	})
	if err != nil {
		panic(err)
	}
}

func (memory *neo4jClient) Clear() {
	_, err := memory.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (n) DETACH DELETE n",
			map[string]interface{}{})
		if err != nil {
			return nil, err
		}

		return nil, result.Err()
	})

	if err != nil {
		panic(err)
	}
}

func (memory *neo4jClient) CreateAggregatorMcl(aggregatorKey string) {
	log.Info("Creating aggregator MCL...")
}

func (memory *neo4jClient) ResetAggregatorsMcl() {
	panic("ResetAggregatorsMcl NOT IMPLEMENTED YET")
}

func (memory *neo4jClient) GetFloatMcl(aggregatorKey string, superstep int64) float64 {
	var floatMcl float64
	_, err := memory.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a: Aggregator {k: $k, s: $s}) RETURN sum(a.v) + 0.0",
			map[string]interface{}{"k": aggregatorKey, "s": superstep})

		if result.Next() {
			floatMcl = result.Record().Values[0].(float64)
			return nil, result.Err()
		}
		if err != nil {
			panic(err)
		}
		return result, err
	})

	if err != nil {
		panic(err)
	}
	return floatMcl
}

func (memory *neo4jClient) AggregateFloatMcl(aggregatorKey string, superstep int64, value float64) {
	_, err := memory.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (a: Aggregator {k: $k, s: $s, v: $v})",
			map[string]interface{}{"k": aggregatorKey, "s": superstep, "v": value})

		if err != nil {
			return nil, err
		}
		if result.Err() != nil {
			return nil, result.Err()
		}

		return nil, err
	})
	if err != nil {
		panic(err)
	}
}
