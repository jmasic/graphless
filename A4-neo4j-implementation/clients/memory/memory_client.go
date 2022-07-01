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

	CountReceiversForSuperstep(superstep int64) int64

	GetMessageRecipients(superstep int64) []int64
	GetMessages(vertexId int64, superstep int64) []interface{}
	PutMessageForAllEdges(recipients []domain.Edge, message interface{}, superstep int64)
	PutMessages(recipients []domain.Edge, messages []interface{}, superstep int64)
	PutMessage(recipient int64, message interface{}, superstep int64)

	GetFloatMcl(aggregatorKey string, superstep int64) float64
	AggregateFloatMcl(aggregatorKey string, superstep int64, value float64)
	Clear()
}

func GetMemoryClient(client ClientType) (Client, error) {
	switch client {
	case Neo4j:
		return newNeo4jClient()
	case Redis:
		return newRedisClient()
	}

	return nil, errors.New("Unsupported memory client!")
}
