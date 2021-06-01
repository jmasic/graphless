package domain

//StartParams: the initial parameters for the main function invocation
type StartParams struct {
	RunId      string                 `json:runId`
	TestRun    bool                   `json:"testRun"`
	ChunkSize  int64                  `json:"chunkSize"`
	Levels     int64                  `json:"levels"`
	Algorithm  GraphAlgorithm         `json:"algorithm"`
	GraphName  string                 `json:"graphName"`
	ExtraArgs  map[string]interface{} `json:"extraArgs"`
	MaxWorkers int64                  `json:"maxWorkers"`
}
