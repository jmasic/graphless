package memory

import (
	"github.com/devLucian93/thesis-go/domain"
)

type inMemoryClient struct {
}

func newInMemoryClient() (Client, error) {
	inMemoryClient := &inMemoryClient{}
	return inMemoryClient, nil
}

func (memory *inMemoryClient) VertexRange(startKey int64, endKey int64) []domain.Vertex {
	panic("VertexRange NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) GetAllVertexIds() []int64 {
	panic("GetAllVertexIds NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) GetVertices(intKeys []int64) []domain.Vertex {
	panic("GetVertices NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) PutVertex(vertex *domain.Vertex) {
	panic("PutVertex NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) PutVertices(vertices []domain.Vertex) {
	panic("PutVertices NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) DeleteVertex(key string) error {
	//to implementDeleteVertex	return errors.New("Unimplemented method 'DeleteVertex'")
	panic("DeleteVertex NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) GetGlobalParams() (*domain.GlobalParams, error) {
	panic("GetGlobalParams NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) PutGlobalParams(gp *domain.GlobalParams) error {
	panic("PutGlobalParams NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) AddActiveVertices(activeVertices []int64) {
	panic("AddActiveVertices NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) RemoveHaltedVertices(haltedVertices []int64) {
	panic("RemoveHaltedVertices NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) GetActiveVertices() []int64 {
	panic("GetActiveVertices NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) GetActiveVerticesCount() int64 {
	panic("GetActiveVerticesCount NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) SetActiveWorkersCount(count int64) {
	panic("SetActiveWorkersCount NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) DecrementActiveWorkersCount() int64 {
	panic("DecrementActiveWorkersCount NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) CountReceiversForSuperstep(superstep int64) int64 {
	panic("CountReceiversForSuperstep NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) GetMessageRecipients(superstep int64) []int64 {
	panic("GetMessageRecipients NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) GetMessages(vertexId int64, superstep int64) []interface{} {
	panic("GetMessages NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) PutMessageForAllEdges(recipients []domain.Edge, message interface{}, superstep int64) {
	panic("PutMessageForAllEdges NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) PutMessages(recipients []domain.Edge, messages []interface{}, superstep int64) {
	panic("PutMessages NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) PutMessage(recipient int64, message interface{}, superstep int64) {
	panic("PutMessage NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) Clear() {
	// no-op
}

func (memory *inMemoryClient) CreateAggregatorMcl(aggregatorKey string) {
	panic("CreateAggregatorMcl NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) ResetAggregatorsMcl() {
	panic("ResetAggregatorsMcl NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) GetFloatMcl(aggregatorKey string, superstep int64) float64 {
	panic("GetFloatMcl NOT IMPLEMENTED YET")
}

func (memory *inMemoryClient) AggregateFloatMcl(aggregatorKey string, superstep int64, value float64) {
	panic("AggregateFloatMcl NOT IMPLEMENTED YET")
}
