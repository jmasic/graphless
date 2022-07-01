package workerimpl

import (
	"github.com/devLucian93/thesis-go/clients/function/api"
	"github.com/devLucian93/thesis-go/clients/memory"
	"github.com/devLucian93/thesis-go/worker"
	"os"
	"time"

	"github.com/devLucian93/thesis-go/domain"
	"github.com/devLucian93/thesis-go/utils"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

type Facade struct {
	FunctionClient functionapi.Client
	memoryClient   memory.Client
}

func (w *Facade) initInfraClients() {
	local := utils.IsLocal()
	log = utils.GetLogger()
	memoryClientType := (map[bool]memory.ClientType{true: memory.Neo4j, false: memory.Redis})[local]
	mcl, err := memory.GetMemoryClient(memoryClientType)
	w.memoryClient = mcl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
}

func (w *Facade) WorkerFunction(workerPayload domain.WorkerPayload) {
	w.initInfraClients()

	aWorker, err := worker.New(workerPayload, w.memoryClient, w.FunctionClient)
	if err != nil {
		panic(err)
	}

	vertices := w.getVertices(&workerPayload)

	if workerPayload.Superstep == 0 {
		for i, _ := range vertices {
			noMsg := make([]interface{}, 0)
			aWorker.Compute(&vertices[i], noMsg)
		}
	} else {
		log.Println("Worker in superstep:", workerPayload.Superstep, "with", len(vertices), "vertices")
		for i, _ := range vertices {
			// Superstep -1 because we need to get the messages from the previous superstep; TODO: fix Unvisited = -1 ugliness
			aWorker.Compute(&vertices[i], w.getMessages(vertices[i].Id, &workerPayload))
		}
	}

	// batch saving vertices
	modifiedVertices := aWorker.GetModifiedVertices()
	vertices = make([]domain.Vertex, len(aWorker.GetModifiedVertices()))
	for i, vertex := range modifiedVertices {
		vertices[i] = *vertex
	}

	if len(vertices) > 0 {
		w.memoryClient.PutVertices(vertices)
	}
	aWorker.HaltVertices()
	aWorker.HaltWorker()
}

func (w *Facade) getMessages(vertexId int64, workerPayload *domain.WorkerPayload) []interface{} {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"runId":        workerPayload.RunId,
			"tag":          "GET_MESSAGES",
			"superstep":    workerPayload.Superstep,
			"pureDuration": utils.MeasureDurationMs(start),
		}).Info("Logging get messages duration")
	}(time.Now())
	return w.memoryClient.GetMessages(vertexId, workerPayload.Superstep-1)
}

func (w *Facade) getVertices(workerPayload *domain.WorkerPayload) []domain.Vertex {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"runId":        workerPayload.RunId,
			"tag":          "GET_VERTICES",
			"superstep":    workerPayload.Superstep,
			"pureDuration": utils.MeasureDurationMs(start),
		}).Info("Logging get vertices duration")
	}(time.Now())
	//log.Info("Getting these vertices for worker: ", workerPayload.VertexIds)
	return w.memoryClient.GetVertices(workerPayload.VertexIds)
}
