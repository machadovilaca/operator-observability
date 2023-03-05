package operatormetrics

import (
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

var operatorRegistry = newRegistry()

type operatorRegisterer struct {
	registeredMetrics    map[string]registeredMetric
	registeredCollectors map[string]registeredCollector
}

type registeredMetric struct {
	collector *prometheus.Collector
	metric    Metric
}

type registeredCollector struct {
	desc   *prometheus.Desc
	metric Metric
}

func newRegistry() operatorRegisterer {
	return operatorRegisterer{
		registeredMetrics:    map[string]registeredMetric{},
		registeredCollectors: map[string]registeredCollector{},
	}
}

// RegisterMetrics registers the metrics with the Prometheus registry.
func RegisterMetrics(allMetrics ...[]Metric) error {
	for _, metricList := range allMetrics {
		for _, metric := range metricList {
			v, err := registerMetric(metric)
			if err != nil {
				return err
			}

			operatorRegistry.registeredMetrics[metric.Name] = registeredMetric{
				collector: &v,
				metric:    metric,
			}
		}
	}

	return nil
}

// RegisterCollector registers the collector with the Prometheus registry.
func RegisterCollector(collectors ...Collector) error {
	for _, collector := range collectors {
		err := metrics.Registry.Register(collector)
		if err != nil {
			return err
		}

		for _, metric := range collector.Metrics {
			operatorRegistry.registeredCollectors[metric.Name] = registeredCollector{
				desc:   prometheus.NewDesc(metric.Name, metric.Help, metric.Labels, metric.ConstLabels),
				metric: metric,
			}
		}
	}

	return nil
}

// ListMetrics returns a list of all registered metrics.
func ListMetrics() []Metric {
	var result []Metric

	for _, rm := range operatorRegistry.registeredMetrics {
		result = append(result, rm.metric)
	}

	for _, rc := range operatorRegistry.registeredCollectors {
		result = append(result, rc.metric)
	}

	return result
}
