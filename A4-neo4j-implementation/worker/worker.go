package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devLucian93/thesis-go/clients/function/api"
	"github.com/devLucian93/thesis-go/clients/memory"
	"math"
	"strconv"
	"time"

	"github.com/devLucian93/thesis-go/domain"
	"github.com/devLucian93/thesis-go/utils"
	"github.com/sirupsen/logrus"
)

var log = utils.GetLogger()

type worker struct {
	runId            string
	superstep        int64
	modifiedVertices []*domain.Vertex
	haltedVertices   []int64
	numberOfVertices int64
	extraArgs        map[string]interface{}
	memoryClient     memory.Client
	functionClient   functionapi.Client
}

type Worker interface {
	//clients.Client
	functionapi.Client
	SaveVertex(vertex *domain.Vertex)
	SendMessage(recipient int64, message interface{})
	SendMessages(recipients []domain.Edge, messages []interface{})
	SendMessageToAllEdges(recipients []domain.Edge, message interface{})
	VoteToHalt(vertex *domain.Vertex)
	HaltWorker()
	//TODO redesign interface. Halt methods and other should not be accessible to the user
	GetModifiedVertices() []*domain.Vertex
	GetRunId() string
	GetSuperstep() int64
	GetNumberOfVertices() int64
	GetExtraArgs() map[string]interface{}
	InvokeOrchestratorFunction(payload string)
	GetArgumentInt(arg interface{}) int64
	GetArgumentFloat(arg interface{}) float64
	ToBool(arg interface{}) bool
	GetFloat(aggregatorKey string, superstep int64) float64
	AggregateFloat(aggregatorKey string, superstep int64, value float64)
}

type Computer interface {
	Worker
	Compute(vertex *domain.Vertex, messages []interface{})
}

func New(payload domain.WorkerPayload, memoryClient memory.Client, functionClient functionapi.Client) (Computer, error) {
	worker := &worker{
		runId:            payload.RunId,
		superstep:        payload.Superstep,
		numberOfVertices: payload.NumberOfVertices,
		memoryClient:     memoryClient,
		functionClient:   functionClient,
		extraArgs:        payload.ExtraArgs,
	}

	switch payload.Algorithm {
	case domain.BFS:
		sourceVertexId := worker.GetArgumentInt(payload.ExtraArgs["sourceVertex"])
		return &BFSWorker{Worker: worker, unvisited: math.MaxInt64, sourceVertexId: sourceVertexId}, nil
	case domain.PAGE_RANK:
		dampingFactor := worker.GetArgumentFloat(payload.ExtraArgs["dampingFactor"])
		numberOfIterations := worker.GetArgumentInt(payload.ExtraArgs["numberOfIterations"])
		danglingNodeSumKey := "danglingNodeSum"
		return &PageRankWorker{
			Worker:             worker,
			dampingFactor:      dampingFactor,
			numberOfIterations: numberOfIterations,
			danglingNodeSumKey: danglingNodeSumKey,
		}, nil
	case domain.SINGLE_SOURCE_SHORTEST_PATH:
		sourceVertexId := worker.GetArgumentInt(payload.ExtraArgs["sourceVertex"])
		return &SingleSourceShortestPathWorker{Worker: worker, sourceVertexId: sourceVertexId}, nil
	case domain.CONNECTED_COMPONENTS:
		directed := worker.ToBool(payload.ExtraArgs["directed"])
		return &ConnectedComponentsWorker{Worker: worker, directed: directed}, nil
	case domain.LOCAL_CLUSTERING_COEFFICIENT:
		directed := worker.ToBool(payload.ExtraArgs["directed"])
		return &LCCWorker{Worker: worker, directed: directed}, nil
	case domain.COMMUNITY_DETECTION:
		directed := worker.ToBool(payload.ExtraArgs["directed"])
		numberOfIterations := worker.GetArgumentInt(payload.ExtraArgs["numberOfIterations"])
		var specializedWorker specializedCDLPWorker
		if directed {
			specializedWorker = DirectedCDLPWorker{Worker: worker, bidirectional: true, unidirectional: false}
		} else {
			specializedWorker = UndirectedCDLPWorker{Worker: worker, bidirectional: true, unidirectional: false}
		}
		return &CommunityDetectionWorker{Worker: worker, numberOfIterations: numberOfIterations, specializedCDLPWorker: specializedWorker}, nil
	default:
		return &BFSWorker{Worker: worker, sourceVertexId: 1}, nil
	}

	return nil, errors.New("Unsupported graph algorithm!")
}

func (w *worker) InvokeFunction(functionName string, payload []byte) error {
	return w.functionClient.InvokeFunction(functionName, payload)
}

