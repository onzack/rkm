package metrics

import "fmt"

// Collector struct: contains field for ClusterName and MetricsPoints
type Collector struct {
	ClusterName  string
	MetricPoints []Metrics
}

// Metrics struct: contains field for Measurement, Tags and
type Metrics struct {
	Measurement string
	Tags        map[string]string
	Value       map[string]interface{}
}

// AddNodesMetricsEntry will create a basic metrics entry with all the needed fields.
// Please do not provide node/clustername within the tags; this function will manage this for you
func (m *Collector) AddNodesMetricsEntry(measurement string, node string, tags map[string]string, value interface{}) {
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

// AddSimpleMetricsEntry will create a basic metrics entry with all the needed fields.
// Please do not provide node/clustername within the tags; this function will manage this for you
func (m *Collector) AddSimpleMetricsEntry(measurement string, value interface{}) {

	tags := map[string]string{
		"clustername": m.ClusterName,
	}

	m.MetricPoints = append(m.MetricPoints, Metrics{
		Measurement: measurement,
		Tags:        tags,
		Value: map[string]interface{}{
			"value": value,
		},
	})
}

// AddComponentsMetricsEntry will create a basic metrics entry with all the needed fields.
// Please do not provide node/clustername within the tags; this function will manage this for you
func (m *Collector) AddComponentsMetricsEntry(measurement string, tags map[string]string, value interface{}) {
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
