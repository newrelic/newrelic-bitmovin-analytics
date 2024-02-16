package bitmovin

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/deser"
	"newrelic/multienv/pkg/model"
	"time"
)

func makeHttpConnectors() (config.RecvConfig, error) {

	var connectors []connect.Connector

	now := time.Now().UTC()
	bitmovinTimestamps.EndTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)
	bitmovinTimestamps.StartTime = bitmovinTimestamps.EndTime.Add(time.Duration(-recv_interval) * time.Second)

	headers := make(map[string]string)
	headers["X-Api-Key"] = bitmovinCreds.APIKey
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	headers["X-Tenant-Org-Id"] = bitmovinCreds.TenantOrg

	for _, metric := range metricTypes {

		requestBody := RequestBody{
			Start:          bitmovinTimestamps.StartTime.Format(time.RFC3339),
			End:            bitmovinTimestamps.EndTime.Format(time.RFC3339),
			LicenseKey:     bitmovinCreds.LicenseKey,
			GroupBy:        []string{"EXPERIMENT_NAME", "CDN_PROVIDER", "COUNTRY", "DEVICE_TYPE"},
			IncludeContext: false,
			Interval:       "MINUTE",
			Dimension:      metric.BMDimension,
			Limit:          200,
			Offset:         0,
			Filters:        metric.Filters,
		}

		jsonValue, err := json.Marshal(requestBody)
		if err != nil {
			log.Println(err)
		}

		connector := connect.MakeHttpPostConnector(BaseURL+metric.URI, jsonValue, headers)
		connector.SetConnectorName(metric.NRMetric)
		connectors = append(connectors, &connector)
	}

	return config.RecvConfig{
		Connectors: connectors,
		Deser:      deser.DeserJson,
	}, nil
}

func InitRecv(pipeConfig *config.PipelineConfig) (config.RecvConfig, error) {
	recv_interval = int(pipeConfig.Interval)
	if recv_interval == 0 {
		log.Warn("Interval not set, using 5 seconds")
		recv_interval = 5
	}

	// Set required bitmovin credentials
	if bitmovinApiKey, ok := pipeConfig.GetString("bitmovin_api_key"); ok {
		bitmovinCreds.APIKey = bitmovinApiKey
	} else {
		return config.RecvConfig{}, errors.New("config key 'bitmovin_api_key' doesn't exist")
	}

	if bitmovinLicenseKey, ok := pipeConfig.GetString("bitmovin_license_key"); ok {
		bitmovinCreds.LicenseKey = bitmovinLicenseKey
	} else {
		return config.RecvConfig{}, errors.New("config key 'bitmovin_license_key' doesn't exist")
	}

	if bitmovinTenantOrg, ok := pipeConfig.GetString("bitmovin_tenant_org"); ok {
		bitmovinCreds.TenantOrg = bitmovinTenantOrg
	} else {
		return config.RecvConfig{}, errors.New("config key 'bitmovin_tenant_org' doesn't exist")
	}

	return makeHttpConnectors()
}

func InitProc(pipeConfig *config.PipelineConfig) (config.ProcConfig, error) {
	recv_interval = int(pipeConfig.Interval)
	if recv_interval == 0 {
		log.Warn("Interval not set, using 5 seconds")
		recv_interval = 5
	}

	return config.ProcConfig{
		Model: Response{},
	}, nil
}

// Proc Generate all kinds of data.
func Proc(data any) []model.MeltModel {

	out := make([]model.MeltModel, 0)

	log.Println("DATA ===== ", data)
	if apiResponse, ok := data.(Response); ok {
		for i := 0; i < apiResponse.Data.Result.RowCount; i++ {
			metricValue := model.Numeric{}

			timestamp, ok := apiResponse.Data.Result.Rows[i][0].(float64)
			if !ok {
				log.Println("Timestamp value is not an float.", timestamp)

			}

			streamFormat := apiResponse.Data.Result.Rows[i][1]
			cdn := apiResponse.Data.Result.Rows[i][2]
			country := apiResponse.Data.Result.Rows[i][3]
			deviceType := apiResponse.Data.Result.Rows[i][4]

			switch apiResponse.Data.Result.Rows[i][5].(type) {
			case int:
				value, ok := apiResponse.Data.Result.Rows[i][5].(int64)
				if !ok {
					log.Println("Value is not an int.", value)

				}
				metricValue = model.MakeIntNumeric(value)
			case float64:

				value, ok := apiResponse.Data.Result.Rows[i][5].(float64)
				if !ok {
					log.Fatalf("Can't convert to float64")
				}
				metricValue = model.MakeFloatNumeric(value)

			default:
				log.Println("Value type unknown ===", apiResponse.Data.Result.Rows[i][5])
				return nil
			}

			meltMetric := createMetricRequest(
				apiResponse.Metric,
				time.Unix(int64(timestamp/1000), 0),
				metricValue,
				streamFormat,
				cdn,
				country,
				deviceType)
			out = append(out, meltMetric)
		}
	} else {
		log.Warn("Unknown type for data = ", data)
	}

	return out
}

func createMetricRequest(metric string, timestamp time.Time, value model.Numeric, streamFormat any, cdn any, country any, deviceType any) model.MeltModel {

	meltMetric := model.MakeGaugeMetric(
		metric, value, timestamp)
	meltMetric.Attributes = map[string]any{
		"stream_format": streamFormat,
		"cdn":           cdn,
		"country":       country,
		"device_type":   deviceType,
	}

	return meltMetric
}
