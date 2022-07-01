package domain

import (
	"encoding/json"
)

//easyjson:json
type WorkerPayload struct {
	VertexIds        []int64
	Superstep        int64
	Algorithm        GraphAlgorithm
	ExtraArgs        map[string]interface{}
	NumberOfVertices int64
	RunId            string
}

func (wp WorkerPayload) MarshalBinary() ([]byte, error) {
	return json.Marshal(wp)
}

func (wp *WorkerPayload) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &wp); err != nil {
		return err
	}

	return nil
}
