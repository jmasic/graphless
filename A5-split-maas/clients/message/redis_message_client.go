package message

import (
	"fmt"
	"github.com/devLucian93/thesis-go/domain"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	REDIS_KEY_MESSAGES = "messages"
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

func (memory *redisClient) CountReceiversForSuperstep(superstep int64) int64 {
	//TODO should delete set for this superstep after count is retrieved. Use pipeline
	count, err := memory.client.SCard(fmt.Sprintf("%v:%v", REDIS_KEY_MESSAGES, superstep)).Result()
	if err != nil {
		panic(err)
	}

	return count
}

func (memory *redisClient) GetMessageRecipients(superstep int64) []int64 {
	var receivers []int64
	err := memory.client.SMembers(fmt.Sprintf("%v:%v", REDIS_KEY_MESSAGES, superstep)).ScanSlice(&receivers)

	if err != nil {
		panic(err)
	}

	//log.Println("Receivers: ", receivers)
	return receivers
}

func (memory *redisClient) GetMessages(vertexId int64, superstep int64) []interface{} {
	return memory.GetMessagesWithRetries(vertexId, superstep, 0)
}

func (memory *redisClient) GetMessagesWithRetries(vertexId int64, superstep int64, retries int64) []interface{} {
	pipe := memory.client.Pipeline()

	resultsCmd := pipe.LRange(fmt.Sprintf("%v:%v:%v", REDIS_KEY_MESSAGES, superstep, vertexId), 0, -1)
	pipe.Unlink(fmt.Sprintf("%v:%v:%v", REDIS_KEY_MESSAGES, superstep, vertexId))

	_, err := pipe.Exec()
	if err != nil {
		if retries < 10 {
			time.Sleep(time.Duration(rand.Intn(2_000)) * time.Millisecond)
			return memory.GetMessagesWithRetries(vertexId, superstep, retries+1)
		}
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
		pipe.RPush(fmt.Sprintf("%v:%v:%v", REDIS_KEY_MESSAGES, superstep, recipient.TargetVertexId), message)
	}
	_, err := pipe.Exec()
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			memory.PutMessageForAllEdges(recipients, message, superstep)
			return
		}
		panic(err)
	}

	memory.client.SAdd(fmt.Sprintf("%v:%v", REDIS_KEY_MESSAGES, superstep), recipientIds...).Err()

	if err != nil {
		panic(err)
	}
}

func (memory *redisClient) PutMessages(recipients []domain.Edge, messages []interface{}, superstep int64) {
	recipientIds := make([]interface{}, len(recipients))
	pipe := memory.client.Pipeline()

	for index, recipient := range recipients {
		recipientIds[index] = recipient.TargetVertexId
		pipe.RPush(fmt.Sprintf("%v:%v:%v", REDIS_KEY_MESSAGES, superstep, recipient.TargetVertexId), messages[index])
	}
	_, err := pipe.Exec()
	if err != nil {
		panic(err)
	}

	memory.client.SAdd(fmt.Sprintf("%v:%v", REDIS_KEY_MESSAGES, superstep), recipientIds...).Err()

	if err != nil {
		panic(err)
	}
}

func (memory *redisClient) PutMessage(recipient int64, message interface{}, superstep int64) {
	pipe := memory.client.Pipeline()

	pipe.RPush(fmt.Sprintf("%v:%v:%v", REDIS_KEY_MESSAGES, superstep, recipient), message)
	pipe.SAdd(fmt.Sprintf("%v:%v", REDIS_KEY_MESSAGES, superstep), recipient)

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
