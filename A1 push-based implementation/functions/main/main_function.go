package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"time"

	lambdaStarter "github.com/aws/aws-lambda-go/lambda"
	"github.com/devLucian93/thesis-go/clients"
	"github.com/devLucian93/thesis-go/domain"
	"github.com/devLucian93/thesis-go/utils"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
)

var storageClient clients.StorageClient
var memoryClient clients.MemoryClient
var functionClient clients.FunctionClient
var functionName = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
var log *logrus.Logger

func init() {
	log = utils.GetLogger()
	mcl, err := clients.GetMemoryClient(clients.REDIS)
	memoryClient = mcl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}

	scl, err := clients.GetStorageClient(clients.S3)
	storageClient = scl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}

	fcl, err := clients.GetFunctionClient(clients.LAMBDA)
	functionClient = fcl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}

}

func mainFunction(params domain.StartParams) (msg string, e error) {
	// dateTimeFormat := "20060102-150405"
	// runId := fmt.Sprintf("%s-%s-%s", time.Now().Format(dateTimeFormat), params.GraphName, params.Algorithm)

	log.WithFields(logrus.Fields{
		"runId": params.RunId,
		"tag":   "STARTED",
	}).Info("Started graph processing engine with params: ", params)

	memoryClient.Clear()

	var globalParams *domain.GlobalParams
	if params.TestRun {

		vertices := buildBinaryTree(params.Levels, int64(params.ExtraArgs["sourceVertex"].(float64)))
		globalParams = &domain.GlobalParams{
			NumberOfVertices:      int64(math.Pow(2, float64(params.Levels+1)) - 1),
			ChunkSize:             params.ChunkSize,
			DataIngestionDuration: time.Now().UnixNano(),
			Superstep:             0,
			NumberOfBuckets:       1,
			Algorithm:             params.Algorithm,
			GraphName:             params.GraphName,
			ExtraArgs:             params.ExtraArgs,
			RunId:                 params.RunId,
			MaxWorkers:            params.MaxWorkers,
		}
		memoryClient.PutGlobalParams(globalParams)
		var activeVertexIds []int64
		for _, vertex := range vertices {
			activeVertexIds = append(activeVertexIds, vertex.Id)
		}
		memoryClient.PutVertices(vertices)
		memoryClient.AddActiveVertices(activeVertexIds)
	} else {
		globalParams = &domain.GlobalParams{
			ChunkSize:             params.ChunkSize,
			DataIngestionDuration: time.Now().UnixNano(),
			Superstep:             0,
			Algorithm:             params.Algorithm,
			GraphName:             params.GraphName,
			ExtraArgs:             params.ExtraArgs,
			RunId:                 params.RunId,
			MaxWorkers:            params.MaxWorkers,
		}
		result, err := storageClient.Get(fmt.Sprintf("-%s-%s", globalParams.GraphName, "properties"))
		if err != nil {
			panic(err)
		}
		globalParams.UnmarshalBinary(result)
		memoryClient.PutGlobalParams(globalParams)

		for i := int64(0); i < globalParams.NumberOfBuckets; i++ {
			var activeVertexIds []int64
			verticesChunk := readGraphChunkCompressed(fmt.Sprintf("-%s-%d", globalParams.GraphName, i))
			for index, vertex := range verticesChunk {
				vertex.Value = math.MaxInt64
				verticesChunk[index] = vertex
				activeVertexIds = append(activeVertexIds, vertex.Id)
			}

			insertGraphChunk(verticesChunk, activeVertexIds)
		}
	}

	invokeOrchestratorFunction(domain.ORCHESTRATOR_INIT)

	return "Hello from Go!", nil
}

func main() {
	// tr := trace.NewTrace()
	// m := metric.NewBuilder().Build()
	// t := thundra.NewBuilder().
	// 	AddPlugin(tr).
	// 	AddPlugin(m).
	// 	Build()
	// lambdaStarter.Start(thundra.Wrap(mainFunction, t))
	lambdaStarter.Start(mainFunction)
}

