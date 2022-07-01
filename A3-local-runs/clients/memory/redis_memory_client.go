package memory

import (
	"errors"
	"fmt"
	"github.com/devLucian93/thesis-go/domain"
	"github.com/devLucian93/thesis-go/utils"
	"github.com/go-redis/redis"
	"log"
	"time"
)

const (
	GLOBAL_PARAMS   = "globalParams"
	VERTEX          = "vertices"
	MESSAGES        = "messages"
	ACTIVE_VERTICES = "activeVertices"
	ACTIVE_WORKERS  = "activeWorkers"
	AGGREGATORS     = "aggregators"
)

type redisClient struct {
	client *redis.Ring
}

func newRedisClient() (Client, error) {
	client := redis.NewRing(&redis.RingOptions{
		// r4.large
		Addrs: map[string]string{
			//TRY AWS PRIVATE IP TO AVOID REGIONAL DATA TRANSFER CHARGES!
			//"redis-6379": "localhost:6379",
			"redis-6379": "18.191.147.238:6379",
			"redis-6380": "18.191.147.238:6380",
			"redis-6381": "18.191.147.238:6381",
			"redis-6382": "18.191.147.238:6382",
			"redis-6383": "18.191.147.238:6383",
			"redis-6384": "18.191.147.238:6384",
			"redis-6385": "18.191.147.238:6385",
			"redis-6386": "18.191.147.238:6386",
			"redis-6387": "18.191.147.238:6387",
			"redis-6388": "18.191.147.238:6388",
			"redis-6389": "18.191.147.238:6389",
			"redis-6390": "18.191.147.238:6390",
			"redis-6391": "18.191.147.238:6391",
			"redis-6392": "18.191.147.238:6392",
			"redis-6393": "18.191.147.238:6393",
			"redis-6394": "18.191.147.238:6394",
		},
		MaxRetries:      20,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
		DialTimeout:     30 * time.Second,
		ReadTimeout:     20 * time.Second,
		WriteTimeout:    20 * time.Second,
		PoolSize:        20,
		Password:        "",
		DB:              0,
	})
	redisCli := &redisClient{client}
	return redisCli, nil
}

func (memory *redisClient) VertexRange(startKey int64, endKey int64) []domain.Vertex {
	var keys []string
	for i := startKey; i < endKey; i++ {
		keys = append(keys, fmt.Sprintf("%v:%v", VERTEX, i))
	}

	vertices := getMultipleKeys(keys, memory)

	return vertices
}

func (memory *redisClient) GetAllVertexIds() []int64 {
	var vertexIds []int64
	err := memory.client.SMembers(VERTEX).ScanSlice(&vertexIds)

	if err != nil {
		panic(err)
	}

	return vertexIds
}

func (memory *redisClient) GetVertices(intKeys []int64) []domain.Vertex {
	var keys []string
	for _, key := range intKeys {
		keys = append(keys, fmt.Sprintf("%v:%v", VERTEX, key))
	}
	vertices := make([]domain.Vertex, 0, len(intKeys))

	vertices = append(vertices, getMultipleKeys(keys, memory)...)

	return vertices
}

func getMultipleKeys(keys []string, memory *redisClient) []domain.Vertex {
	verticesCmd := make([]*redis.StringCmd, len(keys), len(keys))
	vertices := make([]domain.Vertex, 0, len(keys))

	pipe := memory.client.Pipeline()

	for i := 0; i < len(keys); i++ {
		verticesCmd[i] = pipe.Get(keys[i])
	}
	_, err := pipe.Exec()

	if err != nil {
		log.Println("Error, cannot get keys", keys)
		panic(err)
	}

	for _, vertexCmd := range verticesCmd {
		vertex := domain.Vertex{}
		decompressed := utils.ZLibDecompress([]byte(vertexCmd.Val()))
		vertex.UnmarshalBinary(decompressed)
		vertices = append(vertices, vertex)
	}

	return vertices
}

func (memory *redisClient) PutVertex(vertex *domain.Vertex) {
	pipe := memory.client.Pipeline()

	pipe.SAdd(VERTEX, vertex.Id)
	jsonBytes, _ := vertex.MarshalBinary()
	pipe.Set(fmt.Sprintf("%v:%v", VERTEX, vertex.Id), utils.ZLibCompress(jsonBytes), 0)

	_, err := pipe.Exec()

	if err != nil {
		panic(err)
	}
}

