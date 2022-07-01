package orchestrator

import (
	"encoding/json"
	"fmt"
	functionapi "github.com/devLucian93/thesis-go/clients/function/api"
	"github.com/devLucian93/thesis-go/clients/memory"
	"github.com/devLucian93/thesis-go/clients/storage"
	"github.com/devLucian93/thesis-go/utils"
	"math"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/devLucian93/thesis-go/domain"
	"github.com/sirupsen/logrus"
)

type Facade struct {
	FunctionClient functionapi.Client
	storageClient  storage.Client
	memoryClient   memory.Client
}

var log *logrus.Logger

func (o *Facade) initInfraClients() {
	local := utils.IsLocal()
	log = utils.GetLogger()
	storageClientType := (map[bool]storage.ClientType{true: storage.LocalFileSystem, false: storage.S3})[local]
	scl, err := storage.GetStorageClient(storageClientType)
	o.storageClient = scl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
	memoryClientType := (map[bool]memory.ClientType{true: memory.Neo4j, false: memory.Redis})[local]
	mcl, err := memory.GetMemoryClient(memoryClientType)
	o.memoryClient = mcl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
}

func (o *Facade) OrchestratorFunction(payload domain.OrchestratorPayload) (msg string, e error) {
	o.initInfraClients()

	globalParams, err := o.memoryClient.GetGlobalParams()
	if err != nil {
		panic(err)
	}

	if payload.Message == domain.ORCHESTRATOR_INIT {
		o.startProcessing(globalParams)
	} else {
		newActiveVertices := o.memoryClient.GetMessageRecipients(globalParams.Superstep)
		activeVerticesCount := len(newActiveVertices)
		//log.Println("New active vertices are", newActiveVertices)
		log.Println("Active vertices count:", activeVerticesCount)

		if activeVerticesCount > 0 {
			o.startNewSuperstep(globalParams, newActiveVertices)
		} else {
			o.doFinishOperations(globalParams)
		}
	}
	return "Finished executing orchestrator", nil
}

func (o *Facade) startProcessing(globalParams *domain.GlobalParams) {
	dataIngestionDuration := ((time.Now().UnixNano() - globalParams.DataIngestionDuration) / (1e9))
	log.Println("Data ingestion duration", dataIngestionDuration, "s\nLaunching initial worker wave")
	globalParams.DataIngestionDuration = dataIngestionDuration
	globalParams.ExecutionDuration = time.Now().UnixNano()

	err := o.memoryClient.PutGlobalParams(globalParams)
	if err != nil {
		panic("Could not get global params")
	}

	vertexIds := o.memoryClient.GetAllVertexIds()

	o.fanoutWork(vertexIds, globalParams)
}

func (o *Facade) fanoutWork(recipients []int64, globalParams *domain.GlobalParams) {
	numberOfRecipients := int64(len(recipients))
	chunkSize := globalParams.ChunkSize
	activeWorkers := int64(math.Ceil(float64(numberOfRecipients) / float64(chunkSize)))
	if activeWorkers > globalParams.MaxWorkers {
		activeWorkers = globalParams.MaxWorkers
		chunkSize = int64(math.Floor(float64(numberOfRecipients) / float64(activeWorkers)))
	}

	log.WithFields(logrus.Fields{
		"runId":     globalParams.RunId,
		"tag":       "ACTIVE_WORKERS_O",
		"superstep": globalParams.Superstep,
		"workers":   activeWorkers,
	}).Info("Logging number of active workers")

	o.memoryClient.SetActiveWorkersCount(activeWorkers)
	o.invokeWorkers(recipients, globalParams, activeWorkers, chunkSize)
}

func (o *Facade) invokeWorkers(recipients []int64, globalParams *domain.GlobalParams, activeWorkers int64, chunkSize int64) {
	var wg sync.WaitGroup
	for i := int64(0); i < activeWorkers; i++ {
		workerPayload := &domain.WorkerPayload{
			Superstep:        globalParams.Superstep,
			Algorithm:        globalParams.Algorithm,
			ExtraArgs:        globalParams.ExtraArgs,
			NumberOfVertices: globalParams.NumberOfVertices,
			RunId:            globalParams.RunId,
		}
		if (i + 1) != activeWorkers {
			workerPayload.VertexIds = recipients[i*chunkSize : (i+1)*chunkSize]
		} else {
			workerPayload.VertexIds = recipients[i*chunkSize:]
		}
		wg.Add(1)
		go func(payload *domain.WorkerPayload) {
			defer wg.Done()
			o.invokeWorkerFunction(payload)
		}(workerPayload)
	}
	wg.Wait()
}

