package function

import (
	"errors"
	functionapi "github.com/devLucian93/thesis-go/clients/function/api"
)

func GetFunctionClient(clientType functionapi.ClientType) (functionapi.Client, error) {
	switch clientType {
	case functionapi.AwsLambda:
		return newAwsLambdaClient()
	case functionapi.GoFunction:
		return newSimpleInvocationFunctionClient()
	}

	return nil, errors.New("Unsupported function client!")
}
