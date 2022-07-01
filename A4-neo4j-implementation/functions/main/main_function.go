package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"github.com/devLucian93/thesis-go/clients/function"
	functionapi "github.com/devLucian93/thesis-go/clients/function/api"
	"github.com/devLucian93/thesis-go/clients/memory"
	"github.com/devLucian93/thesis-go/clients/storage"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sync"
	"time"

	lambdaStarter "github.com/aws/aws-lambda-go/lambda"
	"github.com/devLucian93/thesis-go/domain"
	"github.com/devLucian93/thesis-go/utils"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
)

var storageClient storage.Client
var memoryClient memory.Client
var functionClient functionapi.Client
var log *logrus.Logger

func init() {
	local := utils.IsLocal()
	log = utils.GetLogger()
	memoryClientType := (map[bool]memory.ClientType{true: memory.Neo4j, false: memory.Neo4j})[local]
	mcl, err := memory.GetMemoryClient(memoryClientType)
	memoryClient = mcl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}

	storageClientType := (map[bool]storage.ClientType{true: storage.LocalFileSystem, false: storage.S3})[local]
	scl, err := storage.GetStorageClient(storageClientType)
	storageClient = scl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}

	functionClientType := (map[bool]functionapi.ClientType{true: functionapi.GoFunction, false: functionapi.AwsLambda})[local]
	fcl, err := function.GetFunctionClient(functionClientType)
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
		err := memoryClient.PutGlobalParams(globalParams)
		if err != nil {
			panic(err)
		}
		var activeVertexIds []int64
		for _, vertex := range vertices {
			activeVertexIds = append(activeVertexIds, vertex.Id)
		}
		memoryClient.CreateVertices(vertices)
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
		propertiesFileName := fmt.Sprintf("-%s-%s", globalParams.GraphName, "properties")
		result, err := storageClient.Get(propertiesFileName)
		if err != nil {
			panic(err)
		}
		err2 := globalParams.UnmarshalBinary(result)
		if err2 != nil {
			panic(err2)
		}
		err3 := memoryClient.PutGlobalParams(globalParams)
		if err3 != nil {
			panic(err3)
		}
		log.Info("Saved global params in memory")

		loadGraphInMemory(globalParams)
		log.Info("Inserted graph in memory")
	}

	invokeOrchestratorFunction(domain.ORCHESTRATOR_INIT)

	return "Hello from Go!", nil
}

func main() {
	if utils.IsLocal() {
		params, err := domain.ReadStartParamsFromFile("local_payload.json")
		if err != nil {
			panic(err)
		}
		mainFunction(params)
	} else {
		lambdaStarter.Start(mainFunction)
	}
}

func loadGraphInMemory(globalParams *domain.GlobalParams) {
	for i := int64(0); i < globalParams.NumberOfBuckets/10+1; i++ {
		var wg sync.WaitGroup
		for j := int64(0); j < 10; j++ {
			nextChunkIndex := i*10 + j
			if nextChunkIndex >= globalParams.NumberOfBuckets {
				continue
			}
			wg.Add(1)
			go func(chunkIndex int64) {
				defer wg.Done()

				var activeVertexIds []int64
				verticesChunk := readGraphChunkCompressed(fmt.Sprintf("-%s-%d", globalParams.GraphName, chunkIndex))
				for index, vertex := range verticesChunk {
					verticesChunk[index] = vertex
					activeVertexIds = append(activeVertexIds, vertex.Id)
				}
				insertGraphChunk(verticesChunk, activeVertexIds)
			}(nextChunkIndex)
		}
		wg.Wait()
		log.Info("Group ", i, " of chunks loaded")
	}
}

func insertGraphChunk(verticesChunk []domain.Vertex, activeVertexIds []int64) {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"tag":          "TEST",
			"pureDuration": fmt.Sprintln(utils.MeasureDurationMs(start), "ms"),
		}).Info("Logging insertGraphChunk duration")
	}(time.Now())

	memoryClient.CreateVertices(verticesChunk)
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
	//log.Info("Decompressed chunk: ", string(decompressedBuf))

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
		}).Info("Logging get chunk from storage duration: ", key)
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

	err = functionClient.InvokeFunction(functionapi.OrchestratorFunction, binaryPayload)

	if err != nil {
		panic(err)
	}
}
