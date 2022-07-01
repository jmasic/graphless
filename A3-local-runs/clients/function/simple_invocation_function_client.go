package function

import (
	"encoding/json"
	functionapi "github.com/devLucian93/thesis-go/clients/function/api"
	"github.com/devLucian93/thesis-go/domain"
	orchestrator "github.com/devLucian93/thesis-go/functions/orchestrator/impl"
	workerimpl "github.com/devLucian93/thesis-go/functions/worker/impl"
)

type simpleInvocationFunctionClient struct {
}

func newSimpleInvocationFunctionClient() (*simpleInvocationFunctionClient, error) {
	return &simpleInvocationFunctionClient{}, nil
}

func (client *simpleInvocationFunctionClient) InvokeFunction(functionName string, payload []byte) error {
	switch functionName {
	case functionapi.OrchestratorFunction:
		var orchestratorPayload domain.OrchestratorPayload
		err := json.Unmarshal(payload, &orchestratorPayload)
		if err != nil {
			panic(err)
		}
		o := orchestrator.Facade{
			FunctionClient: client,
		}
		_, err = o.OrchestratorFunction(orchestratorPayload)
		if err != nil {
			panic(err)
		}
		return nil
	case functionapi.WorkerFunction:
		var workerPayload domain.WorkerPayload
		err := json.Unmarshal(payload, &workerPayload)
		if err != nil {
			panic(err)
		}
		w := workerimpl.Facade{
			FunctionClient: client,
		}
		w.WorkerFunction(workerPayload)
		return nil
	default:
		panic("invoke on '" + functionName + "' not implemented yet")
	}
}
