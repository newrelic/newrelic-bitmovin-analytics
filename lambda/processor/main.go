package main

import (
	"context"
	"encoding/json"

	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"

	"newrelic/multienv/integration"
	"newrelic/multienv/pkg/config"
	nri_lambda "newrelic/multienv/pkg/env/lambda"

	"github.com/mitchellh/mapstructure"
)

var pipeConf config.PipelineConfig
var initConf config.ProcConfig
var initErr error

func init() {
	pipeConf = nri_lambda.LoadConfig()
	initConf, initErr = integration.InitProc(&pipeConf)
}

func HandleRequest(ctx context.Context, event map[string]any) (any, error) {
	if initErr != nil {
		log.Error("Error initializing = ", initErr)
		return nil, initErr
	}

	log.Print("Processor event received = ", event)

	model := initConf.Model
	responsePayload, ok := event["responsePayload"]
	if ok {
		err := mapstructure.Decode(responsePayload, &model)
		if err == nil {
			data := integration.Proc(model)
			log.Print("Sending to SQS processed data = ", data)
			return json.Marshal(data)
		} else {
			// We don't return an error because we don't want AWS to retry, just ignore this event.
			log.Error("Error mapping data = ", err)
		}
	} else {
		log.Error("responsePayload not present")
	}

	return nil, nil
}

func main() {
	aws_lambda.Start(HandleRequest)
}
