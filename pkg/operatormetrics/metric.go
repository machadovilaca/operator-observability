package operatormetrics

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

// Metric is a struct that contains the information needed to register a metric.
type Metric struct {
	Name string
	Help string
	Type MetricType

	ConstLabels map[string]string
	Labels      []string

	StabilityLevel StabilityLevel
}

// MetricType is a string that represents the type of metric.
type MetricType string

// StabilityLevel is a string that represents the stability level of a metric.
type StabilityLevel string

const (
	Counter   MetricType = "Counter"
	Gauge     MetricType = "Gauge"
	Histogram MetricType = "Histogram"
	Summary   MetricType = "Summary"

	Alpha StabilityLevel = "Alpha"
	Beta  StabilityLevel = "Beta"
	GA    StabilityLevel = "GA"
)

// GetCounterMetric returns a prometheus.Counter metric.
func GetCounterMetric(name string) prometheus.Counter {
	return getMetric(name).(prometheus.Counter)
}

// GetCounterMetricWithLabels returns a prometheus.Counter metric with labels.
func GetCounterMetricWithLabels(name string, labels ...string) prometheus.Counter {
	m := getMetric(name).(*prometheus.CounterVec)
	return m.WithLabelValues(labels...)
}

// GetGaugeMetric returns a prometheus.Gauge metric.
func GetGaugeMetric(name string) prometheus.Gauge {
	return getMetric(name).(prometheus.Gauge)
}

// GetGaugeMetricWithLabels returns a prometheus.Gauge metric with labels.
func GetGaugeMetricWithLabels(name string, labels ...string) prometheus.Gauge {
	m := getMetric(name).(*prometheus.GaugeVec)
	return m.WithLabelValues(labels...)
}

// GetHistogramMetric returns a prometheus.Histogram metric.
func GetHistogramMetric(name string) prometheus.Histogram {
	return getMetric(name).(prometheus.Histogram)
}

// GetHistogramMetricWithLabels returns a prometheus.Histogram metric with labels.
func GetHistogramMetricWithLabels(name string, labels ...string) prometheus.Observer {
	m := getMetric(name).(*prometheus.HistogramVec)
	return m.WithLabelValues(labels...)
}

// GetSummaryMetric returns a prometheus.Summary metric.
func GetSummaryMetric(name string) prometheus.Summary {
	return getMetric(name).(prometheus.Summary)
}

// GetSummaryMetricWithLabels returns a prometheus.Summary metric with labels.
func GetSummaryMetricWithLabels(name string, labels ...string) prometheus.Observer {
	m := getMetric(name).(*prometheus.SummaryVec)
	return m.WithLabelValues(labels...)
}

func getMetric(name string) prometheus.Collector {
	metric, ok := operatorRegistry.registeredMetrics[name]
	if !ok {
		panic(fmt.Sprintf("metric %s does not exist", name))
	}
	return *metric.collector
}

func registerMetric(metric Metric) (prometheus.Collector, error) {
	switch metric.Type {
	case Counter:
		if metric.Labels != nil {
			col := prometheus.NewCounterVec(prometheus.CounterOpts(getOpts(metric)), metric.Labels)
			return col, metrics.Registry.Register(col)
		}

		col := prometheus.NewCounter(prometheus.CounterOpts(getOpts(metric)))
		return col, metrics.Registry.Register(col)

	case Gauge:
		if metric.Labels != nil {
			col := prometheus.NewGaugeVec(prometheus.GaugeOpts(getOpts(metric)), metric.Labels)
			return col, metrics.Registry.Register(col)
		}

		col := prometheus.NewGauge(prometheus.GaugeOpts(getOpts(metric)))
		return col, metrics.Registry.Register(col)

	case Histogram:
		if metric.Labels != nil {
			col := prometheus.NewHistogramVec(getHistogramOpts(metric), metric.Labels)
			return col, metrics.Registry.Register(col)
		}

		col := prometheus.NewHistogram(getHistogramOpts(metric))
		return col, metrics.Registry.Register(col)

	case Summary:
		if metric.Labels != nil {
			col := prometheus.NewSummaryVec(getSummaryOpts(metric), metric.Labels)
			return col, metrics.Registry.Register(col)
		}

		col := prometheus.NewSummary(getSummaryOpts(metric))
		return col, metrics.Registry.Register(col)
	}

	return nil, fmt.Errorf("unknown metric type %s", metric.Type)
}

func getOpts(metric Metric) prometheus.Opts {
	return prometheus.Opts{
		Name:        metric.Name,
		Help:        metric.Help,
		ConstLabels: metric.ConstLabels,
	}
}

func getHistogramOpts(metric Metric) prometheus.HistogramOpts {
	return prometheus.HistogramOpts{
		Name:        metric.Name,
		Help:        metric.Help,
		ConstLabels: metric.ConstLabels,
		Buckets:     prometheus.DefBuckets, // FIXME: make configurable
	}
}

func getSummaryOpts(metric Metric) prometheus.SummaryOpts {
	return prometheus.SummaryOpts{
		Name:        metric.Name,
		Help:        metric.Help,
		ConstLabels: metric.ConstLabels,
		Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, // FIXME: make configurable
		MaxAge:      10 * 60,                                                // FIXME: make configurable
		AgeBuckets:  5,                                                      // FIXME: make configurable
		BufCap:      1000,                                                   // FIXME: make configurable
	}
}
