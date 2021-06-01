package clients

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/devLucian93/thesis-go/domain"
	"github.com/devLucian93/thesis-go/utils"
	"github.com/go-redis/redis"
)

type MemoryClient interface {
	VertexRange(startKey int64, endKey int64) []domain.Vertex
	GetVertices(intKeys []int64) []domain.Vertex
	GetAllVertexIds() []int64
	PutVertex(vertex *domain.Vertex)
	PutVertices(vertices []domain.Vertex)
	DeleteVertex(key string) error

	GetGlobalParams() (*domain.GlobalParams, error)
	PutGlobalParams(*domain.GlobalParams) error

	SetActiveWorkersCount(count int64)
	DecrementActiveWorkersCount(finishedWorkers int64) int64

	AddActiveVertices(activeVertices []int64)
	GetActiveVertices() []int64
	RemoveHaltedVertices(haltedVertices []int64)
	GetActiveVerticesCount() int64

	CountReceiversForSuperstep(superstep int64) int64

	GetMessageRecipients(superstep int64) []int64
	GetMessages(vertexId int64, superstep int64) []interface{}
	PutMessageForAllEdges(recipients []domain.Edge, message interface{}, superstep int64)
	PutMessages(recipients []domain.Edge, messages []interface{}, superstep int64)
	PutMessage(recipient int64, message interface{}, superstep int64)

	CreateAggregatorMcl(aggregatorKey string)
	AggregateFloatMcl(aggregatorKey string, value float64, superstep int64) float64
	ResetAggregatorsMcl()
	Clear()
}

