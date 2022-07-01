package memory

import (
	"github.com/devLucian93/thesis-go/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	log "github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	NEO4J_VERTEX_CHUNK_SIZE int = 200
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
 * From here, the API of the Neo4j memory client starts
 */
func (neo *neo4jClient) GetAllVertexIds() []int64 {
	log.Println("Getting all vertex ids...")
	var vertexIds []int64
	_, err := neo.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (v: Vertex) RETURN COUNT(v)",
			map[string]interface{}{})
		if err != nil {
			return nil, err
		}
		var totalVertexCount int64
		if result.Next() {
			totalVertexCount = result.Record().Values[0].(int64)
			vertexIds = make([]int64, totalVertexCount)
		} else {
			return nil, result.Err()
		}
		log.Println("Will load", totalVertexCount, "ids...")

		var lastId int64 = -1
		var chunkSize int64 = 50_000
		var vertexCounter int64 = 0
		for vertexCounter < totalVertexCount {
			result, err = transaction.Run(
				"MATCH (v: Vertex) WHERE id(v) > $last_id RETURN id(v), v.i LIMIT $limit",
				map[string]interface{}{"last_id": lastId, "limit": chunkSize})

			if err != nil {
				return nil, err
			}
			if result.Err() != nil {
				return nil, result.Err()
			}

			var i int64
			for i = 0; result.Next(); i++ {
				lastId = result.Record().Values[0].(int64)
				vertexIds[vertexCounter] = result.Record().Values[1].(int64)
				vertexCounter++
			}
			log.Println("Read all vertex ids:", vertexCounter, "/", totalVertexCount)
		}

		return vertexIds, result.Err()
	})
	log.Println("Loaded all vertex ids...")

	if err != nil {
		panic(err)
	}

	return vertexIds
}

func (neo *neo4jClient) GetVertices(vertexIds []int64) <-chan []domain.Vertex {
	channelSize := len(vertexIds)/NEO4J_VERTEX_CHUNK_SIZE + 1
	vertexChannel := make(chan []domain.Vertex, channelSize)
	if len(vertexIds) == 0 {
		return vertexChannel
	}

	go func() {
		defer close(vertexChannel)
		for i := 0; len(vertexIds) > 0; i++ {
			nextChunkSize := NEO4J_VERTEX_CHUNK_SIZE
			if len(vertexIds) < nextChunkSize {
				nextChunkSize = len(vertexIds)
			}
			neo.getVertexChunk(nextChunkSize, vertexIds, vertexChannel, 0)

			newStart := NEO4J_VERTEX_CHUNK_SIZE
			if len(vertexIds) < newStart {
				newStart = len(vertexIds)
			}
			vertexIds = vertexIds[newStart:]
		}
		vertexChannel <- nil
	}()

	return vertexChannel
}

func (neo *neo4jClient) getVertexChunk(nextChunkSize int, vertexIds []int64, vertexChannel chan []domain.Vertex, retries int) {
	vertexChunk := make([]domain.Vertex, nextChunkSize)
	_, err := neo.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"UNWIND $ids AS id "+
				"MATCH (v: Vertex {i: id.v}) RETURN v.b, v.v",
			map[string]interface{}{"ids": mapVertexIds(vertexIds[:nextChunkSize])})
		if err != nil {
			return nil, err
		}
		for j := 0; result.Next(); j++ {
			record := result.Record().Values
			vertex := &domain.Vertex{}
			if err := vertex.UnmarshalBinary(record[0].([]byte)); err != nil {
				panic(err)
			}
			vertex.Value = record[1]
			vertexChunk[j] = *vertex
		}
		return nil, err
	})
	if err != nil {
		log.Error("Got error `", err, "` while reading")
		if retries < 10 {
			time.Sleep(time.Duration(rand.Intn(2_000)) * time.Millisecond)
			neo.getVertexChunk(nextChunkSize, vertexIds, vertexChannel, retries+1)
			return
		}
		panic(err)
	}
	vertexChannel <- vertexChunk
}

func mapVertexIds(vertexIds []int64) []map[string]interface{} {
	var vertexResult = make([]map[string]interface{}, len(vertexIds))
	for i, vertexId := range vertexIds {
		vm := make(map[string]interface{})
		vm["v"] = vertexId
		vertexResult[i] = vm
	}
	return vertexResult
}

// NOTE: 20x improvement in local environment when creating vertices with UNWIND vs for-loop
func (neo *neo4jClient) CreateVertices(vertices []domain.Vertex) {
	//_, err := neo.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
	//	_, err := transaction.Run(
	//		"UNWIND $vertices AS vertex "+
	//			"CREATE (v: Vertex {i: vertex.id, v: $initial_value, b: vertex.binary})",
	//		map[string]interface{}{"vertices": mapVertices(vertices), "initial_value": math.MaxInt64})
	//	return nil, err
	//})
	//if err != nil {
	//	panic(err)
	//}

	// NOTE: value-reset implementation
	_, err := neo.session.Run(
		"MATCH (v: Vertex) CALL { WITH v SET v.v = $initial_value} IN TRANSACTIONS OF 2000 ROWS;",
		map[string]interface{}{"initial_value": math.MaxInt64})
	if err != nil {
		panic(err)
	}
}

func mapVertices(vertices []domain.Vertex) []map[string]interface{} {
	var vertexResult = make([]map[string]interface{}, len(vertices))
	for i, vertex := range vertices {
		vm := make(map[string]interface{})
		vm["id"] = vertex.Id
		binary, err := vertex.MarshalBinary()
		if err != nil {
			panic(err)
		}
		vm["binary"] = binary
		vertexResult[i] = vm
	}
	return vertexResult
}

