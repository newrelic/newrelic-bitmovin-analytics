package bitmovin

import (
	_ "fmt"
)

var recv_interval = 0
var bitmovinCreds = BitmovinCreds{}

// BitmovinCredentials offers methods for interaction with Bitmovin Analytics API
type BitmovinCreds struct {
	APIKey     string
	LicenseKey string
	TenantOrg  string
}

type Metric struct {
	Name        string
	URI         string
	NRMetric    string
	BMDimension string
	Metric      string
	Filters     []map[string]any
	OrderBy     []OrderBy
	Interval    string
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
	Start      string            `json:"start"`
	End        string            `json:"end"`
	LicenseKey string            `json:"licenseKey"`
	Dimension  *string           `json:"dimension,omitempty"`
	Metric     *string           `json:"metric,omitempty"`
	Filters    *[]map[string]any `json:"filters,omitempty"`
	OrderBy    *[]OrderBy        `json:"orderBy,omitempty"`
	Interval   *string           `json:"interval,omitempty"`
}

type Message struct {
	ID   string `json:"id"`
	Date string `json:"date"`
	Type string `json:"type"`
	Text string `json:"text"`
}

type OrderBy struct {
	Name  string `json:"name"`
	Order string `json:"order"`
}

// BaseURL Bitmovin package constant
const (
	BaseURL = "https://api.bitmovin.com"
)

func getBitmovinRequestHeaders() map[string][]string {
	return map[string][]string{
		"X-Api-Key":       {bitmovinCreds.APIKey},
		"Content-Type":    {"application/json"},
		"Accept":          {"application/json"},
		"X-Tenant-Org-Id": {bitmovinCreds.TenantOrg}}
}
