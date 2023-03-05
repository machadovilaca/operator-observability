package operatormetrics

import "github.com/prometheus/client_golang/prometheus"

// Collector registers a prometheus.Collector with a set of metrics in the
// Prometheus registry. The metrics are collected by calling the CollectCallback
// function.
type Collector struct {
	// Metrics is a list of metrics to be collected by the collector.
	Metrics []Metric

	// CollectCallback is a function that returns a list of CollectionResults.
	// The CollectionResults are used to populate the metrics in the collector.
	CollectCallback func() []CollectionResult
}

// CollectionResult is a single metric value with a set of labels.
type CollectionResult struct {
	Name   string
	Value  float64
	Labels []string
}

// Describe implements the prometheus.Collector interface.
func (c Collector) Describe(_ chan<- *prometheus.Desc) {}

// Collect implements the prometheus.Collector interface.
func (c Collector) Collect(ch chan<- prometheus.Metric) {
	collectedMetrics := c.CollectCallback()

	for _, cm := range collectedMetrics {
		rc, ok := operatorRegistry.registeredCollectors[cm.Name]
		if !ok {
			continue
		}

		promType := prometheus.UntypedValue
		switch rc.metric.Type {
		case Counter:
			promType = prometheus.CounterValue
		case Gauge:
			promType = prometheus.GaugeValue
		}

		mv, _ := prometheus.NewConstMetric(rc.desc, promType, cm.Value, cm.Labels...)
		ch <- mv
	}
}
