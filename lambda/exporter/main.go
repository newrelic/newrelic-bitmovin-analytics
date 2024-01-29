package main

import (
	"context"
	"encoding/base64"
	"encoding/json"

	aws_events "github.com/aws/aws-lambda-go/events"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"

	"newrelic/multienv/pkg/config"
	nri_lambda "newrelic/multienv/pkg/env/lambda"
	"newrelic/multienv/pkg/export"
	"newrelic/multienv/pkg/model"
)

var pipeConf config.PipelineConfig

func init() {
	pipeConf = nri_lambda.LoadConfig()
}

func HandleRequest(ctx context.Context, event aws_events.SQSEvent) error {

	samples := make([]model.MeltModel, 0)

	for i, record := range event.Records {
		var obj map[string]any
		json.Unmarshal([]byte(record.Body), &obj)
		responsePayload, ok := obj["responsePayload"].(string)
		if ok {
			// Decode Base64 and Unmarshal JSON
			sDec, err := base64.StdEncoding.DecodeString(responsePayload)
			if err != nil {
				log.Error("Error decoding = ", err)
				continue
			}

			responsePayload = string(sDec)

			var model []model.MeltModel
			err = json.Unmarshal([]byte(responsePayload), &model)
			if err != nil {
				log.Error("Error mapping MELT model = ", err)
				continue
			}

			log.Printf("(%d) SQS decoded response payload = %v", i, model)

			samples = append(samples, model...)

		} else {
			log.Error("Couldn't get a string from 'responsePayload'")
		}
	}

	// Export data
	exporter := export.SelectExporter(pipeConf.Exporter)
	err := exporter(pipeConf, samples)
	if err != nil {
		log.Error("Exporter failed = ", err)
		//TODO: handle error condition, refill buffer? Discard data? Retry?
	}

	return nil
}

func main() {
	aws_lambda.Start(HandleRequest)
}
