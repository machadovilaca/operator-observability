package operatormetrics_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"

	"github.com/machadovilaca/operator-observability/pkg/operatormetrics"
)

var _ = Describe("Collector", func() {
	var (
		testCounterOpts = operatormetrics.MetricOpts{
			Name: "collector_test_counter",
			Help: "A test counter",
		}
		testGaugeOpts = operatormetrics.MetricOpts{
			Name: "collector_test_gauge",
			Help: "A test gauge",
		}
		testCounterOpts2 = operatormetrics.MetricOpts{
			Name: "collector_test_counter_2",
			Help: "A test counter",
		}
	)

	Describe("Collect", func() {
		BeforeEach(func() {
			err := operatormetrics.CleanRegistry()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should collect metrics from registered collectors", func() {
			counter := operatormetrics.NewCounter(testCounterOpts)
			gauge := operatormetrics.NewGauge(testGaugeOpts)

			collector := operatormetrics.Collector{
				Metrics: []operatormetrics.Metric{counter, gauge},
				CollectCallback: func() []operatormetrics.CollectorResult {
					return []operatormetrics.CollectorResult{
						{Metric: counter, Labels: nil, Value: 5},
						{Metric: gauge, Labels: nil, Value: 10},
					}
				},
			}

			err := operatormetrics.RegisterCollector(collector)
			Expect(err).NotTo(HaveOccurred())

			ch := make(chan prometheus.Metric, 2)
			go collector.Collect(ch)

			metricCounter := <-ch
			metricGauge := <-ch

			Expect(metricCounter.Desc().String()).To(ContainSubstring(testCounterOpts.Name))
			Expect(metricGauge.Desc().String()).To(ContainSubstring(testGaugeOpts.Name))
		})

		It("should skip unregistered collectors", func() {
			counter := operatormetrics.NewCounter(testCounterOpts2)

			collector := operatormetrics.Collector{
				Metrics: []operatormetrics.Metric{counter},
				CollectCallback: func() []operatormetrics.CollectorResult {
					return []operatormetrics.CollectorResult{
						{Metric: counter, Labels: nil, Value: 5},
					}
				},
			}

			ch := make(chan prometheus.Metric, 1)
			go collector.Collect(ch)

			Expect(ch).NotTo(Receive())
		})

		It("should collect metrics with const labels added on collection time", func() {
			counter := operatormetrics.NewCounter(testCounterOpts)

			collector := operatormetrics.Collector{
				Metrics: []operatormetrics.Metric{counter},
				CollectCallback: func() []operatormetrics.CollectorResult {
					return []operatormetrics.CollectorResult{
						{
							Metric: counter,
							Labels: nil,
							ConstLabels: map[string]string{
								"important_info": "xpto",
							},
							Value: 5,
						},
					}
				},
			}

			err := operatormetrics.RegisterCollector(collector)
			Expect(err).NotTo(HaveOccurred())

			ch := make(chan prometheus.Metric, 1)
			go collector.Collect(ch)

			metricCounter := <-ch

			Expect(metricCounter.Desc().String()).To(ContainSubstring(testCounterOpts.Name))
			Expect(metricCounter.Desc().String()).To(ContainSubstring("constLabels: {important_info=\"xpto\"}"))
		})

		It("should collect metrics with custom timestamps", func() {
			counter := operatormetrics.NewCounter(testCounterOpts)

			collector := operatormetrics.Collector{
				Metrics: []operatormetrics.Metric{counter},
				CollectCallback: func() []operatormetrics.CollectorResult {
					return []operatormetrics.CollectorResult{
						{Metric: counter, Labels: nil, Value: 5, Timestamp: time.UnixMilli(1000)},
					}
				},
			}

			err := operatormetrics.RegisterCollector(collector)
			Expect(err).NotTo(HaveOccurred())

			ch := make(chan prometheus.Metric, 1)
			go collector.Collect(ch)

			metricCounter := <-ch

			Expect(metricCounter.Desc().String()).To(ContainSubstring(testCounterOpts.Name))

			dto := &io_prometheus_client.Metric{}
			err = metricCounter.Write(dto)
			Expect(err).NotTo(HaveOccurred())

			Expect(dto.TimestampMs).To(Equal(proto.Int64(1000)))
		})
	})
})
