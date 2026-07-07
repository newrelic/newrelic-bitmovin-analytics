package bitmovin

import (
	"fmt"

	"github.com/newrelic/newrelic-labs-sdk/pkg/integration"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/exporters"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/pipeline"
	"github.com/spf13/viper"
)

func InitPipelines(i *integration.LabsIntegration) error {
	// Load credentials
	apiKey := viper.GetString("bitmovinApiKey")
	if apiKey == "" {
		return fmt.Errorf("missing bitmovin API key")
	}

	licenseKey := viper.GetString("bitmovinLicenseKey")
	if licenseKey == "" {
		return fmt.Errorf("missing bitmovin license key")
	}

	tenantOrg := viper.GetString("bitmovinTenantOrg")
	if tenantOrg == "" {
		return fmt.Errorf("missing bitmovin tenant org")
	}

	// Load bitmovin query config
	var queries []BitmovinQuery

	err := viper.UnmarshalKey("queries", &queries)
	if err != nil {
		return fmt.Errorf("parse queries failed: %w", err)
	}

	// Create metrics pipeline
	mp := pipeline.NewMetricsPipeline()

	// Create newrelic exporter
	newrelicExporter := exporters.NewNewRelicExporter(
		"newrelic-api",
		i,
	)
	mp.AddExporter(newrelicExporter)

	// Create bitmovin metric receiver
	bitmovinMetricReceiver := NewBitmovinReceiver(
		&BitmovinCredentials{
			apiKey,
			licenseKey,
			tenantOrg,
		},
		i.Interval,
		queries,
	)
	mp.AddReceiver(bitmovinMetricReceiver)

	// Register the pipeline with the integration
	i.AddPipeline(mp)

	return nil
}
