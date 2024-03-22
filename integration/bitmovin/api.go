package bitmovin

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/deser"
	"newrelic/multienv/pkg/model"
	"time"
)

func createBuilderRequest(metric Metric) (*http.Request, error) {
	end := time.Now().UTC()
	bufferDuration := time.Duration(recv_interval + 60)
	start := end.Add(-bufferDuration * time.Second)

	seconds := end.Sub(start).Seconds()

	log.Println("TIME DIFF ", seconds)
	log.Println(start.Format(time.RFC3339Nano))
	log.Println(end.Format(time.RFC3339Nano))

	requestBody := &RequestBody{
		Start:      start.Format(time.RFC3339Nano),
		End:        end.Format(time.RFC3339Nano),
		LicenseKey: bitmovinCreds.LicenseKey,
	}

	if metric.Filters != nil {
		requestBody.Filters = &metric.Filters
	}

	if metric.OrderBy != nil {
		requestBody.OrderBy = &metric.OrderBy
	}

	if metric.Interval != "" {
		requestBody.Interval = &metric.Interval
	}

	if metric.Metric != "" {
		requestBody.Metric = &metric.Metric
	}

	if metric.BMDimension != "" {
		requestBody.Dimension = &metric.BMDimension
	}

	jsonValue, err := json.Marshal(requestBody)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(jsonValue))

	req, err := http.NewRequest("POST", BaseURL+metric.URI, bytes.NewReader(jsonValue))
	req.Header = getBitmovinRequestHeaders()
	if err != nil {
		return nil, err
	}
	return req, nil
}

func GetMaxConcurrentViewersBuilder(conf *connect.HttpConfig) (*http.Request, error) {

	metric := Metric{
		Name:        "CONCURRENTS",
		URI:         "/v1/analytics/metrics/max_concurrentviewers",
		NRMetric:    "bitmovin.cnt_max_concurrent_viewers",
		BMDimension: "",
		Metric:      "max_concurrentviewers",
		Filters:     []map[string]any{},
		Interval:    "MINUTE",
		OrderBy:     []OrderBy{{Name: "MINUTE", Order: "DESC"}},
	}

	return createBuilderRequest(metric)
}

func GetRebufferPercentageBuilder(conf *connect.HttpConfig) (*http.Request, error) {

	metric := Metric{
		Name:        "REBUFFER",
		URI:         "/v1/analytics/queries/avg",
		NRMetric:    "bitmovin.avg_rebuffer_percentage",
		Metric:      "",
		BMDimension: "REBUFFER_PERCENTAGE",
		Filters:     []map[string]any{},
		Interval:    "MINUTE",
		OrderBy:     []OrderBy{{Name: "MINUTE", Order: "DESC"}},
	}

	return createBuilderRequest(metric)
}

func GetPlayerAttemptsBuilder(conf *connect.HttpConfig) (*http.Request, error) {

	metric := Metric{
		Name:        "PLAY ATTEMPTS",
		URI:         "/v1/analytics/queries/count",
		NRMetric:    "bitmovin.cnt_play_attempts",
		BMDimension: "PLAY_ATTEMPTS",
		Interval:    "MINUTE",
	}

	return createBuilderRequest(metric)
}

func GetVideoStartFailsBuilder(conf *connect.HttpConfig) (*http.Request, error) {

	metric := Metric{
		Name:        "VIDEO START FAILURES",
		URI:         "/v1/analytics/queries/count",
		NRMetric:    "bitmovin.cnt_video_start_failures",
		Metric:      "",
		BMDimension: "VIDEOSTART_FAILED",
		Filters: []map[string]any{{"name": "VIDEOSTART_FAILED_REASON", "operator": "NE", "value": "PAGE_CLOSED"},
			{"name": "VIDEOSTART_FAILED", "operator": "EQ", "value": true}},
		Interval: "MINUTE",
		OrderBy:  []OrderBy{{Name: "MINUTE", Order: "DESC"}},
	}

	return createBuilderRequest(metric)
}

func GetVideoStartTimeBuilder(conf *connect.HttpConfig) (*http.Request, error) {

	metric := Metric{
		Name:        "VIDEO START TIME",
		URI:         "/v1/analytics/queries/median",
		NRMetric:    "bitmovin.avg_video_startup_time_ms",
		BMDimension: "VIDEO_STARTUPTIME",
		Filters:     []map[string]any{{"name": "VIDEO_STARTUPTIME", "operator": "GT", "value": 0}},
		Interval:    "MINUTE",
	}

	return createBuilderRequest(metric)
}

func GetVideoBitrateBuilder(conf *connect.HttpConfig) (*http.Request, error) {

	metric := Metric{
		Name:        "VIDEO BITRATE",
		URI:         "/v1/analytics/queries/avg",
		NRMetric:    "bitmovin.avg_video_bitrate_mbps",
		Metric:      "",
		BMDimension: "VIDEO_BITRATE",
		Filters:     []map[string]any{{"name": "VIDEO_BITRATE", "operator": "GT", "value": 0}},
		Interval:    "MINUTE",
		OrderBy:     []OrderBy{{Name: "MINUTE", Order: "DESC"}},
	}

	return createBuilderRequest(metric)
}

