package metrics

import "encoding/json"

type GetHostMetricsOutput struct {
	Metrics []MetricValue `json:"metrics"`
}

func (o *GetHostMetricsOutput) String() string {
	bytea, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

type MetricValue struct {
	HostID string  `json:"hostId,omitempty"`
	Name   string  `json:"name,omitempty"`
	Time   int64   `json:"time"`
	Value  float64 `json:"value"`
}