func GetMemoryClient(client MemoryClientType) (MemoryClient, error) {
	switch client {
	case REDIS:
		client := redis.NewRing(&redis.RingOptions{
			//r4.large
			Addrs: map[string]string{
				//TRY AWS PRIVATE IP TO AVOID REGIONAL DATA TRANSFER CHARGES!
				"redis-6379": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6379",
				"redis-6380": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6380",
				"redis-6381": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6381",
				"redis-6382": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6382",
				"redis-6383": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6383",
				"redis-6384": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6384",
				"redis-6385": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6385",
				"redis-6386": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6386",
				"redis-6387": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6387",
				"redis-6388": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6388",
				"redis-6389": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6389",
				"redis-6390": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6390",
				"redis-6391": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6391",
				"redis-6392": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6392",
				"redis-6393": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6393",
				"redis-6394": "ec2-18-222-182-74.us-east-2.compute.amazonaws.com:6394",
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
		redis := &RedisClient{client}
		return redis, nil
	}

	return nil, errors.New("Unsupported memory client!")
}

const (
	GLOBAL_PARAMS   = "globalParams"
	VERTEX          = "vertices"
	MESSAGES        = "messages"
	ACTIVE_VERTICES = "activeVertices"
	ACTIVE_WORKERS  = "activeWorkers"
	AGGREGATORS     = "aggregators"
)

type RedisClient struct {
	client *redis.Ring
}

//TODO Old impl, not working
func (memory *RedisClient) VertexRange(startKey int64, endKey int64) []domain.Vertex {
	var keys []string
	for i := startKey; i < endKey; i++ {
		keys = append(keys, fmt.Sprintf("%v:%v", VERTEX, i))
	}

	vertices := getMultipleKeys(keys, memory)

	return vertices
}

func (memory *RedisClient) GetAllVertexIds() []int64 {
	var vertexIds []int64
	err := memory.client.SMembers(VERTEX).ScanSlice(&vertexIds)

	if err != nil {
		panic(err)
	}

	return vertexIds
}

func (memory *RedisClient) GetVertices(intKeys []int64) []domain.Vertex {
	var keys []string
	for _, key := range intKeys {
		keys = append(keys, fmt.Sprintf("%v:%v", VERTEX, key))
	}
	vertices := make([]domain.Vertex, 0, len(intKeys))

	vertices = append(vertices, getMultipleKeys(keys, memory)...)

	return vertices
}

func getMultipleKeys(keys []string, memory *RedisClient) []domain.Vertex {
	//TODO Readapt for getMessageReceivers
	verticesCmd := make([]*redis.StringCmd, len(keys), len(keys))
	vertices := make([]domain.Vertex, 0, len(keys))

	pipe := memory.client.Pipeline()

	for i := 0; i < len(keys); i++ {
		verticesCmd[i] = pipe.Get(keys[i])
	}
	_, err := pipe.Exec()

	if err != nil {
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

func (memory *RedisClient) PutVertex(vertex *domain.Vertex) {

	pipe := memory.client.Pipeline()

	pipe.SAdd(VERTEX, vertex.Id)
	jsonBytes, _ := vertex.MarshalBinary()
	pipe.Set(fmt.Sprintf("%v:%v", VERTEX, vertex.Id), utils.ZLibCompress(jsonBytes), 0)

	_, err := pipe.Exec()

	if err != nil {
		panic(err)
	}
}

func (memory *RedisClient) PutVertices(vertices []domain.Vertex) {
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

func (memory *RedisClient) DeleteVertex(key string) error {
	//to implement
	return errors.New("Unimplemented method 'DeleteVertex'")
}

func (memory *RedisClient) GetGlobalParams() (*domain.GlobalParams, error) {
	gp := &domain.GlobalParams{}
	gpBytes, err := memory.client.Get(GLOBAL_PARAMS).Bytes()
	if err != nil {
		return nil, err
	}

	err = gp.UnmarshalBinary(gpBytes)
	return gp, err
}

func (memory *RedisClient) PutGlobalParams(gp *domain.GlobalParams) error {
	log.Println("Saving global params")
	gpBytes, err := gp.MarshalBinary()
	if err != nil {
		return err
	}

	return memory.client.Set(GLOBAL_PARAMS, gpBytes, 0).Err()
}

func (memory *RedisClient) AddActiveVertices(activeVertices []int64) {

	s := make([]interface{}, len(activeVertices))
	for i, v := range activeVertices {
		s[i] = v
	}
	err := memory.client.SAdd(ACTIVE_VERTICES, s...).Err()
	if err != nil {
		panic(err)
	}
}

func (memory *RedisClient) RemoveHaltedVertices(haltedVertices []int64) {

	s := make([]interface{}, len(haltedVertices))
	for i, v := range haltedVertices {
		s[i] = v
	}
	err := memory.client.SRem(ACTIVE_VERTICES, s...).Err()

	if err != nil {
		panic(err)
	}
}

func (memory *RedisClient) GetActiveVertices() []int64 {
	var activeVertexIds []int64
	err := memory.client.SMembers(ACTIVE_VERTICES).ScanSlice(&activeVertexIds)

	if err != nil {
		panic(err)
	}

	return activeVertexIds
}

func (memory *RedisClient) GetActiveVerticesCount() int64 {

	count, err := memory.client.SCard(ACTIVE_VERTICES).Result()

	if err != nil {
		panic(err)
	}

	return count
}

func (memory *RedisClient) SetActiveWorkersCount(count int64) {
	_, err := memory.client.Set(ACTIVE_WORKERS, count, 0).Result()
	if err != nil {
		panic(err)
	}
}

func (memory *RedisClient) DecrementActiveWorkersCount(finishedComputations int64) int64 {

	value, err := memory.client.DecrBy(ACTIVE_WORKERS, finishedComputations).Result()

	if err != nil {
		panic(err)
	}

	return value
}

func (memory *RedisClient) CountReceiversForSuperstep(superstep int64) int64 {
	//TODO should delete set for this superstep after count is retrieved. Use pipeline
	count, err := memory.client.SCard(fmt.Sprintf("%v:%v", MESSAGES, superstep)).Result()
	if err != nil {
		panic(err)
	}

	return count

}

func (memory *RedisClient) GetMessageRecipients(superstep int64) []int64 {
	var receivers []int64
	err := memory.client.SMembers(fmt.Sprintf("%v:%v", MESSAGES, superstep)).ScanSlice(&receivers)

	if err != nil {
		panic(err)
	}

	//log.Println("Receivers: ", receivers)
	return receivers
}

func (memory *RedisClient) GetMessages(vertexId int64, superstep int64) []interface{} {

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

func (memory *RedisClient) PutMessageForAllEdges(recipients []domain.Edge, message interface{}, superstep int64) {

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

//Would work if the client knew not to wait for a reply. Sent a message to the author of the redis client
func toggleServerReplies(redisdb *redis.Ring, key string) *redis.BoolCmd {
	cmd := redis.NewBoolCmd("client", "reply", key)
	redisdb.Process(cmd)
	return cmd
}

func (memory *RedisClient) PutMessages(recipients []domain.Edge, messages []interface{}, superstep int64) {

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

func (memory *RedisClient) PutMessage(recipient int64, message interface{}, superstep int64) {
	pipe := memory.client.Pipeline()

	pipe.RPush(fmt.Sprintf("%v:%v:%v", MESSAGES, superstep, recipient), message)
	pipe.SAdd(fmt.Sprintf("%v:%v", MESSAGES, superstep), recipient)

	_, err := pipe.Exec()

	if err != nil {
		panic(err)
	}

}

func (memory *RedisClient) Clear() {
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

func (memory *RedisClient) CreateAggregatorMcl(aggregatorKey string) {
	err := memory.client.SAdd(AGGREGATORS, aggregatorKey).Err()
	if err != nil {
		panic(err)
	}
}

func (memory *RedisClient) ResetAggregatorsMcl() {
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

func (memory *RedisClient) AggregateFloatMcl(aggregatorKey string, value float64, superstep int64) float64 {
	result, err := memory.client.IncrByFloat(fmt.Sprintf("%v:%v", aggregatorKey, superstep), value).Result()
	if err != nil {
		panic(err)
	}
	return result
}

// var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

// func main() {
// 	mcl, _ := GetMemoryClient(REDIS)
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
