package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/devLucian93/thesis-go/domain"
	"github.com/devLucian93/thesis-go/utils"
	"github.com/sirupsen/logrus"

	lambdaStarter "github.com/aws/aws-lambda-go/lambda"

	"github.com/devLucian93/thesis-go/clients"
)

var storageClient clients.StorageClient
var memoryClient clients.MemoryClient
var functionClient clients.FunctionClient
var queueClient clients.QueueClient
var functionName = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
var log *logrus.Logger

const INITIAL_SUPERSTEP = 0

func init() {
	log = utils.GetLogger()
	scl, err := clients.GetStorageClient(clients.S3)
	storageClient = scl
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	mcl, err := clients.GetMemoryClient(clients.REDIS)
	memoryClient = mcl
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	fcl, err := clients.GetFunctionClient(clients.LAMBDA)
	functionClient = fcl
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	qcl, err := clients.GetQueueClient(clients.QUEUE_REDIS)
	queueClient = qcl
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}

func OrchestratorFunction(payload domain.OrchestratorPayload) (msg string, e error) {

	globalParams, err := memoryClient.GetGlobalParams()
	if err != nil {
		panic(err)
	}
	if payload.Message == domain.ORCHESTRATOR_INIT {
		startProcessing(globalParams)
	} else {
		newActiveVertices := memoryClient.GetMessageRecipients(globalParams.Superstep)
		//fmt.Println("New active vertices are", newActiveVertices)
		if len(newActiveVertices) > 0 {
			memoryClient.AddActiveVertices(newActiveVertices)
		}
		activeVerticesCount := memoryClient.GetActiveVerticesCount()
		log.Println("Active vertices count ", activeVerticesCount)

		if activeVerticesCount > 0 {
			startNewSuperstep(globalParams)
		} else {
			doFinishOperations(globalParams)
		}
	}
	return "Finished executing orchestrator", nil
}

func main() {
	lambdaStarter.Start(OrchestratorFunction)
}

func startProcessing(globalParams *domain.GlobalParams) {
	dataIngestionDuration := ((time.Now().UnixNano() - globalParams.DataIngestionDuration) / (1e9))
	log.Println("Data ingestion duration ", dataIngestionDuration, "s\nLaunching initial worker wave")
	globalParams.DataIngestionDuration = dataIngestionDuration
	globalParams.ExecutionDuration = time.Now().UnixNano()
	memoryClient.PutGlobalParams(globalParams)

	fanoutWork([]int64{}, globalParams)

}

func fanoutWork(recipients []int64, globalParams *domain.GlobalParams) {
	numberOfRecipients := int64(len(recipients))
	var workersToStart int64
	if globalParams.Superstep == 0 {
		if workersToStart = int64(math.Ceil(float64(globalParams.NumberOfVertices) / float64(globalParams.ChunkSize))); workersToStart > globalParams.MaxWorkers {
			workersToStart = globalParams.MaxWorkers
		}
		memoryClient.SetActiveWorkersCount(workersToStart)
		invokeWorkers(workersToStart, globalParams)

	} else {
		if workersToStart = int64(math.Ceil(float64(numberOfRecipients) / float64(globalParams.ChunkSize))); workersToStart > globalParams.MaxWorkers {
			workersToStart = globalParams.MaxWorkers
		}
		memoryClient.SetActiveWorkersCount(workersToStart)

		PushWorkerTasks(globalParams, recipients)
		SortWorkerTasks(globalParams)

		invokeWorkers(workersToStart, globalParams)
	}
	log.WithFields(logrus.Fields{
		"runId":     globalParams.RunId,
		"tag":       "ACTIVE_WORKERS",
		"superstep": globalParams.Superstep,
		"workers":   workersToStart,
	}).Info("Logging number of active workers")

}

func PushWorkerTasks(globalParams *domain.GlobalParams, tasks []int64) {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"runId":        globalParams.RunId,
			"tag":          "PUSH_TASKS",
			"superstep":    globalParams.Superstep,
			"pureDuration": utils.MeasureDurationMs(start),
		}).Info("Pushing worker tasks")
	}(time.Now())

	if err := queueClient.PushWorkerTasks(tasks); err != nil {
		panic(err)
	}
}

