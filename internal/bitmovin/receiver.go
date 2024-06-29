package bitmovin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/connectors"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/log"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/model"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/pipeline"
	"github.com/spf13/viper"
)

const (
	BaseURL = "https://api.bitmovin.com"
)

type BitmovinAuthenticator struct {
	credentials 		*BitmovinCredentials
}

func NewBitmovinAuthenticator(credentials *BitmovinCredentials) (
	*BitmovinAuthenticator,
) {
	return &BitmovinAuthenticator{ credentials }
}

func (b *BitmovinAuthenticator) Authenticate(
	connector *connectors.HttpConnector,
	req *http.Request,
) error {
	req.Header.Add("X-Api-Key", b.credentials.apiKey)
	req.Header.Add("X-Tenant-Org-Id", b.credentials.tenantOrg)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	return nil
}

func bitmovinPostBodyBuilder(
	licenseKey string,
	queryParams *BitmovinQueryParams,
	recvInterval time.Duration,
) connectors.HttpBodyBuilder {
	return func () (any, error) {
		end := time.Now().UTC()
		bufferDuration := time.Duration(recvInterval + 60)
		start := end.Add(-bufferDuration * time.Second)

		requestBody := &BitmovinRequestBody{
			Start:      start.Format(time.RFC3339Nano),
			End:        end.Format(time.RFC3339Nano),
			LicenseKey: licenseKey,
		}

		if queryParams.Filters != nil {
			requestBody.Filters = &queryParams.Filters
		}

		if queryParams.GroupBy != nil {
			requestBody.GroupBy = &queryParams.GroupBy
		}

		if queryParams.OrderBy != nil {
			requestBody.OrderBy = &queryParams.OrderBy
		}

		if queryParams.Interval != "" {
			requestBody.Interval = &queryParams.Interval
		}

		if queryParams.Metric != "" {
			requestBody.Metric = &queryParams.Metric
		}

		if queryParams.BMDimension != "" {
			requestBody.Dimension = &queryParams.BMDimension
		}

		requestBody.Percentile = queryParams.Percentile

		if log.IsDebugEnabled() {
			log.Debugf("request payload follows:")
			log.PrettyPrintJson(requestBody)
		}

		jsonValue, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(jsonValue), nil
	}
}

func bitmovinResponseDecoderBuilder(
	queryParams *BitmovinQueryParams,
	metricPrefix string,
) (
	pipeline.MetricsDecoderFunc,
) {
	return func(
		receiver pipeline.MetricsReceiver,
		in io.ReadCloser,
		out chan <- model.Metric,
	) error {
		apiResponse := BitmovinResponse{}

		log.Debugf("decoding bitmovin JSON response")

		dec := json.NewDecoder(in)

		err := dec.Decode(&apiResponse)
		if err != nil {
			return err
		}

		if log.IsDebugEnabled() {
			log.PrettyPrintJson(apiResponse)
		}

		// Generate the metrics
		LOOP:

		for i := 0; i < len(apiResponse.Data.Result.Rows); i += 1 {
			if len(apiResponse.Data.Result.Rows[i]) == 0 {
				log.Debugf("skipping response row %d because it is empty", i)
				continue
			}

			if apiResponse.Data.Result.Rows[i][0] == nil {
				log.Debugf("skipping row %d: first column is nil", i)
				break
			}

			// JSON numbers are always decoded as floats
			timestamp, ok := apiResponse.Data.Result.Rows[i][0].(float64)
			if !ok {
				log.Warnf("skipping row %d: timestamp column is not a float: %v", i, timestamp)
				continue
			}

			// Collect the dimensions.  They always come between the timestamp
			// and the value.

			j := 1
			dimensions := make(map[string]string)

			for _, dimension := range queryParams.GroupBy {
				val, ok := apiResponse.Data.Result.Rows[i][j].(string)
				if !ok {
					log.Warnf("skipping row %d: dimension column %s is not a string: %v", j, dimension, val)
					continue LOOP
				}
				dimensions[dimension] = val
				j += 1
			}

			// JSON numbers are always decoded as floats
			val, ok := apiResponse.Data.Result.Rows[i][j].(float64)
			if !ok {
				log.Warnf("skipping row %d: value column is not a float: %v", i, val)
				continue
			}

			metric := model.NewGaugeMetric(
				fmt.Sprintf("%s%s", metricPrefix, queryParams.NRMetric),
				model.MakeNumeric(val),
				time.Unix(int64(timestamp/1000), 0),
			)

			if len(dimensions) > 0 {
				for k, v := range(dimensions) {
					metric.Attributes[k] = v
				}
			}

			out <- metric
		}

		return nil
	}
}

