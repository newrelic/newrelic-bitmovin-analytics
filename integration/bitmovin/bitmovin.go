package bitmovin

import (
	_ "fmt"
	"time"
)

var recv_interval = 0
var bitmovinTimestamps = BitmovinTimestamps{}
var bitmovinCreds = BitmovinCreds{}

// BitmovinCredentials offers methods for interaction with Bitmovin Analytics API
type BitmovinCreds struct {
	APIKey     string
	LicenseKey string
	TenantOrg  string
}

type BitmovinTimestamps struct {
	StartTime time.Time
	EndTime   time.Time
}

type Metric struct {
	Name        string
	URI         string
	NRMetric    string
	BMDimension string
	Filters     []map[string]any
}

type Result struct {
	RowCount     int       `json:"rowCount"`
	Rows         []RowData `json:"rows"`
	ColumnLabels []Column  `json:"columnLabels"`
	JobID        string    `json:"jobId"`
	TaskCount    int       `json:"taskCount"`
}

type Column struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type RowData []any

type Data struct {
	Result   Result    `json:"result"`
	Messages []Message `json:"messages"`
}

type Response struct {
	RequestID string `json:"requestId"`
	Status    string `json:"status"`
	Data      Data   `json:"data"`
	Metric    string `json:"metric"`
}

type RequestBody struct {
	Start          string           `json:"start"`
	End            string           `json:"end"`
	LicenseKey     string           `json:"licenseKey"`
	Dimension      string           `json:"dimension"`
	GroupBy        []string         `json:"groupBy"`
	IncludeContext bool             `json:"includeContext"`
	Limit          int              `json:"limit"`
	Offset         int              `json:"offset"`
	Interval       string           `json:"interval"`
	Filters        []map[string]any `json:"filters"`
}

type Message struct {
	ID   string `json:"id"`
	Date string `json:"date"`
	Type string `json:"type"`
	Text string `json:"text"`
}

// BaseURL Bitmovin package constant
const (
	BaseURL = "https://api.bitmovin.com"
)

var metricTypes = []Metric{
	{
		Name:        "CONCURRENTS",
		URI:         "/v1/analytics/metrics/max-concurrentviewers",
		NRMetric:    "bitmovin.cnt_max_concurrent_viewers",
		BMDimension: "",
		Filters:     []map[string]any{},
	},
	{
		Name:        "REBUFFER",
		URI:         "/v1/analytics/queries/avg",
		NRMetric:    "bitmovin.avg_rebuffer_percentage",
		BMDimension: "REBUFFER_PERCENTAGE",
		Filters:     []map[string]any{},
	},
	{
		Name:        "PLAY ATTEMPTS",
		URI:         "/v1/analytics/queries/count",
		NRMetric:    "bitmovin.cnt_play_attempts",
		BMDimension: "PLAY_ATTEMPTS",
		Filters:     []map[string]any{},
	},
	{
		Name:        "VIDEO START FAILURES",
		URI:         "/v1/analytics/queries/count",
		NRMetric:    "bitmovin.cnt_video_start_failures",
		BMDimension: "VIDEOSTART_FAILED",
		Filters: []map[string]any{{"name": "VIDEOSTART_FAILED_REASON", "operator": "NE", "value": "PAGE_CLOSED"},
			{"name": "VIDEOSTART_FAILED", "operator": "EQ", "value": true}},
	},
	{
		Name:        "VIDEO START TIME",
		URI:         "/v1/analytics/queries/avg",
		NRMetric:    "bitmovin.avg_video_startup_time_ms",
		BMDimension: "VIDEO_STARTUPTIME",
		Filters:     []map[string]any{},
	},
	{
		Name:        "VIDEO BITRATE",
		URI:         "/v1/analytics/queries/avg",
		NRMetric:    "bitmovin.avg_video_bitrate_mbps",
		BMDimension: "VIDEO_BITRATE",
		Filters:     []map[string]any{},
	},
}
