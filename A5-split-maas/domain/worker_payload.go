package domain

import (
	"encoding/json"
)

//easyjson:json
type WorkerPayload struct {
	VertexIds           []int64                `json:"VertexIds"`
	Superstep           int64                  `json:"superstep"`
	Algorithm           GraphAlgorithm         `json:"algorithm"`
	ExtraArgs           map[string]interface{} `json:"extraArgs"`
	NumberOfVertices    int64                  `json:"numberOfVertices"`
	RunId               string                 `json:"runId"`
	MemoryClientConfig  MemoryClientConfig     `json:"memoryClientConfig"`
	MessageClientConfig MessageClientConfig    `json:"messageClientConfig"`
	StorageClientConfig StorageClientConfig    `json:"storageClientConfig"`
}

func (wp WorkerPayload) MarshalBinary() ([]byte, error) {
	return json.Marshal(wp)
}

func (wp *WorkerPayload) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, &wp)
	return err
}
