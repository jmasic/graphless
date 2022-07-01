package main

import (
	"context"
	lambdaStarter "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/devLucian93/thesis-go/clients/function"
	"github.com/devLucian93/thesis-go/clients/function/api"
	"github.com/devLucian93/thesis-go/domain"
	workerimpl "github.com/devLucian93/thesis-go/functions/worker/impl"
	"github.com/devLucian93/thesis-go/utils"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var functionClient functionapi.Client
var coldStart bool
var log *logrus.Logger

func init() {
	coldStart = true
	log = utils.GetLogger()

	fcl, err := function.GetFunctionClient(functionapi.AwsLambda)
	functionClient = fcl
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
}

func WorkerFunction(ctx context.Context, payload domain.WorkerPayload) {
	if coldStart {
		lc, _ := lambdacontext.FromContext(ctx)
		start := time.Now()
		coldStart = false
		defer func(start time.Time) {
			log.WithFields(logrus.Fields{
				"runId":        payload.RunId,
				"requestId":    lc.AwsRequestID,
				"tag":          "COLD_START",
				"superstep":    payload.Superstep,
				"pureDuration": utils.MeasureDurationMs(start),
			}).Info("Logging pure worker duration")
		}(start)
	}

	w := workerimpl.Facade{
		FunctionClient: functionClient,
	}
	w.WorkerFunction(payload)
}

func main() {
	lambdaStarter.Start(WorkerFunction)
}
