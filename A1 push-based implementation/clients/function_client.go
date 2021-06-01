package clients

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type FunctionClient interface {
	InvokeFunction(functionName string, payload []byte) error
}

func GetFunctionClient(client FunctionClientType) (FunctionClient, error) {
	switch client {
	case LAMBDA:
		sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-2")}))
		lambda := lambda.New(sess)
		return &LambdaClient{lambda}, nil
	}

	return nil, errors.New("Unsupported function client!")
}

type LambdaClient struct {
	*lambda.Lambda
}

func (client *LambdaClient) InvokeFunction(functionName string, payload []byte) error {
	_, err := client.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: aws.String("Event"),
		Payload:        payload,
	})
	return err
}
