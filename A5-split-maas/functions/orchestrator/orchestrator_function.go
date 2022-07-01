package main

import (
	lambdaStarter "github.com/aws/aws-lambda-go/lambda"
	"github.com/devLucian93/thesis-go/clients/function"
	functionapi "github.com/devLucian93/thesis-go/clients/function/api"
	"github.com/devLucian93/thesis-go/domain"
	orchestrator "github.com/devLucian93/thesis-go/functions/orchestrator/impl"
	"github.com/devLucian93/thesis-go/utils"
	"github.com/sirupsen/logrus"
	"os"
)

var functionClient functionapi.Client
var log *logrus.Logger

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
}

func OrchestratorFunction(payload domain.OrchestratorPayload) (msg string, e error) {
	o := orchestrator.Facade{
		FunctionClient: functionClient,
	}
	return o.OrchestratorFunction(payload)
}

func main() {
	lambdaStarter.Start(OrchestratorFunction)
}
