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
		//log.Println("New active vertices are", newActiveVertices)
		if len(newActiveVertices) > 0 {
			o.memoryClient.AddActiveVertices(newActiveVertices)
		}
		activeVerticesCount := o.memoryClient.GetActiveVerticesCount()
		log.Println("Active vertices count:", activeVerticesCount)

		if activeVerticesCount > 0 {
			o.startNewSuperstep(globalParams)
		} else {
			o.doFinishOperations(globalParams)
		}
	}
	return "Finished executing orchestrator", nil
}

func (o *Facade) startProcessing(globalParams *domain.GlobalParams) {
	dataIngestionDuration := ((time.Now().UnixNano() - globalParams.DataIngestionDuration) / (1e9))
	log.Println("Data ingestion duration ", dataIngestionDuration, "s\nLaunching initial worker wave")
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
		"tag":       "ACTIVE_WORKERS",
		"superstep": globalParams.Superstep,
		"workers":   activeWorkers,
	}).Info("Logging number of active workers")

	o.memoryClient.SetActiveWorkersCount(activeWorkers)
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
		o.invokeWorkerFunction(workerPayload)
	}
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

func (o *Facade) startNewSuperstep(globalParams *domain.GlobalParams) {
	activeVertices := o.memoryClient.GetActiveVertices()
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

	var allIdValuePairs []domain.IdValuePair
	var vertexChunk []domain.Vertex
	for i := 0; i < len(vertexIds); i += chunkSize {
		if i+chunkSize < len(vertexIds) {
			vertexChunk = o.memoryClient.GetVertices(vertexIds[i : i+chunkSize])
		} else {
			vertexChunk = o.memoryClient.GetVertices(vertexIds[i:])
		}

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
	o.storageClient.Put(fmt.Sprintf("-%s-%s-%s-%s", globalParams.RunId, globalParams.GraphName, globalParams.Algorithm, "results"), joinedResult)
	log.Println("Uploaded results to remote storage")
	log.WithFields(logrus.Fields{
		"runId": globalParams.RunId,
		"tag":   "FINISHED",
	}).Info("Finished processing graph! Global params", globalParams)
}
