package metrics

import "fmt"

type MetricsCollector struct {
	ClusterName  string
	MetricPoints []Metrics
}

type Metrics struct {
	Measurement string
	Tags        map[string]string
	Value       map[string]interface{}
}

// AddMetricsEntry will create a basic metrics entry with all the needed fields.
// Please do not provide node/clustername within the tags; this function will manage this for you
func (m *MetricsCollector) AddMetricsEntry(measurement string, node string, tags map[string]string, value interface{}) {
	tags["node"] = node
	tags["clustername"] = m.ClusterName

	m.MetricPoints = append(m.MetricPoints, Metrics{
		Measurement: measurement,
		Tags:        tags,
		Value: map[string]interface{}{
			"value": value,
		},
	})
}

func (m Metrics) String() string {
	return fmt.Sprintf("adding point to measurement=%s with tags=%+v and value=%v", m.Measurement, m.Tags, m.Value["value"])
}