func insertGraphChunk(verticesChunk []domain.Vertex, activeVertexIds []int64) {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"tag":          "TEST",
			"pureDuration": fmt.Sprintln(utils.MeasureDurationMs(start), "ms"),
		}).Info("Logging insertGraphChunk duration")
	}(time.Now())

	memoryClient.PutVertices(verticesChunk)
	memoryClient.AddActiveVertices(activeVertexIds)
}

func buildBinaryTree(levels int64, sourceVertexId int64) []domain.Vertex {
	var vertices []domain.Vertex
	id := int64(1)
	vertices = append(vertices, domain.Vertex{
		Id:    sourceVertexId,
		Edges: []domain.Edge{domain.Edge{TargetVertexId: 1, Value: 1}, domain.Edge{TargetVertexId: 2, Value: 2}},
		Value: math.MaxInt64,
	})
	for i := int64(1); i <= levels; i++ {
		for j := int64(0); j < int64(math.Pow(2, float64(i))); j++ {
			if i < levels {
				vertices = append(vertices, domain.Vertex{
					Id:    id,
					Edges: []domain.Edge{domain.Edge{TargetVertexId: 2*id + 1, Value: 1}, domain.Edge{TargetVertexId: 2*id + 2, Value: 2}},
					Value: math.MaxInt64})
			} else {
				vertices = append(vertices, domain.Vertex{Id: id, Edges: []domain.Edge{}, Value: math.MaxInt64})
			}
			id++
		}
	}

	return vertices
}

func readGraphChunk(key string) []domain.Vertex {

	buf, err := storageClient.Get(key)
	if err != nil {
		log.Println("Error")
		os.Exit(1)
	}
	v := &domain.VertexList{}
	easyjson.Unmarshal(buf, v)
	return []domain.Vertex(*v)
}

func readGraphChunkCompressed(key string) []domain.Vertex {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"tag":          "TEST",
			"pureDuration": fmt.Sprintln(utils.MeasureDurationMs(start), "ms"),
		}).Info("Logging read graph chunk compressed duration")
	}(time.Now())
	var vertices []domain.Vertex

	compressedBuf, err := getChunkFromStorage(key)
	if err != nil {
		log.Println("Error")
		os.Exit(1)
	}

	zlibReader, err := zlib.NewReader(bytes.NewReader(compressedBuf))
	defer zlibReader.Close()

	if err != nil {
		panic(err)
	}
	decompressedBuf, err := readDecompressed(zlibReader)

	if err != nil {
		log.Println("Error decompressing chunk ", key)
		panic(err)
	}
	vertices, err = unmarshalJson(decompressedBuf)
	if err != nil {
		panic(err)
	}

	return vertices
}

func getChunkFromStorage(key string) ([]byte, error) {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"tag":          "TEST",
			"pureDuration": fmt.Sprintln(utils.MeasureDurationMs(start), "ms"),
		}).Info("Logging get chunk from storage duration")
	}(time.Now())
	return storageClient.Get(key)
}

func readDecompressed(zlibReader io.ReadCloser) ([]byte, error) {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"tag":          "TEST",
			"pureDuration": fmt.Sprintln(utils.MeasureDurationMs(start), "ms"),
		}).Info("Logging decompressing duration")
	}(time.Now())
	return ioutil.ReadAll(zlibReader)
}

func unmarshalJson(decompressedBuf []byte) ([]domain.Vertex, error) {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"tag":          "TEST",
			"pureDuration": fmt.Sprintln(utils.MeasureDurationMs(start), "ms"),
		}).Info("Logging unmarshalling json duration")
	}(time.Now())

	v := &domain.VertexList{}
	err := easyjson.Unmarshal(decompressedBuf, v)
	return []domain.Vertex(*v), err
}

func invokeOrchestratorFunction(message string) {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"tag":          "TEST",
			"pureDuration": fmt.Sprintln(utils.MeasureDurationMs(start), "ms"),
		}).Info("Logging Invoking orchestrator function duration")
	}(time.Now())
	log.Println("Invoking orchestrator function")
	binaryPayload, err := json.Marshal(domain.OrchestratorPayload{message})

	if err != nil {
		panic(err)
	}

	err = functionClient.InvokeFunction(os.Getenv("ORCHESTRATOR_FUNCTION_NAME"), binaryPayload)

	if err != nil {
		panic(err)
	}

}
