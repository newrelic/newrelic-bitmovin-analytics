package main

import (
	"context"

	"github.com/newrelic/newrelic-labs-sdk/pkg/integration"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/log"
	"github.com/newrelic/nri-bitmovin-analytics/internal/bitmovin"
)

var (
	/* Args below are populated via ldflags at build time */
	gIntegrationID      = "com.newrelic.labs.nri-bitmovin-analytics"
	gIntegrationName    = "New Relic Infrastructure Bitmovin Integration"
	gIntegrationVersion = "0.1.0"
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

func main() {
	// Create a new background context to use
	ctx := context.Background()

	// Create the integration with options
	i, err := integration.NewStandaloneIntegration(
		&gBuildInfo,
		gBuildInfo.Name,
		"",
		integration.WithInterval(60),
		integration.WithLicenseKey(),
	)
	fatalIfErr(err)

	err = bitmovin.InitPipelines(i)
	fatalIfErr(err)

	// Run the integration
	defer i.Shutdown(ctx)
 	err = i.Run(ctx)
	fatalIfErr(err)
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatalf(err)
	}
}
