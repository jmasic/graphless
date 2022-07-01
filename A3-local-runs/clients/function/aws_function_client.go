package function

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	log "github.com/sirupsen/logrus"
	"os"
)

type awsLambdaClient struct {
	*lambda.Lambda
}

func newAwsLambdaClient() (*awsLambdaClient, error) {
	config := aws.NewConfig().
		WithRegion("us-east-2")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *config,
		SharedConfigState: session.SharedConfigEnable,
	}))
	log.Info("Session: ", sess)
	client := lambda.New(sess)
	return &awsLambdaClient{client}, nil
}

func (client *awsLambdaClient) InvokeFunction(functionName string, payload []byte) error {
	awsFunctionName := os.Getenv(functionName)
	_, err := client.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String(awsFunctionName),
		InvocationType: aws.String("Event"),
		Payload:        payload,
	})
	return err
}