func SortWorkerTasks(globalParams *domain.GlobalParams) {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"runId":        globalParams.RunId,
			"tag":          "SORT_TASKS",
			"superstep":    globalParams.Superstep,
			"pureDuration": utils.MeasureDurationMs(start),
		}).Info("Sorting worker tasks")
	}(time.Now())
	if err := queueClient.SortTasks(); err != nil {
		panic(err)
	}
}

func invokeWorkers(workersToStart int64, globalParams *domain.GlobalParams) {
	for i := int64(0); i < workersToStart; i++ {
		workerPayload := &domain.WorkerPayload{
			Superstep:        globalParams.Superstep,
			Algorithm:        globalParams.Algorithm,
			ChunkSize:        globalParams.ChunkSize,
			ExtraArgs:        globalParams.ExtraArgs,
			NumberOfVertices: globalParams.NumberOfVertices,
			RunId:            globalParams.RunId,
		}
		invokeWorkerFunction(workerPayload)
	}
}

func invokeWorkerFunction(workerPayload *domain.WorkerPayload) {
	//log.Println("Invoking worker function")
	binaryPayload, err := workerPayload.MarshalBinary()
	if err != nil {
		panic(err)
	}
	err = functionClient.InvokeFunction(os.Getenv("WORKER_FUNCTION_NAME"), binaryPayload)
	if err != nil {
		panic(err)
	}

}

func startNewSuperstep(globalParams *domain.GlobalParams) {
	activeVertices := memoryClient.GetActiveVertices()
	globalParams.Superstep++
	memoryClient.PutGlobalParams(globalParams)
	log.Printf("Start of new superstep: %d. Active vertices: %d\n", globalParams.Superstep, len(activeVertices))

	fanoutWork(activeVertices, globalParams)

}

func doFinishOperations(globalParams *domain.GlobalParams) {
	chunkSize := 1000
	executionDuration := time.Now().UnixNano() - globalParams.ExecutionDuration
	globalParams.ExecutionDuration = executionDuration
	globalParams.Finished = true
	memoryClient.PutGlobalParams(globalParams)
	globalParamsJsonBytes, _ := json.Marshal(globalParams)
	storageClient.Put(fmt.Sprintf("-%s-%s-%s-%s", globalParams.RunId,
		globalParams.GraphName, globalParams.Algorithm, "results-metadata"), string(globalParamsJsonBytes[:]))

	vertexIds := memoryClient.GetAllVertexIds()

	var allIdValuePairs []domain.IdValuePair
	var vertexChunk []domain.Vertex
	for i := 0; i < len(vertexIds); i += chunkSize {
		if i+chunkSize < len(vertexIds) {
			vertexChunk = memoryClient.GetVertices(vertexIds[i : i+chunkSize])
		} else {
			vertexChunk = memoryClient.GetVertices(vertexIds[i:])
		}

		for _, vertex := range vertexChunk {
			allIdValuePairs = append(allIdValuePairs, domain.IdValuePair{Id: vertex.Id, Value: vertex.Value})
		}
	}
	sort.Slice(allIdValuePairs, func(i, j int) bool {
		return allIdValuePairs[i].Id < allIdValuePairs[j].Id
	})

	finalResults := make([]string, globalParams.NumberOfVertices, globalParams.NumberOfVertices)
	for i, idValuePair := range allIdValuePairs {
		if math.IsInf(idValuePair.Value.(float64), 0) {
			idValuePair.Value = "infinity"
		} else if uint64(idValuePair.Value.(float64)) > uint64(math.MaxInt64) {
			//Because JSON treats all numbers as floats and using math.MaxInt64 leads to a float value 1 greater
			idValuePair.Value = math.MaxInt64
		}

		finalResults[i] = fmt.Sprintf("%d %v", idValuePair.Id, idValuePair.Value)
	}

	joinedResult := strings.Join(finalResults, "\n") + "\n"
	storageClient.Put(fmt.Sprintf("-%s-%s-%s-%s", globalParams.RunId, globalParams.GraphName, globalParams.Algorithm, "results"), joinedResult)
	log.Println("Uploaded results to remote storage")
	log.WithFields(logrus.Fields{
		"runId": globalParams.RunId,
		"tag":   "FINISHED",
	}).Info("Finished processing graph! Global params", globalParams)
}