func addReceiver(
	mp *pipeline.MetricsPipeline,
	authenticator *BitmovinAuthenticator,
	licenseKey string,
	metricPrefix string,
	recvInterval time.Duration,
	id string,
	queryParams *BitmovinQueryParams,
) {
	mp.AddReceiver(
		pipeline.NewSimpleReceiver(
			id,
			BaseURL + queryParams.URI,
			pipeline.WithAuthenticator(authenticator),
			pipeline.WithMethod("POST"),
			pipeline.WithBody(bitmovinPostBodyBuilder(
				licenseKey,
				queryParams,
				recvInterval,
			)),
			pipeline.WithMetricsDecoder(bitmovinResponseDecoderBuilder(
				queryParams,
				metricPrefix,
			)),
		),
	)
}

func buildFilters(query *BitmovinQuery) []map[string]any {
	filters := make([]map[string]any, 0)

	if query.Filters == nil {
		return filters
	}

	for k, v := range *query.Filters {
		filters = append(filters, map[string]any {
			"name": strings.ToUpper(k),
			"operator": strings.ToUpper(v.Operator),
			"value": v.Value,
		})
	}

	return filters
}

func buildGroupBy(query *BitmovinQuery) []string {
	groupBy := make([]string, 0)

	if query.Dimensions == nil {
		return groupBy
	}

	groupBy = *query.Dimensions

	return groupBy
}

func getInterval(query *BitmovinQuery) string {
	if query.Interval == nil {
		return "MINUTE"
	}

	return strings.ToUpper(*query.Interval)
}

func buildOrderBy(query *BitmovinQuery) []BitmovinOrderBy {
	if query.OrderBy == nil {
		return []BitmovinOrderBy{{Name: "MINUTE", Order: "DESC"}}
	}

	return *query.OrderBy
}