func (neo *neo4jClient) SaveVertices(vertices []domain.Vertex) {
	_, err := neo.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"UNWIND $vertices AS vertex "+
				"MATCH (v: Vertex {i: vertex.i}) SET v.v = vertex.v",
			map[string]interface{}{"vertices": mapVerticesToValues(vertices)})
		return nil, err
	})
	if err != nil {
		panic(err)
	}
}

func mapVerticesToValues(vertices []domain.Vertex) []map[string]interface{} {
	var vertexResult = make([]map[string]interface{}, len(vertices))
	for i, vertex := range vertices {
		vm := make(map[string]interface{})
		vm["i"] = vertex.Id
		vm["v"] = vertex.Value
		vertexResult[i] = vm
	}
	return vertexResult
}

func (neo *neo4jClient) GetGlobalParams() (*domain.GlobalParams, error) {
	gp := &domain.GlobalParams{}
	_, err := neo.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (gp: GlobalParams) RETURN gp",
			map[string]interface{}{})
		if err != nil {
			return nil, err
		}
		if result.Next() {
			nodeProps := result.Record().Values[0].(dbtype.Node).Props
			gp.RunId = nodeProps["ri"].(string)
			gp.Superstep = nodeProps["ss"].(int64)
			gp.NumberOfVertices = nodeProps["nv"].(int64)
			gp.NumberOfEdges = nodeProps["ne"].(int64)
			gp.NumberOfBuckets = nodeProps["nb"].(int64)
			gp.ChunkSize = nodeProps["cs"].(int64)
			gp.Finished = nodeProps["fn"].(bool)
			gp.DataIngestionDuration = nodeProps["id"].(int64)
			gp.ExecutionDuration = nodeProps["ed"].(int64)
			gp.Algorithm = domain.GraphAlgorithm(nodeProps["al"].(string))
			gp.GraphName = nodeProps["gn"].(string)
			err := gp.UnmarshalExtraArgs(nodeProps["ea"].([]byte))
			if err != nil {
				return nil, err
			}
			gp.MaxWorkers = nodeProps["mw"].(int64)
			return nil, nil
		}
		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	}
	return gp, nil
}

func (neo *neo4jClient) PutGlobalParams(gp *domain.GlobalParams) error {
	//log.Println("Saving global params")
	extraArgsBytes, err := gp.MarshalExtraArgs()
	if err != nil {
		return err
	}

	_, err = neo.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"MERGE (gp: GlobalParams) "+
				"ON CREATE SET gp.ri = $run_id, gp.ss = $superstep, gp.nv = $no_vertices, gp.ne = $no_edges, "+
				"	gp.nb = $no_buckets, gp.cs = $chunk_size, gp.fn = $finished, gp.id = $ingestion_duration, "+
				"	gp.ed = $exec_duration, gp.al = $algorithm, gp.gn = $graph_name, gp.ea = $extra_args, "+
				"	gp.mw = $max_workers "+
				"ON MATCH SET gp.id = $ingestion_duration, gp.ed = $exec_duration, gp.ss = $superstep, "+
				"	gp.fn = $finished",
			map[string]interface{}{"run_id": gp.RunId, "superstep": gp.Superstep, "no_vertices": gp.NumberOfVertices,
				"no_edges": gp.NumberOfEdges, "no_buckets": gp.NumberOfBuckets, "chunk_size": gp.ChunkSize,
				"finished": gp.Finished, "ingestion_duration": gp.DataIngestionDuration,
				"exec_duration": gp.ExecutionDuration, "algorithm": gp.Algorithm, "graph_name": gp.GraphName,
				"extra_args": extraArgsBytes, "max_workers": gp.MaxWorkers})
		return nil, err
	})
	return err
}

func (neo *neo4jClient) SetActiveWorkersCount(count int64) {
	_, err := neo.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"MATCH (fw: FinishedWorker) DETACH DELETE fw",
			map[string]interface{}{"count": count})
		if err != nil {
			return nil, err
		}

		_, err = transaction.Run(
			"MERGE (aw: ActiveWorkers) SET aw.count = $count",
			map[string]interface{}{"count": count})
		return nil, err
	})
	if err != nil {
		panic(err)
	}
}

// DecrementActiveWorkersCount decrements the active worker count and returns 0 if it's the last worker
func (neo *neo4jClient) DecrementActiveWorkersCount() int64 {
	var count int64 = 1
	var timestamp int64
	for true {
		_, err := neo.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
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

	_, err := neo.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
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

func (neo *neo4jClient) Clear() {
	//_, err := neo.session.Run(
	//	"MATCH (n) CALL { WITH n DETACH DELETE n} IN TRANSACTIONS OF 100000 ROWS;",
	//	map[string]interface{}{})
	//if err != nil {
	//	panic(err)
	//}

	// NOTE: value-reset implementation
	_, err := neo.session.Run(
		"MATCH (n) WHERE NOT n:Vertex CALL { WITH n DETACH DELETE n} IN TRANSACTIONS OF 100000 ROWS;",
		map[string]interface{}{})
	if err != nil {
		panic(err)
	}
}

func (neo *neo4jClient) GetFloatMcl(aggregatorKey string, superstep int64) float64 {
	var floatMcl float64
	_, err := neo.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
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

func (neo *neo4jClient) AggregateFloatMcl(aggregatorKey string, superstep int64, value float64) {
	_, err := neo.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (a: Aggregator {k: $k, s: $s}) SET a.v = $v",
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
