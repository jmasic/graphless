package message

import (
	"errors"
	"github.com/devLucian93/thesis-go/domain"
)

type Client interface {
	CountReceiversForSuperstep(superstep int64) int64
	GetMessageRecipients(superstep int64) []int64
	GetMessages(vertexId int64, superstep int64) []interface{}

	PutMessageForAllEdges(recipients []domain.Edge, message interface{}, superstep int64)
	PutMessages(recipients []domain.Edge, messages []interface{}, superstep int64)
	PutMessage(recipient int64, message interface{}, superstep int64)

	Clear()
}

func GetMessageClient(client ClientType, dbConfig domain.DatabaseConfig) (Client, error) {
	switch client {
	case Neo4j:
		return newNeo4jClient(dbConfig)
	case RedisShard:
		return newRedisClient(dbConfig)
	case RedisCluster:
		return newRedisClusterClient(dbConfig)
	}

	return nil, errors.New("Unsupported memory client!")
}
