package main

import (
	"newrelic/multienv/integration"
	"newrelic/multienv/pkg/env/infra"
	"newrelic/multienv/pkg/export"
	"os"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func main() {
	pipeConf, err := infra.LoadConfig()
	if err != nil {
		log.Error("Error loading config: ", err)
		os.Exit(1)
	}

	recvConfig, err := integration.InitRecv(&pipeConf)
	if err != nil {
		log.Error("Error initializing receiver: ", err)
		os.Exit(2)
	}

	procConfig, err := integration.InitProc(&pipeConf)
	if err != nil {
		log.Error("Error initializing processor: ", err)
		os.Exit(3)
	}

	res, reqErr := recvConfig.Connector.Request()
	if reqErr.Err != nil {
		log.Error("Error connecting: ", reqErr.Err)
		os.Exit(4)
	}

	var deserBuffer map[string]any
	errDes := recvConfig.Deser(res, &deserBuffer)
	if errDes != nil {
		log.Error("Error deserializing: ", errDes)
		os.Exit(5)
	}

	model := procConfig.Model
	errMap := mapstructure.Decode(deserBuffer, &model)
	if errMap != nil {
		log.Error("Error mapping: ", errMap)
		os.Exit(6)
	}

	meltData := integration.Proc(model)

	exporter := export.SelectExporter(pipeConf.Exporter)

	errExp := exporter(pipeConf, meltData)
	if errExp != nil {
		log.Error("Error exporting: ", errExp)
		os.Exit(7)
	}
}
