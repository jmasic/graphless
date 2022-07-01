package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"github.com/devLucian93/thesis-go/clients/function"
	functionapi "github.com/devLucian93/thesis-go/clients/function/api"
	"github.com/devLucian93/thesis-go/clients/memory"
	"github.com/devLucian93/thesis-go/clients/message"
	"github.com/devLucian93/thesis-go/clients/storage"
	"io"
	"io/ioutil"
	"math"
	"os"
	"time"

	lambdaStarter "github.com/aws/aws-lambda-go/lambda"
	"github.com/devLucian93/thesis-go/domain"
	"github.com/devLucian93/thesis-go/utils"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
)

var storageClient storage.Client
var memoryClient memory.Client
var messageClient message.Client
var functionClient functionapi.Client
var log *logrus.Logger

var numberOfLoadingWorkers int64

func init() {
	local := utils.IsLocal()
	log = utils.GetLogger()

	functionClientType := (map[bool]functionapi.ClientType{true: functionapi.GoFunction, false: functionapi.AwsLambda})[local]
	fcl, err := function.GetFunctionClient(functionClientType)
	functionClient = fcl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}

	numberOfLoadingWorkers = (map[bool]int64{true: 15, false: 10})[local]
}

func initInfraClients(params domain.StartParams) {
	memoryClientType := memory.ResolveClientType(params.MemoryClientConfig.ClientType)
	mcl, err := memory.GetMemoryClient(memoryClientType, params.MemoryClientConfig.DbConfig)
	memoryClient = mcl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}

	messageClientType := message.ResolveClientType(params.MessageClientConfig.ClientType)
	msgCl, err := message.GetMessageClient(messageClientType, params.MessageClientConfig.DbConfig)
	messageClient = msgCl
	if err != nil {
		log.Println("Error while instantiating message client: ", err)
		os.Exit(1)
	}

	storageClientType := storage.ResolveClientType(params.StorageClientConfig.ClientType)
	scl, err := storage.GetStorageClient(storageClientType, params.StorageClientConfig.StorageConfig)
	storageClient = scl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}

	memoryClient.Clear()
	messageClient.Clear()
}

func mainFunction(params domain.StartParams) (msg string, e error) {
	log.WithFields(logrus.Fields{
		"runId": params.RunId,
		"tag":   "STARTED",
	}).Info("Started graph processing engine with params: ", params)

	initInfraClients(params)

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

	invokeOrchestratorFunction(domain.ORCHESTRATOR_INIT, params)

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
	chunkChannel := make(chan int64, globalParams.NumberOfBuckets+numberOfLoadingWorkers)
	finishChannel := make(chan bool, numberOfLoadingWorkers)
	defer close(chunkChannel)
	defer close(finishChannel)

	// locally: from 100s to 70-80s (20-30% improvement)
	var i int64
	for i = 0; i < numberOfLoadingWorkers; i++ {
		go readAndInsertGraphChunk(chunkChannel, globalParams, finishChannel)
	}
	for i = 0; i < globalParams.NumberOfBuckets; i++ {
		chunkChannel <- i
	}
	for i = 0; i < numberOfLoadingWorkers; i++ {
		chunkChannel <- -1
	}
	for i = 0; i < numberOfLoadingWorkers; i++ {
		var _ = <-finishChannel
	}

	log.Info("All chunks loaded")

	// NOTE: value-reset implementation
	//vertices := make([]domain.Vertex, 1)
	//memoryClient.CreateVertices(vertices)
	//log.Info("All chunks loaded")
}

func readAndInsertGraphChunk(chunkChannel <-chan int64, globalParams *domain.GlobalParams, finishChannel chan<- bool) {
	for chunk := <-chunkChannel; chunk >= 0; chunk = <-chunkChannel {
		verticesChunk := readGraphChunkCompressed(fmt.Sprintf("-%s-%d", globalParams.GraphName, chunk))
		var activeVertexIds []int64
		for index, vertex := range verticesChunk {
			vertex.Value = math.MaxInt64
			verticesChunk[index] = vertex
			activeVertexIds = append(activeVertexIds, vertex.Id)
		}
		insertGraphChunk(verticesChunk, activeVertexIds)
	}
	finishChannel <- true
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

func invokeOrchestratorFunction(message string, params domain.StartParams) {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"tag":          "TEST",
			"pureDuration": fmt.Sprintln(utils.MeasureDurationMs(start), "ms"),
		}).Info("Logging Invoking orchestrator function duration")
	}(time.Now())
	log.Println("Invoking orchestrator function")
	binaryPayload, err := json.Marshal(domain.OrchestratorPayload{
		message,
		params.MemoryClientConfig,
		params.MessageClientConfig,
		params.StorageClientConfig,
	})

	if err != nil {
		panic(err)
	}

	err = functionClient.InvokeFunction(functionapi.OrchestratorFunction, binaryPayload)

	if err != nil {
		panic(err)
	}
}
