package memory

import (
	"errors"
	"github.com/devLucian93/thesis-go/domain"
)

type Client interface {
	GetVertices(intKeys []int64) <-chan []domain.Vertex
	GetAllVertexIds() []int64
	CreateVertices(vertices []domain.Vertex)
	SaveVertices(vertices []domain.Vertex)

	GetGlobalParams() (*domain.GlobalParams, error)
	PutGlobalParams(*domain.GlobalParams) error

	SetActiveWorkersCount(count int64)
	DecrementActiveWorkersCount() int64

	GetFloatMcl(aggregatorKey string, superstep int64) float64
	AggregateFloatMcl(aggregatorKey string, superstep int64, value float64)
	Clear()
}

func GetMemoryClient(client ClientType, config domain.DatabaseConfig) (Client, error) {
	switch client {
	case Neo4j:
		return newNeo4jClient(config)
	case Redis:
		return newRedisClient(config)
	}

	return nil, errors.New("Unsupported memory client!")
}
