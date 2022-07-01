package functionapi

type Client interface {
	InvokeFunction(functionName string, payload []byte) error
}

type ClientType int

const (
	AwsLambda  ClientType = iota
	GoFunction ClientType = iota
)

const (
	OrchestratorFunction string = "ORCHESTRATOR_FUNCTION_NAME"
	WorkerFunction       string = "WORKER_FUNCTION_NAME"
)