func (memory *redisClient) PutVertices(vertices []domain.Vertex) {
	vertexIds := make([]interface{}, len(vertices))
	pipe := memory.client.Pipeline()
	for index, vertex := range vertices {
		vertexIds[index] = vertex.Id
		jsonBytes, _ := vertex.MarshalBinary()
		pipe.Set(fmt.Sprintf("%v:%v", VERTEX, vertex.Id), utils.ZLibCompress(jsonBytes), 0)
	}
	_, err := pipe.Exec()

	if err != nil {
		panic(err)
	}

	memory.client.SAdd(VERTEX, vertexIds...).Err()

	if err != nil {
		panic(err)
	}
}

func (memory *redisClient) DeleteVertex(key string) error {
	//to implement
	return errors.New("Unimplemented method 'DeleteVertex'")
}

func (memory *redisClient) GetGlobalParams() (*domain.GlobalParams, error) {
	gp := &domain.GlobalParams{}
	gpBytes, err := memory.client.Get(GLOBAL_PARAMS).Bytes()
	if err != nil {
		return nil, err
	}

	err = gp.UnmarshalBinary(gpBytes)
	return gp, err
}

func (memory *redisClient) PutGlobalParams(gp *domain.GlobalParams) error {
	log.Println("Saving global params")
	gpBytes, err := gp.MarshalBinary()
	if err != nil {
		return err
	}

	return memory.client.Set(GLOBAL_PARAMS, gpBytes, 0).Err()
}

func (memory *redisClient) AddActiveVertices(activeVertices []int64) {
	if len(activeVertices) == 0 {
		return
	}
	s := make([]interface{}, len(activeVertices))
	for i, v := range activeVertices {
		s[i] = v
	}
	err := memory.client.SAdd(ACTIVE_VERTICES, s...).Err()
	if err != nil {
		panic(err)
	}
}

func (memory *redisClient) RemoveHaltedVertices(haltedVertices []int64) {
	s := make([]interface{}, len(haltedVertices))
	for i, v := range haltedVertices {
		s[i] = v
	}
	err := memory.client.SRem(ACTIVE_VERTICES, s...).Err()

	if err != nil {
		panic(err)
	}
}

func (memory *redisClient) GetActiveVertices() []int64 {
	var activeVertexIds []int64
	err := memory.client.SMembers(ACTIVE_VERTICES).ScanSlice(&activeVertexIds)

	if err != nil {
		panic(err)
	}

	return activeVertexIds
}

func (memory *redisClient) GetActiveVerticesCount() int64 {
	count, err := memory.client.SCard(ACTIVE_VERTICES).Result()

	if err != nil {
		panic(err)
	}

	return count
}

func (memory *redisClient) SetActiveWorkersCount(count int64) {
	_, err := memory.client.Set(ACTIVE_WORKERS, count, 0).Result()
	if err != nil {
		panic(err)
	}
}

func (memory *redisClient) DecrementActiveWorkersCount() int64 {
	value, err := memory.client.DecrBy(ACTIVE_WORKERS, 1).Result()

	if err != nil {
		panic(err)
	}

	return value
}

func (memory *redisClient) CountReceiversForSuperstep(superstep int64) int64 {
	//TODO should delete set for this superstep after count is retrieved. Use pipeline
	count, err := memory.client.SCard(fmt.Sprintf("%v:%v", MESSAGES, superstep)).Result()
	if err != nil {
		panic(err)
	}

	return count
}

func (memory *redisClient) GetMessageRecipients(superstep int64) []int64 {
	var receivers []int64
	err := memory.client.SMembers(fmt.Sprintf("%v:%v", MESSAGES, superstep)).ScanSlice(&receivers)

	if err != nil {
		panic(err)
	}

	//log.Println("Receivers: ", receivers)
	return receivers
}

func (memory *redisClient) GetMessages(vertexId int64, superstep int64) []interface{} {
	pipe := memory.client.Pipeline()

	resultsCmd := pipe.LRange(fmt.Sprintf("%v:%v:%v", MESSAGES, superstep, vertexId), 0, -1)
	pipe.Unlink(fmt.Sprintf("%v:%v:%v", MESSAGES, superstep, vertexId))

	_, err := pipe.Exec()
	// results, err := memory.client.LRange(fmt.Sprintf("%v:%v:%v", MESSAGES, superstep, vertexId), 0, -1).Result()
	if err != nil {
		panic(err)
	}
	results := resultsCmd.Val()
	messages := make([]interface{}, len(results), len(results))
	for index, message := range results {
		messages[index] = message
	}

	return messages
}

