package bitmovin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/connectors"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/log"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/model"
	"github.com/spf13/viper"
)

const (
	DEFAULT_NULL_VALUE_STRING = "<NULL>"
	BaseURL                   = "https://api.bitmovin.com"
	RESULT_LIMIT              = 200 // The maximum allowed by the Bitmovin API
)

type BitmovinAuthenticator struct {
	credentials *BitmovinCredentials
}

func NewBitmovinAuthenticator(credentials *BitmovinCredentials) *BitmovinAuthenticator {
	return &BitmovinAuthenticator{credentials}
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
) connectors.HttpBodyBuilder {
	return func() (any, error) {
		requestBody := &BitmovinRequestBody{
			Start:      queryParams.StartTime,
			End:        queryParams.EndTime,
			LicenseKey: licenseKey,
			Limit:      queryParams.Limit,
			Offset:     queryParams.Offset,
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

func buildFilters(query *BitmovinQuery) []map[string]any {
	filters := make([]map[string]any, 0)

	if query.Filters == nil {
		return filters
	}

	for k, v := range *query.Filters {
		filters = append(filters, map[string]any{
			"name":     strings.ToUpper(k),
			"operator": strings.ToUpper(v.Operator),
			"value":    v.Value,
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

func getBitmovinTimeout() uint {
	timeout := viper.GetUint("bitmovinTimeout")
	if timeout <= 0 {
		timeout = 10
	}
	return timeout
}

func calcQueryStartEnd(recvInterval time.Duration) (string, string) {
	end := time.Now().UTC()
	bufferDuration := time.Duration(recvInterval + 60)
	start := end.Add(-bufferDuration * time.Second)
	return start.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano)
}

func buildQuery(
	queryPos int,
	query *BitmovinQuery,
	recvInterval time.Duration,
) *BitmovinQueryParams {
	var queryParams *BitmovinQueryParams

	if query.Name == "" {
		log.Warnf("automatic metric names are DEPRECATED, add a 'name' parameter to your config file for the query at position %d", queryPos+1)
	}

	switch query.Type {
	case "max_concurrentviewers":
		queryParams = &BitmovinQueryParams{
			URI:         "/v1/analytics/metrics/max_concurrentviewers",
			NRMetric:    "max_concurrent_viewers",
			Metric:      "max_concurrentviewers",
			BMDimension: "",
		}
	case "avg_concurrentviewers":
		queryParams = &BitmovinQueryParams{
			URI:         "/v1/analytics/metrics/avg_concurrentviewers",
			NRMetric:    "avg_concurrent_viewers",
			Metric:      "avg_concurrentviewers",
			BMDimension: "",
		}
	case "avg_dropped_frames":
		queryParams = &BitmovinQueryParams{
			URI:         "/v1/analytics/metrics/avg_dropped_frames",
			NRMetric:    "avg_dropped_frames",
			Metric:      "avg_dropped_frames",
			BMDimension: "",
		}
	case "count":
		queryParams = &BitmovinQueryParams{
			URI:         "/v1/analytics/queries/count",
			NRMetric:    fmt.Sprintf("cnt_%s", strings.ToLower(query.Metric)),
			Metric:      "",
			BMDimension: query.Metric,
		}
	case "sum":
		queryParams = &BitmovinQueryParams{
			URI:         "/v1/analytics/queries/sum",
			NRMetric:    fmt.Sprintf("sum_%s", strings.ToLower(query.Metric)),
			Metric:      "",
			BMDimension: query.Metric,
		}
	case "average":
		queryParams = &BitmovinQueryParams{
			URI:         "/v1/analytics/queries/avg",
			NRMetric:    fmt.Sprintf("avg_%s", strings.ToLower(query.Metric)),
			Metric:      "",
			BMDimension: query.Metric,
		}
	case "min":
		queryParams = &BitmovinQueryParams{
			URI:         "/v1/analytics/queries/min",
			NRMetric:    fmt.Sprintf("min_%s", strings.ToLower(query.Metric)),
			Metric:      "",
			BMDimension: query.Metric,
		}
	case "max":
		queryParams = &BitmovinQueryParams{
			URI:         "/v1/analytics/queries/max",
			NRMetric:    fmt.Sprintf("max_%s", strings.ToLower(query.Metric)),
			Metric:      "",
			BMDimension: query.Metric,
		}
	case "stddev":
		queryParams = &BitmovinQueryParams{
			URI:         "/v1/analytics/queries/stddev",
			NRMetric:    fmt.Sprintf("stddev_%s", strings.ToLower(query.Metric)),
			Metric:      "",
			BMDimension: query.Metric,
		}
	case "percentile":
		// @TODO: add percentile
		queryParams = &BitmovinQueryParams{
			URI:         "/v1/analytics/queries/percentile",
			NRMetric:    fmt.Sprintf("p%d_%s", *query.Percentile, strings.ToLower(query.Metric)),
			Metric:      "",
			BMDimension: query.Metric,
			Percentile:  query.Percentile,
		}
	case "variance":
		// @TODO: add percentile
		queryParams = &BitmovinQueryParams{
			URI:         "/v1/analytics/queries/variance",
			NRMetric:    fmt.Sprintf("var_%s", strings.ToLower(query.Metric)),
			Metric:      "",
			BMDimension: query.Metric,
		}
	case "median":
		queryParams = &BitmovinQueryParams{
			URI:         "/v1/analytics/queries/median",
			NRMetric:    fmt.Sprintf("med_%s", strings.ToLower(query.Metric)),
			Metric:      "",
			BMDimension: query.Metric,
		}
	}

	// Set attrs common to all queries
	queryParams.Name = query.Name
	queryParams.Filters = buildFilters(query)
	queryParams.GroupBy = buildGroupBy(query)
	queryParams.Interval = getInterval(query)
	queryParams.OrderBy = buildOrderBy(query)
	queryParams.Limit = RESULT_LIMIT
	queryParams.Offset = 0
	queryParams.StartTime, queryParams.EndTime = calcQueryStartEnd(recvInterval)

	return queryParams
}

func newQueryConnector(
	authenticator *BitmovinAuthenticator,
	queryParams *BitmovinQueryParams,
	licenseKey string,
) *connectors.HttpConnector {
	connector := connectors.NewHttpGetConnector(BaseURL + queryParams.URI)
	connector.SetAuthenticator(authenticator)
	connector.SetMethod("POST")
	connector.SetBody(bitmovinPostBodyBuilder(
		licenseKey,
		queryParams,
	))
	connector.SetTimeout(time.Duration(getBitmovinTimeout()) * time.Second)

	return connector
}

// TODO: refactor this function, it's a bit spaghetti
func decodeAndSendResponse(
	queryParams *BitmovinQueryParams,
	metricPrefix string,
	in io.ReadCloser,
	out chan<- model.Metric,
) (uint, error) {
	apiResponse := BitmovinResponse{}

	log.Debugf("decoding bitmovin JSON response")

	nullValueString := viper.GetString("bitmovinNullValueString")
	if nullValueString == "" {
		nullValueString = DEFAULT_NULL_VALUE_STRING
	}

	dec := json.NewDecoder(in)

	err := dec.Decode(&apiResponse)
	if err != nil {
		return 0, err
	}

	if log.IsDebugEnabled() {
		log.PrettyPrintJson(apiResponse)
	}

	bitmovinMetricName := queryParams.Metric
	if bitmovinMetricName == "" {
		bitmovinMetricName = queryParams.BMDimension
	}

	numResults := uint(0)

	// Generate the metrics
LOOP:

	for i := 0; i < len(apiResponse.Data.Result.Rows); i += 1 {
		colCount := len(apiResponse.Data.Result.Rows[i])
		if colCount == 0 {
			log.Debugf(
				"skipping row %d for metric %s because it is empty",
				i,
				bitmovinMetricName,
			)
			continue
		}

		// Validate we have the number of columns we think we should
		// timestamp column + # of dimension columns + value column
		expectedColCount := 1 + len(queryParams.GroupBy) + 1
		if colCount != expectedColCount {
			log.Warnf(
				"skipping row %d for metric %s: unexpected number of columns: %d: expected number of columns: %d",
				i,
				bitmovinMetricName,
				colCount,
				expectedColCount,
			)
			continue
		}

		// JSON numbers are always decoded as floats
		timestamp, ok := apiResponse.Data.Result.Rows[i][0].(float64)
		if !ok {
			log.Warnf(
				"skipping row %d for metric %s: timestamp column is not a float: %T",
				i,
				bitmovinMetricName,
				apiResponse.Data.Result.Rows[i][0],
			)
			continue
		}

		// Collect the dimensions.  They always come between the timestamp
		// and the value.

		j := 1
		dimensions := make(map[string]interface{})

		for _, dimension := range queryParams.GroupBy {
			dimVal := apiResponse.Data.Result.Rows[i][j]
			switch v := dimVal.(type) {
			// Bitmovin dimensions can be 3 data types: string, int, bool
			// JSON numbers are always decoded as floats
			case string, bool, float64:
				dimensions[dimension] = v
				j += 1
			case nil:
				// @todo: is there a better way to handle null group
				// values?
				dimensions[dimension] = nullValueString
				j += 1
			default:
				log.Warnf(
					"skipping row %d for metric %s: dimension column %s is not a valid type: %T",
					i,
					bitmovinMetricName,
					dimension,
					dimVal,
				)
				continue LOOP
			}
		}

		// The metric value is a number and JSON numbers are always decoded
		// as floats
		val, ok := apiResponse.Data.Result.Rows[i][j].(float64)
		if !ok {
			log.Warnf(
				"skipping row %d for metric %s: value column is not a float: %T",
				i,
				bitmovinMetricName,
				apiResponse.Data.Result.Rows[i][j],
			)
			continue
		}

		var metricName string

		if queryParams.Name == "" {
			metricName = fmt.Sprintf("%s%s", metricPrefix, queryParams.NRMetric)
		} else {
			metricName = fmt.Sprintf("%s%s", metricPrefix, queryParams.Name)
		}

		metric := model.NewGaugeMetric(
			metricName,
			model.MakeNumeric(val),
			time.Unix(int64(timestamp/1000), 0),
		)

		if len(dimensions) > 0 {
			for k, v := range dimensions {
				metric.Attributes[k] = v
			}
		}

		out <- metric
		numResults += 1
	}

	return numResults, nil
}

type BitmovinMetricReceiver struct {
	queries      []BitmovinQuery
	credentials  *BitmovinCredentials
	recvInterval time.Duration
}

func NewBitmovinReceiver(
	credentials *BitmovinCredentials,
	recvInterval time.Duration,
	queries []BitmovinQuery,
) *BitmovinMetricReceiver {
	return &BitmovinMetricReceiver{
		queries:      queries,
		credentials:  credentials,
		recvInterval: recvInterval,
	}
}

//// MetricsReceiver interface implementation \\\\

func (r *BitmovinMetricReceiver) GetId() string {
	return "bitmovin-metric-receiver"
}

func (r *BitmovinMetricReceiver) PollMetrics(
	ctx context.Context,
	metricChan chan<- model.Metric,
) error {

	metricPrefix := viper.GetString("bitmovinMetricPrefix")
	authenticator := NewBitmovinAuthenticator(r.credentials)

	for queryIndex := range r.queries {
		log.Debugf("Request query at index %d", queryIndex)

		queryParams := buildQuery(queryIndex, &r.queries[queryIndex], r.recvInterval)
		lastPage := false

		for !lastPage {
			connector := newQueryConnector(
				authenticator,
				queryParams,
				r.credentials.licenseKey,
			)
			data, err := connector.Request()
			if err != nil {
				return err
			}

			numResults, err := decodeAndSendResponse(
				queryParams,
				metricPrefix,
				data,
				metricChan,
			)
			if err != nil {
				return err
			}

			// We may not have all results
			if numResults == queryParams.Limit {
				log.Debugf("Slide the offset and get the next page")
				queryParams.Offset += queryParams.Limit
			} else if numResults < queryParams.Limit {
				log.Debugf("No more pages")
				lastPage = true
			} else {
				// Something went really wrong
				return fmt.Errorf("We have more results than the expected, data might be inconsistent or corrupted")
			}
		}
	}

	return nil
}
