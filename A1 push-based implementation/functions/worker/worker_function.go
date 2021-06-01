package main

import (
	"context"
	"os"
	"time"

	lambdaStarter "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/devLucian93/thesis-go/clients"
	"github.com/devLucian93/thesis-go/domain"
	"github.com/devLucian93/thesis-go/utils"
	"github.com/devLucian93/thesis-go/worker"
	"github.com/sirupsen/logrus"
)

var functionClient clients.FunctionClient
var memoryClient clients.MemoryClient
var haltedVerticesCounter int64
var coldStart bool
var log *logrus.Logger

func init() {
	coldStart = true
	log = utils.GetLogger()
	mcl, err := clients.GetMemoryClient(clients.REDIS)
	memoryClient = mcl
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

func WorkerFunction(ctx context.Context, workerPayload domain.WorkerPayload) {
	if coldStart {
		lc, _ := lambdacontext.FromContext(ctx)
		start := time.Now()
		coldStart = false
		defer func(start time.Time) {
			log.WithFields(logrus.Fields{
				"runId":        workerPayload.RunId,
				"requestId":    lc.AwsRequestID,
				"tag":          "COLD_START",
				"superstep":    workerPayload.Superstep,
				"pureDuration": utils.MeasureDurationMs(start),
			}).Info("Logging pure worker duration")
		}(start)
	}

	worker, err := worker.New(workerPayload, memoryClient, functionClient)
	if err != nil {
		panic(err)
	}

	vertices := getVertices(&workerPayload)

	if workerPayload.Superstep == 0 {
		for i, _ := range vertices {
			noMsg := make([]interface{}, 0)
			worker.Compute(&vertices[i], noMsg)
		}
	} else {
		//log.Println("Worker in superstep: ", SUPERSTEP)
		for i, _ := range vertices {
			//Superstep -1 because we need to get the messages from the previous superstep; TODO: fix Unvisited = -1 ugliness
			worker.Compute(&vertices[i], getMessages(vertices[i].Id, &workerPayload))
		}

	}

	//batch saving vertices
	modifiedVertices := worker.GetModifiedVertices()
	vertices = make([]domain.Vertex, len(worker.GetModifiedVertices()))
	for i, vertex := range modifiedVertices {
		vertices[i] = *vertex
	}

	if len(vertices) > 0 {
		memoryClient.PutVertices(vertices)
	}
	worker.HaltVertices()
	worker.HaltWorker()
}

func main() {
	lambdaStarter.Start(WorkerFunction)
}

func getMessages(vertexId int64, workerPayload *domain.WorkerPayload) []interface{} {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"runId":        workerPayload.RunId,
			"tag":          "GET_MESSAGES",
			"superstep":    workerPayload.Superstep,
			"pureDuration": utils.MeasureDurationMs(start),
		}).Info("Logging get messages duration")
	}(time.Now())
	return memoryClient.GetMessages(vertexId, workerPayload.Superstep-1)
}

func getVertices(workerPayload *domain.WorkerPayload) []domain.Vertex {
	defer func(start time.Time) {
		log.WithFields(logrus.Fields{
			"runId":        workerPayload.RunId,
			"tag":          "GET_VERTICES",
			"superstep":    workerPayload.Superstep,
			"pureDuration": utils.MeasureDurationMs(start),
		}).Info("Logging get vertices duration")
	}(time.Now())
	return memoryClient.GetVertices(workerPayload.VertexIds)
}
