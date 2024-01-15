package operatormetrics

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"
)

var _ = Describe("Collector", func() {
	var (
		testCounterOpts = MetricOpts{
			Name: "collector_test_counter",
			Help: "A test counter",
		}
		testGaugeOpts = MetricOpts{
			Name: "collector_test_gauge",
			Help: "A test gauge",
		}
		testCounterOpts2 = MetricOpts{
			Name: "collector_test_counter_2",
			Help: "A test counter",
		}
	)

	Describe("Collect", func() {
		BeforeEach(func() {
			CleanRegistry()
		})

		It("should collect metrics from registered collectors", func() {
			counter := NewCounter(testCounterOpts)
			gauge := NewGauge(testGaugeOpts)

			collector := Collector{
				Metrics: []Metric{counter, gauge},
				CollectCallback: func() []CollectorResult {
					return []CollectorResult{
						{Metric: counter, Labels: nil, Value: 5},
						{Metric: gauge, Labels: nil, Value: 10},
					}
				},
			}

			err := RegisterCollector(collector)
			Expect(err).NotTo(HaveOccurred())

			ch := make(chan prometheus.Metric, 2)
			go collector.Collect(ch)

			metricCounter := <-ch
			metricGauge := <-ch

			Expect(metricCounter.Desc().String()).To(ContainSubstring(testCounterOpts.Name))
			Expect(metricGauge.Desc().String()).To(ContainSubstring(testGaugeOpts.Name))
		})

		It("should skip unregistered collectors", func() {
			counter := NewCounter(testCounterOpts2)

			collector := Collector{
				Metrics: []Metric{counter},
				CollectCallback: func() []CollectorResult {
					return []CollectorResult{
						{Metric: counter, Labels: nil, Value: 5},
					}
				},
			}

			ch := make(chan prometheus.Metric, 1)
			go collector.Collect(ch)

			Expect(ch).NotTo(Receive())
		})
	})
})