func GetAverageViewTimeBuilder(conf *connect.HttpConfig) (*http.Request, error) {

	metric := Metric{
		Name:        "AVERAGE VIEW TIME",
		URI:         "/v1/analytics/queries/avg",
		NRMetric:    "bitmovin.avg_view_time",
		Metric:      "",
		BMDimension: "VIEWTIME",
		Filters:     []map[string]any{},
		Interval:    "MINUTE",
		OrderBy:     []OrderBy{{Name: "MINUTE", Order: "DESC"}},
	}

	return createBuilderRequest(metric)
}

func InitRecv(pipeConfig *config.PipelineConfig) (config.RecvConfig, error) {
	recv_interval = int(pipeConfig.Interval)
	if recv_interval == 0 {
		log.Warn("Interval not set, using 5 seconds")
		recv_interval = 60
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

	// Connectors
	getConcurrentViewers := connect.MakeHttpConnectorWithBuilder(GetMaxConcurrentViewersBuilder)
	getConcurrentViewers.SetConnectorName("bitmovin.cnt_max_concurrent_viewers")

	getRebufferPercentage := connect.MakeHttpConnectorWithBuilder(GetRebufferPercentageBuilder)
	getRebufferPercentage.SetConnectorName("bitmovin.avg_rebuffer_percentage")

	getPlayAttempts := connect.MakeHttpConnectorWithBuilder(GetPlayerAttemptsBuilder)
	getPlayAttempts.SetConnectorName("bitmovin.cnt_play_attempts")

	getVideoStartFails := connect.MakeHttpConnectorWithBuilder(GetVideoStartFailsBuilder)
	getVideoStartFails.SetConnectorName("bitmovin.cnt_video_start_failures")

	getVideoStartTime := connect.MakeHttpConnectorWithBuilder(GetVideoStartTimeBuilder)
	getVideoStartTime.SetConnectorName("bitmovin.avg_video_startup_time_ms")

	getVideoBitrate := connect.MakeHttpConnectorWithBuilder(GetVideoBitrateBuilder)
	getVideoBitrate.SetConnectorName("bitmovin.avg_video_bitrate_mbps")

	getViewTime := connect.MakeHttpConnectorWithBuilder(GetAverageViewTimeBuilder)
	getViewTime.SetConnectorName("bitmovin.avg_view_time")

	return config.RecvConfig{
		Connectors: []connect.Connector{
			&getConcurrentViewers,
			&getRebufferPercentage,
			&getPlayAttempts,
			&getVideoStartFails,
			&getVideoStartTime,
			&getVideoBitrate,
			&getViewTime,
		},
		Deser: deser.DeserJson,
	}, nil
}

func InitProc(pipeConfig *config.PipelineConfig) (config.ProcConfig, error) {
	recv_interval = int(pipeConfig.Interval)
	if recv_interval == 0 {
		log.Warn("Interval not set, using 60 seconds")
		recv_interval = 60
	}

	return config.ProcConfig{
		Model: Response{},
	}, nil
}

// Proc Generate all kinds of data.
func Proc(data any) []model.MeltModel {

	out := make([]model.MeltModel, 0)

	if apiResponse, ok := data.(Response); ok {
		for i := 0; i < len(apiResponse.Data.Result.Rows); i++ {

			if len(apiResponse.Data.Result.Rows[i]) > 0 {
				metricValue := model.Numeric{}

				if apiResponse.Data.Result.Rows[i][0] == nil {
					return nil
				} else {

					timestamp, ok := apiResponse.Data.Result.Rows[i][0].(float64)
					if !ok {
						log.Println("Timestamp value is not an float.", timestamp)

					}

					switch apiResponse.Data.Result.Rows[i][1].(type) {
					case int:
						value, ok := apiResponse.Data.Result.Rows[i][1].(int64)
						if !ok {
							log.Println("Value is not an int.", value)

						}
						metricValue = model.MakeIntNumeric(value)
					case float64:

						value, ok := apiResponse.Data.Result.Rows[i][1].(float64)
						if !ok {
							log.Fatalf("Can't convert to float64")
						}
						metricValue = model.MakeFloatNumeric(value)

					default:
						log.Println("Value type unknown ===", apiResponse.Data.Result.Rows[i][1])
						return nil
					}

					meltMetric := createMetricRequest(
						apiResponse.Metric,
						time.Unix(int64(timestamp/1000), 0),
						metricValue)
					out = append(out, meltMetric)
				}
			} else {
				log.Warn("No data to send to NR = ", data)
			}
		}
	} else {
		log.Warn("Unknown type for data = ", data)
	}

	return out
}

func createMetricRequest(metric string, timestamp time.Time, value model.Numeric) model.MeltModel {

	meltMetric := model.MakeGaugeMetric(
		metric, value, timestamp)
	meltMetric.Attributes = map[string]any{
		"instrumentation.name": "nri-bitmovin-analytics",
	}

	return meltMetric
}