func (o *Facade) invokeWorkerFunction(workerPayload *domain.WorkerPayload) {
	//log.Println("Invoking worker function")
	binaryPayload, err := workerPayload.MarshalBinary()
	if err != nil {
		panic(err)
	}
	err = o.FunctionClient.InvokeFunction(functionapi.WorkerFunction, binaryPayload)
	if err != nil {
		panic(err)
	}
}

func (o *Facade) startNewSuperstep(globalParams *domain.GlobalParams, activeVertices []int64) {
	globalParams.Superstep++
	err := o.memoryClient.PutGlobalParams(globalParams)
	if err != nil {
		log.Println("Error while putting global params in start superstep: ", err)
	}
	log.Printf("Start of new superstep: %d. Active vertices: %d\n", globalParams.Superstep, len(activeVertices))

	o.fanoutWork(activeVertices, globalParams)
}

func (o *Facade) doFinishOperations(globalParams *domain.GlobalParams) {
	chunkSize := 1000
	executionDuration := time.Now().UnixNano() - globalParams.ExecutionDuration

	globalParams.ExecutionDuration = executionDuration
	globalParams.Finished = true
	err := o.memoryClient.PutGlobalParams(globalParams)
	if err != nil {
		log.Println("Error while putting global params in finishing operations: ", err)
	}
	globalParamsJsonBytes, _ := json.Marshal(globalParams)
	o.storageClient.Put(fmt.Sprintf("-%s-%s-%s-%s", globalParams.RunId,
		globalParams.GraphName, globalParams.Algorithm, "results-metadata.json"), string(globalParamsJsonBytes[:]))

	vertexIds := o.memoryClient.GetAllVertexIds()

	allIdValuePairs := o.readAllVerticesInIdValuePairs(vertexIds, chunkSize)

	finalResults := make([]string, globalParams.NumberOfVertices, globalParams.NumberOfVertices)
	for i, idValuePair := range allIdValuePairs {
		finalResults[i] = fmt.Sprintf("%d %v", idValuePair.Id, prepareResult(idValuePair))
	}

	joinedResult := strings.Join(finalResults, "\n") + "\n"
	err = o.storageClient.Put(fmt.Sprintf("-%s-%s-%s-%s", globalParams.RunId, globalParams.GraphName, globalParams.Algorithm, "results"), joinedResult)
	if err != nil {
		log.Error("Couldn't upload results to storage")
		panic(err)
	}
	log.Println("Uploaded results to remote storage")
	log.WithFields(logrus.Fields{
		"runId": globalParams.RunId,
		"tag":   "FINISHED",
	}).Info("Finished processing graph! Global params", globalParams)
}

func (o *Facade) readAllVerticesInIdValuePairs(vertexIds []int64, chunkSize int) []domain.IdValuePair {
	var allIdValuePairs []domain.IdValuePair

	vertexChannel := o.memoryClient.GetVertices(vertexIds)
	for vertexChunk := range vertexChannel {
		for _, vertex := range vertexChunk {
			allIdValuePairs = append(allIdValuePairs, domain.IdValuePair{Id: vertex.Id, Value: vertex.Value})
		}
	}

	sort.Slice(allIdValuePairs, func(i, j int) bool {
		return allIdValuePairs[i].Id < allIdValuePairs[j].Id
	})
	if len(allIdValuePairs) < 30 {
		log.Info("All pairs: ", allIdValuePairs)
	}
	return allIdValuePairs
}

func prepareResult(idValuePair domain.IdValuePair) interface{} {
	value := idValuePair.Value
	switch v := value.(type) {
	case float64:
		if math.IsInf(value.(float64), 0) {
			value = "infinity"
		} else if uint64(value.(float64)) > uint64(math.MaxInt64) {
			// Because JSON treats all numbers as floats and using math.MaxInt64 leads to a float value 1 greater
			value = math.MaxInt64
		}
		return value
	default:
		return v
	}
}
