package domain

import (
	"encoding/json"
)

//easyjson:json
type GlobalParams struct {
	RunId                 string                 `json:"runId"`
	Superstep             int64                  `json:"superstep"`
	NumberOfVertices      int64                  `json:"numberOfVertices"`
	NumberOfEdges         int64                  `json:"numberOfEdges"`
	NumberOfBuckets       int64                  `json:"numberOfBuckets"`
	ChunkSize             int64                  `json:"chunkSize"`
	Finished              bool                   `json:"finished"`
	DataIngestionDuration int64                  `json:"dataIngestionDuration"`
	ExecutionDuration     int64                  `json:"executionDuration"`
	Algorithm             GraphAlgorithm         `json:"algorithm"`
	GraphName             string                 `json:"graphName"`
	ExtraArgs             map[string]interface{} `json:"extraArgs"`
	MaxWorkers            int64                  `json:"maxWorkers"`
}

func (gp GlobalParams) MarshalBinary() ([]byte, error) {
	return json.Marshal(gp)
}

func (gp *GlobalParams) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &gp); err != nil {
		return err
	}
	return nil
}

func (gp GlobalParams) MarshalExtraArgs() ([]byte, error) {
	return json.Marshal(gp.ExtraArgs)
}

func (gp *GlobalParams) UnmarshalExtraArgs(data []byte) error {
	if err := json.Unmarshal(data, &gp.ExtraArgs); err != nil {
		return err
	}
	return nil
}

func (gp GlobalParams) String() string {
	b, _ := json.Marshal(gp)
	return string(b[:])
}
