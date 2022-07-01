package memory

import (
	"fmt"
	"github.com/devLucian93/thesis-go/domain"
	"github.com/devLucian93/thesis-go/utils"
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const (
	GLOBAL_PARAMS               = "globalParams"
	VERTEX                      = "vertices"
	ACTIVE_WORKERS              = "activeWorkers"
	REDIS_VERTEX_CHUNK_SIZE int = 200
)

type redisClient struct {
	client *redis.Ring
}

func newRedisClient(config domain.DatabaseConfig) (Client, error) {
	client := redis.NewRing(&redis.RingOptions{
		// r4.large
		Addrs:           generateRedisShardAddresses(config),
		MaxRetries:      20,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
		DialTimeout:     30 * time.Second,
		ReadTimeout:     20 * time.Second,
		WriteTimeout:    20 * time.Second,
		PoolSize:        20,
		Password:        config.Password,
		DB:              0,
	})
	redisCli := &redisClient{client}
	return redisCli, nil
}

func generateRedisShardAddresses(config domain.DatabaseConfig) map[string]string {
	ipAddress := config.Ip
	basePort := config.Port
	shardsCount := config.ShardsCount
	addresses := make(map[string]string)
	for i := 0; i < shardsCount; i++ {
		port := strconv.Itoa(basePort + i)
		key := "redis-" + port
		addresses[key] = ipAddress + ":" + port
	}
	return addresses
}

func (memory *redisClient) GetAllVertexIds() []int64 {
	var vertexIds []int64
	err := memory.client.SMembers(VERTEX).ScanSlice(&vertexIds)

	if err != nil {
		panic(err)
	}

	return vertexIds
}

func (memory *redisClient) GetVertices(vertexIds []int64) <-chan []domain.Vertex {
	channelSize := len(vertexIds)/REDIS_VERTEX_CHUNK_SIZE + 1
	vertexChannel := make(chan []domain.Vertex, channelSize)
	if len(vertexIds) == 0 {
		return vertexChannel
	}

	go func() {
		defer close(vertexChannel)
		for i := 0; len(vertexIds) > 0; i++ {
			nextChunkSize := REDIS_VERTEX_CHUNK_SIZE
			if len(vertexIds) < nextChunkSize {
				nextChunkSize = len(vertexIds)
			}
			keys := make([]string, nextChunkSize)
			for j := 0; j < nextChunkSize; j++ {
				keys[j] = fmt.Sprintf("%v:%v", VERTEX, vertexIds[j])
			}
			vertices := getMultipleKeys(keys, memory, 0)
			vertexChannel <- vertices

			newStart := REDIS_VERTEX_CHUNK_SIZE
			if len(vertexIds) < newStart {
				newStart = len(vertexIds)
			}
			vertexIds = vertexIds[newStart:]
		}
		vertexChannel <- nil
	}()

	return vertexChannel
}

func getMultipleKeys(keys []string, memory *redisClient, retries int) []domain.Vertex {
	//TODO Readapt for getMessageReceivers
	verticesCmd := make([]*redis.StringCmd, len(keys), len(keys))
	vertices := make([]domain.Vertex, 0, len(keys))

	pipe := memory.client.Pipeline()

	for i := 0; i < len(keys); i++ {
		verticesCmd[i] = pipe.Get(keys[i])
	}
	_, err := pipe.Exec()

	if err != nil {
		log.Println("Error, cannot get keys", keys)
		if retries < 10 {
			time.Sleep(time.Duration(rand.Intn(2_000)) * time.Millisecond)
			return getMultipleKeys(keys, memory, retries+1)
		}
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

func (memory *redisClient) CreateVertices(vertices []domain.Vertex) {
	memory.SaveVertices(vertices)
}

func (memory *redisClient) SaveVertices(vertices []domain.Vertex) {
	if len(vertices) == 0 {
		return
	}

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

	err = memory.client.SAdd(VERTEX, vertexIds...).Err()
	if err != nil {
		panic(err)
	}
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