func (memory *redisClient) PutMessageForAllEdges(recipients []domain.Edge, message interface{}, superstep int64) {
	recipientIds := make([]interface{}, len(recipients))
	pipe := memory.client.Pipeline()

	for index, recipient := range recipients {
		recipientIds[index] = recipient.TargetVertexId
		pipe.RPush(fmt.Sprintf("%v:%v:%v", MESSAGES, superstep, recipient.TargetVertexId), message)
	}
	_, err := pipe.Exec()
	if err != nil {
		panic(err)
	}

	memory.client.SAdd(fmt.Sprintf("%v:%v", MESSAGES, superstep), recipientIds...).Err()

	if err != nil {
		panic(err)
	}
}

func (memory *redisClient) PutMessages(recipients []domain.Edge, messages []interface{}, superstep int64) {
	recipientIds := make([]interface{}, len(recipients))
	pipe := memory.client.Pipeline()

	for index, recipient := range recipients {
		recipientIds[index] = recipient.TargetVertexId
		pipe.RPush(fmt.Sprintf("%v:%v:%v", MESSAGES, superstep, recipient.TargetVertexId), messages[index])
	}
	_, err := pipe.Exec()
	if err != nil {
		panic(err)
	}

	memory.client.SAdd(fmt.Sprintf("%v:%v", MESSAGES, superstep), recipientIds...).Err()

	if err != nil {
		panic(err)
	}
}

func (memory *redisClient) PutMessage(recipient int64, message interface{}, superstep int64) {
	pipe := memory.client.Pipeline()

	pipe.RPush(fmt.Sprintf("%v:%v:%v", MESSAGES, superstep, recipient), message)
	pipe.SAdd(fmt.Sprintf("%v:%v", MESSAGES, superstep), recipient)

	_, err := pipe.Exec()

	if err != nil {
		panic(err)
	}
}

func (memory *redisClient) Clear() {
	err := memory.client.ForEachShard(func(client *redis.Client) error {
		return client.FlushDB().Err()
	})
	if err != nil {
		panic(err)
	}

	err = memory.client.ForEachShard(func(client *redis.Client) error {
		purged, err := purgeRedisMemory(client).Result()
		if purged {
			log.Println("Memory purged succesfully for a shard")
		} else {
			log.Println("Memory purged failed for client")
		}
		return err
	})
	if err != nil {
		panic(err)
	}

}

//MEMORY_PURGE only works for redis >4 with jemalloc. It should recover dirty pages after a flushdb. Normally, redis doesn't clear memory even after
//a flushdb, due to malloc behaviour and paging
func purgeRedisMemory(redisdb *redis.Client) *redis.BoolCmd {
	cmd := redis.NewBoolCmd("MEMORY", "PURGE")
	redisdb.Process(cmd)

	return cmd
}

func (memory *redisClient) CreateAggregatorMcl(aggregatorKey string) {
	err := memory.client.SAdd(AGGREGATORS, aggregatorKey).Err()
	if err != nil {
		panic(err)
	}
}

func (memory *redisClient) ResetAggregatorsMcl() {
	aggregators := memory.client.SMembers(AGGREGATORS).Val()
	pipe := memory.client.Pipeline()

	for _, aggregator := range aggregators {
		pipe.Set(aggregator, 0, 0)
	}

	_, err := pipe.Exec()
	if err != nil {
		panic(err)
	}
}

func (memory *redisClient) GetFloatMcl(aggregatorKey string, superstep int64) float64 {
	result, err := memory.client.IncrByFloat(fmt.Sprintf("%v:%v", aggregatorKey, superstep), 0.0).Result()
	if err != nil {
		panic(err)
	}
	return result
}

func (memory *redisClient) AggregateFloatMcl(aggregatorKey string, superstep int64, value float64) {
	_, err := memory.client.IncrByFloat(fmt.Sprintf("%v:%v", aggregatorKey, superstep), value).Result()
	if err != nil {
		panic(err)
	}
}

// var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

// func main() {
// 	mcl, _ := GetMemoryClient(Redis)
// 	flag.Parse()
// 	if *cpuprofile != "" {
// 		f, err := os.Create(*cpuprofile)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		pprof.StartCPUProfile(f)
// 		defer pprof.StopCPUProfile()
// 	}
// 	mcl.GetVertices(mcl.GetAllVertexIds())

// }
