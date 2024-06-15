package bitmovin

type BitmovinCredentials struct {
	apiKey			string
	licenseKey		string
	tenantOrg		string
}

type BitmovinQueryParams struct {
	Name        string
	URI         string
	NRMetric    string
	BMDimension string
	Metric      string
	Filters     []map[string]any
	GroupBy		[]string
	OrderBy     []BitmovinOrderBy
	Interval    string
	Percentile	*int64
}

type BitmovinColumn struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type BitmovinRowData []any

type BitmovinResult struct {
	RowCount     int       `json:"rowCount"`
	Rows         []BitmovinRowData `json:"rows"`
	ColumnLabels []BitmovinColumn  `json:"columnLabels"`
	JobID        string    `json:"jobId"`
	TaskCount    int       `json:"taskCount"`
}

type BitmovinMessage struct {
	ID   string `json:"id"`
	Date string `json:"date"`
	Type string `json:"type"`
	Text string `json:"text"`
}

type BitmovinData struct {
	Result   BitmovinResult    `json:"result"`
	Messages []BitmovinMessage `json:"messages"`
}

type BitmovinResponse struct {
	RequestID string `json:"requestId"`
	Status    string `json:"status"`
	Data      BitmovinData   `json:"data"`
	Metric    string `json:"metric"`
}

type BitmovinOrderBy struct {
	Name  string `json:"name"`
	Order string `json:"order"`
}

type BitmovinRequestBody struct {
	Start      string            	`json:"start"`
	End        string            	`json:"end"`
	LicenseKey string            	`json:"licenseKey"`
	Dimension  *string           	`json:"dimension,omitempty"`
	Metric     *string           	`json:"metric,omitempty"`
	Filters    *[]map[string]any 	`json:"filters,omitempty"`
	GroupBy	   *[]string		 	`json:"groupBy,omitempty"`
	OrderBy    *[]BitmovinOrderBy  	`json:"orderBy,omitempty"`
	Interval   *string           	`json:"interval,omitempty"`
	Percentile *int64				`json:"percentile,omitempty"`
}

type BitmovinFilter struct {
	Operator string
	Value any
}

type BitmovinQuery struct {
	Type 		string 						`json:"type"`
	Metric 		string 						`json:"metric"`
	Interval 	*string 					`json:"interval,omitempty"`
	Dimensions 	*[]string 					`json:"dimensions,omitempty"`
	Filters 	*map[string]BitmovinFilter 	`json:"filters,omitempty"`
	GroupBy 	*[]string 					`json:"groupBy,omitempty"`
	OrderBy    	*[]BitmovinOrderBy  		`json:"orderBy,omitempty"`
	Percentile	*int64						`json:"percentile,omitempty"`
}
