package message

import (
	"fmt"
	"github.com/devLucian93/thesis-go/domain"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

type redisClusterClient struct {
	client *redis.ClusterClient
}

func newRedisClusterClient(config domain.DatabaseConfig) (Client, error) {
	// TODO: This didn't perform well, so we haven't invested time in refactoring the database config to accommodate
	// 		 the options for redis clusters. Ideally, it should be something like a list of URIs
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			// memory db
			//"redis-small.qttbzi.clustercfg.memorydb.us-east-2.amazonaws.com:6379",
			// elasticache
			"redis-cache-0001-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0002-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0003-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0004-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0005-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0006-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0007-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0008-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0009-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0010-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0011-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0012-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0013-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0014-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0015-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
			"redis-cache-0016-001.qttbzi.0001.use2.cache.amazonaws.com:6379",
		},
		MaxRedirects:    8,
		MaxRetries:      20,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
		DialTimeout:     20 * time.Second,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		PoolSize:        20,
		Password:        "",
	})
	redisCli := &redisClusterClient{clusterClient}
	return redisCli, nil
}

func (memory *redisClusterClient) CountReceiversForSuperstep(superstep int64) int64 {
	//TODO should delete set for this superstep after count is retrieved. Use pipeline
	count, err := memory.client.SCard(fmt.Sprintf("%v:%v", REDIS_KEY_MESSAGES, superstep)).Result()
	if err != nil {
		panic(err)
	}

	return count
}

func (memory *redisClusterClient) GetMessageRecipients(superstep int64) []int64 {
	var receivers []int64
	err := memory.client.SMembers(fmt.Sprintf("%v:%v", REDIS_KEY_MESSAGES, superstep)).ScanSlice(&receivers)

	if err != nil {
		panic(err)
	}

	//log.Println("Receivers: ", receivers)
	return receivers
}

func (memory *redisClusterClient) GetMessages(vertexId int64, superstep int64) []interface{} {
	return memory.GetMessagesWithRetries(vertexId, superstep, 0)
}

func (memory *redisClusterClient) GetMessagesWithRetries(vertexId int64, superstep int64, retriesCount int64) []interface{} {
	pipe := memory.client.Pipeline()

	resultsCmd := pipe.LRange(fmt.Sprintf("%v:%v:%v", REDIS_KEY_MESSAGES, superstep, vertexId), 0, -1)
	pipe.Unlink(fmt.Sprintf("%v:%v:%v", REDIS_KEY_MESSAGES, superstep, vertexId))

	_, err := pipe.Exec()
	if err != nil {
		if retriesCount < 5 {
			time.Sleep(500 * time.Millisecond)
			return memory.GetMessagesWithRetries(vertexId, superstep, retriesCount+1)
		} else {
			panic(err)
		}
	}
	results := resultsCmd.Val()
	messages := make([]interface{}, len(results), len(results))
	for index, message := range results {
		messages[index] = message
	}

	return messages
}

func (memory *redisClusterClient) PutMessageForAllEdges(recipients []domain.Edge, message interface{}, superstep int64) {
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

func (memory *redisClusterClient) PutMessages(recipients []domain.Edge, messages []interface{}, superstep int64) {
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

func (memory *redisClusterClient) PutMessage(recipient int64, message interface{}, superstep int64) {
	pipe := memory.client.Pipeline()

	pipe.RPush(fmt.Sprintf("%v:%v:%v", REDIS_KEY_MESSAGES, superstep, recipient), message)
	pipe.SAdd(fmt.Sprintf("%v:%v", REDIS_KEY_MESSAGES, superstep), recipient)

	_, err := pipe.Exec()

	if err != nil {
		panic(err)
	}
}

func (memory *redisClusterClient) Clear() {
	cmd := redis.NewBoolCmd("FLUSHALL")
	err := memory.client.Process(cmd)
	if err != nil {
		panic(err)
	}
}