func (w *worker) SendMessageToAllEdges(recipients []domain.Edge, message interface{}) {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"runId":        w.runId,
			"tag":          "SEND_MESSAGE_ALL_EDGES",
			"superstep":    w.superstep,
			"pureDuration": utils.MeasureDurationMs(start),
		}).Info("Logging send message to all edges duration")
	}(time.Now())
	w.memoryClient.PutMessageForAllEdges(recipients, message, w.superstep)
}

func (w *worker) SendMessages(recipients []domain.Edge, messages []interface{}) {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"runId":        w.runId,
			"tag":          "SEND_MESSAGES",
			"superstep":    w.superstep,
			"pureDuration": utils.MeasureDurationMs(start),
		}).Info("Logging send messages duration")
	}(time.Now())
	w.memoryClient.PutMessages(recipients, messages, w.superstep)
}

func (w *worker) SendMessage(recipient int64, message interface{}) {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"runId":        w.runId,
			"tag":          "SEND_SINGLE_MESSAGE",
			"superstep":    w.superstep,
			"pureDuration": utils.MeasureDurationMs(start),
		}).Info("Logging send single message duration")
	}(time.Now())
	w.memoryClient.PutMessage(recipient, message, w.superstep)
}

func (w *worker) SaveVertex(vertex *domain.Vertex) {
	//log.Printf("Vertex %d voting to halt", vertex.Id)
	w.modifiedVertices = append(w.modifiedVertices, vertex)
}

func (w *worker) GetModifiedVertices() []*domain.Vertex {
	return w.modifiedVertices
}

func (w *worker) VoteToHalt(vertex *domain.Vertex) {
	//log.Printf("Vertex %d voting to halt", vertex.Id)
	w.haltedVertices = append(w.haltedVertices, vertex.Id)
}

func (w *worker) HaltWorker() {
	activeWorkersCount := w.memoryClient.DecrementActiveWorkersCount()
	log.WithFields(logrus.Fields{
		"runId":     w.runId,
		"tag":       "ACTIVE_WORKERS_W",
		"superstep": w.superstep,
		"workers":   activeWorkersCount,
	}).Info("Logging number of active workers")

	if activeWorkersCount == 0 {
		w.InvokeOrchestratorFunction("SUPERSTEP FINISHED")
	}
}

func (w *worker) GetRunId() string {
	return w.runId
}

func (w *worker) GetSuperstep() int64 {
	return w.superstep
}

func (w *worker) GetNumberOfVertices() int64 {
	return w.numberOfVertices
}

func (w *worker) GetExtraArgs() map[string]interface{} {
	return w.extraArgs
}

func (w *worker) InvokeOrchestratorFunction(payload string) {
	log.Println("Invoking orchestrator function")
	binaryPayload, err := json.Marshal(domain.OrchestratorPayload{payload})

	if err != nil {
		panic(err)
	}

	err = w.InvokeFunction(functionapi.OrchestratorFunction, binaryPayload)

	if err != nil {
		panic(err)
	}
}

func (worker *worker) GetArgumentInt(arg interface{}) int64 {
	switch v := arg.(type) {
	case float64:
		//When unmarshalling into an interface{}, Go decodes JSON numbers as float64 by default
		//If the value was math.MaxInt64 the decoded value will not correspond

		if uint64(v) > uint64(math.MaxInt64) {
			v := math.MaxInt64
			return int64(v)
		}

		return int64(v)
	case int:
		return int64(v)
	case int64:
		return int64(v)
	case string:
		converted, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		return converted
	}
	panic(errors.New(fmt.Sprintf("Could not convert arg %v to int64", arg)))
}

func (worker *worker) GetArgumentFloat(arg interface{}) float64 {
	switch v := arg.(type) {
	case float64:
		return v
	case int64:
		return float64(v)
	case int:
		return float64(v)
	case string:
		converted, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(err)
		}
		return converted
	}
	panic(errors.New(fmt.Sprintf("Could not convert arg of type %T and value %v to float64", arg, arg)))
}

func (worker *worker) ToBool(arg interface{}) bool {
	switch v := arg.(type) {
	case bool:
		return v
	case string:
		converted, err := strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
		return converted
	}
	panic(errors.New(fmt.Sprintf("Could not convert arg of type %T and value %v to bool", arg, arg)))
}

func (worker *worker) GetFloat(aggregatorKey string, superstep int64) float64 {
	return worker.memoryClient.GetFloatMcl(aggregatorKey, superstep)
}
func (worker *worker) AggregateFloat(aggregatorKey string, superstep int64, value float64) {
	worker.memoryClient.AggregateFloatMcl(aggregatorKey, superstep, value)
}
