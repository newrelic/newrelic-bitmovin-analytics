package main

import (
	"context"

	aws_lambda "github.com/aws/aws-lambda-go/lambda"

	"newrelic/multienv/integration"
	"newrelic/multienv/pkg/config"
	nri_lambda "newrelic/multienv/pkg/env/lambda"

	log "github.com/sirupsen/logrus"
)

var pipeConf config.PipelineConfig
var initConf config.RecvConfig
var initErr error

func init() {
	pipeConf = nri_lambda.LoadConfig()
	initConf, initErr = integration.InitRecv(&pipeConf)
}

// TODO: maybe we can get the scheduler rule from the context
// https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html

func HandleRequest(ctx context.Context, event any) (map[string]any, error) {
	if initErr != nil {
		log.Error("Error initializing = ", initErr)
		return nil, initErr
	}

	log.Print("Event received: ", event)

	data, reqErr := initConf.Connector.Request()
	if reqErr.Err != nil {
		log.Error("Http Get error = ", reqErr.Err.Error())
		return nil, reqErr.Err
	}

	log.Print("Data received: ", data)

	var deserBuffer map[string]any
	desErr := initConf.Deser(data, &deserBuffer)
	if desErr != nil {
		log.Error("Error deserializing data = ", desErr)
		return nil, desErr
	}

	log.Print("Data deserialized: ", deserBuffer)

	return deserBuffer, nil
}

func main() {
	aws_lambda.Start(HandleRequest)
}
