package domain

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

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

func ReadStartParamsFromFile(fileName string) (StartParams, error) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	var startParams StartParams
	json.Unmarshal(fileBytes, &startParams)
	return startParams, nil
}
