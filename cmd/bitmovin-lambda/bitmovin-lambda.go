package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/newrelic/newrelic-bitmovin-analytics/internal/bitmovin"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/log"
)

var (
	/* Args below are populated via ldflags at build time */
	gIntegrationID      = "com.newrelic.labs.newrelic-bitmovin-analytics"
	gIntegrationName    = "New Relic Bitmovin Analytics Integration Lambda"
	gIntegrationVersion = "2.0.0"
	gGitCommit          = ""
	gBuildDate          = ""
	gBuildInfo			= integration.BuildInfo{
		Id:        gIntegrationID,
		Name:      gIntegrationName,
		Version:   gIntegrationVersion,
		GitCommit: gGitCommit,
		BuildDate: gBuildDate,
	}
)

type BitmovinLambdaResult struct {
  Success           bool
  Message           error
}

func HandleRequest(ctx context.Context, event any) (
  BitmovinLambdaResult,
  error,
) {
	// Create the integration with options
	i, err := integration.NewStandaloneIntegration(
		&gBuildInfo,
		gBuildInfo.Name,
		integration.WithInterval(60),
		integration.WithLicenseKey(),
	)
	if err != nil {
		log.Errorf("failed to create integration: %v", err)
		return BitmovinLambdaResult{ false, err }, err
	}

	err = bitmovin.InitPipelines(i)
	if err != nil {
		log.Errorf("failed to initialize pipelines: %v", err)
		return BitmovinLambdaResult{ false, err }, err
	}

	// Run the integration
	defer i.Shutdown(ctx)

 	err = i.Run(ctx)
	if err != nil {
		log.Errorf("integration failed: %v", err)
		return BitmovinLambdaResult{ false, err }, err
	}

  return BitmovinLambdaResult{true, nil}, nil
}

func main() {
  lambda.Start(HandleRequest)
}
