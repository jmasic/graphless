package workerimpl

import (
	"github.com/devLucian93/thesis-go/clients/function/api"
	"github.com/devLucian93/thesis-go/clients/memory"
	"github.com/devLucian93/thesis-go/clients/message"
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
	messageClient  message.Client
}

func (w *Facade) initInfraClients(workerPayload domain.WorkerPayload) {
	log = utils.GetLogger()

	memoryClientType := memory.ResolveClientType(workerPayload.MemoryClientConfig.ClientType)
	mcl, err := memory.GetMemoryClient(memoryClientType, workerPayload.MemoryClientConfig.DbConfig)
	w.memoryClient = mcl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}

	messageClientType := message.ResolveClientType(workerPayload.MessageClientConfig.ClientType)
	msgCl, err := message.GetMessageClient(messageClientType, workerPayload.MessageClientConfig.DbConfig)
	w.messageClient = msgCl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
}

func (w *Facade) WorkerFunction(workerPayload domain.WorkerPayload) {
	w.initInfraClients(workerPayload)

	aWorker, err := worker.New(workerPayload, w.memoryClient, w.messageClient, w.FunctionClient)
	if err != nil {
		panic(err)
	}

	vertexChannel := w.getVertices(&workerPayload)

	for verticesChunk := range vertexChannel {
		if verticesChunk == nil {
			break
		}
		for _, vertex := range verticesChunk {
			messages := w.getMessages(vertex.Id, &workerPayload)
			if workerPayload.Superstep == 0 {
				messages = make([]interface{}, 0) // no messages at the first step
			}
			aWorker.Compute(vertex, messages)
		}
	}

	// batch saving vertices
	modifiedVertices := aWorker.GetModifiedVertices()
	vertices := make([]domain.Vertex, len(modifiedVertices))
	for i, vertex := range modifiedVertices {
		vertices[i] = *vertex
	}
	if len(vertices) > 0 {
		w.memoryClient.SaveVertices(vertices)
	}

	aWorker.HaltWorker(workerPayload)
}

func (w *Facade) toIds(vertices []domain.Vertex) []int64 {
	ids := make([]int64, len(vertices))
	for i, vertex := range vertices {
		ids[i] = vertex.Id
	}
	return ids
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
	return w.messageClient.GetMessages(vertexId, workerPayload.Superstep-1)
}

func (w *Facade) getVertices(workerPayload *domain.WorkerPayload) <-chan []domain.Vertex {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"runId":        workerPayload.RunId,
			"tag":          "GET_VERTICES",
			"superstep":    workerPayload.Superstep,
			"pureDuration": utils.MeasureDurationMs(start),
		}).Info("Logging get vertices duration")
	}(time.Now())
	//log.Info("Getting these vertices for worker: ", workerPayload.VertexIds)
	vertexIds := workerPayload.VertexIds
	return w.memoryClient.GetVertices(vertexIds)
}