func addReceiverWithQuery(
	mp *pipeline.MetricsPipeline,
	authenticator *BitmovinAuthenticator,
	licenseKey string,
	metricPrefix string,
	recvInterval time.Duration,
	query *BitmovinQuery,
) error {
	var queryParams *BitmovinQueryParams

	switch query.Type {
	case "max_concurrentviewers":
		queryParams = &BitmovinQueryParams{
			URI:         	"/v1/analytics/metrics/max_concurrentviewers",
			NRMetric:    	"max_concurrent_viewers",
			Metric:		 	"max_concurrentviewers",
			BMDimension: 	"",
			Filters:     	buildFilters(query),
			GroupBy:	 	buildGroupBy(query),
			Interval:    	getInterval(query),
			OrderBy:     	buildOrderBy(query),
		}
	case "avg_concurrentviewers":
		queryParams = &BitmovinQueryParams{
			URI:         	"/v1/analytics/metrics/avg_concurrentviewers",
			NRMetric:    	"avg_concurrent_viewers",
			Metric:		 	"avg_concurrentviewers",
			BMDimension: 	"",
			Filters:     	buildFilters(query),
			GroupBy:	 	buildGroupBy(query),
			Interval:    	getInterval(query),
			OrderBy:     	buildOrderBy(query),
		}
	case "avg_dropped_frames":
		queryParams = &BitmovinQueryParams{
			URI:         	"/v1/analytics/metrics/avg_dropped_frames",
			NRMetric:    	"avg_dropped_frames",
			Metric:		 	"avg_dropped_frames",
			BMDimension: 	"",
			Filters:     	buildFilters(query),
			GroupBy:	 	buildGroupBy(query),
			Interval:    	getInterval(query),
			OrderBy:     	buildOrderBy(query),
		}
	case "count":
		queryParams = &BitmovinQueryParams{
			URI:         	"/v1/analytics/queries/count",
			NRMetric:    	fmt.Sprintf("cnt_%s", strings.ToLower(query.Metric)),
			Metric:			"",
			BMDimension: 	query.Metric,
			Filters:     	buildFilters(query),
			GroupBy:	 	buildGroupBy(query),
			Interval:    	getInterval(query),
			OrderBy:     	buildOrderBy(query),
		}
	case "sum":
		queryParams = &BitmovinQueryParams{
			URI:         	"/v1/analytics/queries/sum",
			NRMetric:    	fmt.Sprintf("sum_%s", strings.ToLower(query.Metric)),
			Metric:			"",
			BMDimension: 	query.Metric,
			Filters:     	buildFilters(query),
			GroupBy:	 	buildGroupBy(query),
			Interval:    	getInterval(query),
			OrderBy:     	buildOrderBy(query),
		}
	case "average":
		queryParams = &BitmovinQueryParams{
			URI:         	"/v1/analytics/queries/avg",
			NRMetric:    	fmt.Sprintf("avg_%s", strings.ToLower(query.Metric)),
			Metric:			"",
			BMDimension: 	query.Metric,
			Filters:     	buildFilters(query),
			GroupBy:	 	buildGroupBy(query),
			Interval:    	getInterval(query),
			OrderBy:     	buildOrderBy(query),
		}
	case "min":
		queryParams = &BitmovinQueryParams{
			URI:         	"/v1/analytics/queries/min",
			NRMetric:    	fmt.Sprintf("min_%s", strings.ToLower(query.Metric)),
			Metric:			"",
			BMDimension: 	query.Metric,
			Filters:     	buildFilters(query),
			GroupBy:	 	buildGroupBy(query),
			Interval:    	getInterval(query),
			OrderBy:     	buildOrderBy(query),
		}
	case "max":
		queryParams = &BitmovinQueryParams{
			URI:         	"/v1/analytics/queries/max",
			NRMetric:    	fmt.Sprintf("max_%s", strings.ToLower(query.Metric)),
			Metric:			"",
			BMDimension: 	query.Metric,
			Filters:     	buildFilters(query),
			GroupBy:	 	buildGroupBy(query),
			Interval:    	getInterval(query),
			OrderBy:     	buildOrderBy(query),
		}
	case "stddev":
		queryParams = &BitmovinQueryParams{
			URI:         	"/v1/analytics/queries/stddev",
			NRMetric:    	fmt.Sprintf("stddev_%s", strings.ToLower(query.Metric)),
			Metric:			"",
			BMDimension: 	query.Metric,
			Filters:     	buildFilters(query),
			GroupBy:	 	buildGroupBy(query),
			Interval:    	getInterval(query),
			OrderBy:     	buildOrderBy(query),
		}
	case "percentile":
		// @TODO: add percentile
		queryParams = &BitmovinQueryParams{
			URI:         	"/v1/analytics/queries/percentile",
			NRMetric:    	fmt.Sprintf("p%d_%s", *query.Percentile, strings.ToLower(query.Metric)),
			Metric:			"",
			BMDimension: 	query.Metric,
			Filters:     	buildFilters(query),
			GroupBy:	 	buildGroupBy(query),
			Interval:    	getInterval(query),
			OrderBy:     	buildOrderBy(query),
			Percentile:		query.Percentile,
		}
	case "variance":
		// @TODO: add percentile
		queryParams = &BitmovinQueryParams{
			URI:         	"/v1/analytics/queries/variance",
			NRMetric:    	fmt.Sprintf("var_%s", strings.ToLower(query.Metric)),
			Metric:			"",
			BMDimension: 	query.Metric,
			Filters:     	buildFilters(query),
			GroupBy:	 	buildGroupBy(query),
			Interval:    	getInterval(query),
			OrderBy:     	buildOrderBy(query),
		}
	case "median":
		queryParams = &BitmovinQueryParams{
			URI:         	"/v1/analytics/queries/median",
			NRMetric:    	fmt.Sprintf("med_%s", strings.ToLower(query.Metric)),
			Metric:			"",
			BMDimension: 	query.Metric,
			Filters:     	buildFilters(query),
			GroupBy:	 	buildGroupBy(query),
			Interval:    	getInterval(query),
			OrderBy:     	buildOrderBy(query),
		}
	}

	addReceiver(
		mp,
		authenticator,
		licenseKey,
		metricPrefix,
		recvInterval,
		queryParams.NRMetric,
		queryParams,
	)

	return nil
}

func setupReceivers(
	mp *pipeline.MetricsPipeline,
	credentials *BitmovinCredentials,
	recvInterval time.Duration,
	queries []BitmovinQuery,
) error {
	metricPrefix := viper.GetString("bitmovinMetricPrefix")
	authenticator := NewBitmovinAuthenticator(credentials)

	for _, query := range queries {
		err := addReceiverWithQuery(
			mp,
			authenticator,
			credentials.licenseKey,
			metricPrefix,
			recvInterval,
			&query,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
