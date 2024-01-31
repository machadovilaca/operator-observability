package operatormetrics

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Registry", func() {
	var (
		testCounterOpts = MetricOpts{
			Name: "test_counter",
			Help: "A test counter",
		}
		testGaugeOpts = MetricOpts{
			Name: "test_gauge",
			Help: "A test gauge",
		}
	)

	Describe("RegisterMetrics", func() {
		BeforeEach(func() {
			err := CleanRegistry()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should register metrics without error", func() {
			counter := NewCounter(testCounterOpts)
			gauge := NewGauge(testGaugeOpts)

			err := RegisterMetrics([]Metric{counter, gauge})
			Expect(err).NotTo(HaveOccurred())

			Expect(operatorRegistry.registeredMetrics).To(HaveLen(2))
			Expect(operatorRegistry.registeredMetrics).To(HaveKey(testCounterOpts.Name))
			Expect(operatorRegistry.registeredMetrics).To(HaveKey(testGaugeOpts.Name))
		})

		It("should replace metrics with the same name in different RegisterMetrics call", func() {
			counter := NewCounter(testCounterOpts)

			err := RegisterMetrics([]Metric{counter})
			Expect(err).NotTo(HaveOccurred())

			err = RegisterMetrics([]Metric{counter})
			Expect(err).NotTo(HaveOccurred())

			Expect(operatorRegistry.registeredMetrics).To(HaveLen(1))
			Expect(operatorRegistry.registeredMetrics).To(HaveKey(testCounterOpts.Name))
		})
	})

	Describe("ListMetrics", func() {
		BeforeEach(func() {
			err := CleanRegistry()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return a list of all registered metrics", func() {
			counter := NewCounter(testCounterOpts)
			gauge := NewGauge(testGaugeOpts)

			err := RegisterMetrics([]Metric{counter, gauge})
			Expect(err).NotTo(HaveOccurred())

			metrics := ListMetrics()
			Expect(metrics).To(HaveLen(2))
			Expect(metrics).To(ContainElement(counter))
			Expect(metrics).To(ContainElement(gauge))
		})
	})

	Describe("RegisterCollector", func() {
		BeforeEach(func() {
			err := CleanRegistry()
			Expect(err).NotTo(HaveOccurred())
		})

		var customResourceCollectorCallback = func() []CollectorResult {
			return []CollectorResult{}
		}

		It("should register collectors without error", func() {
			collector := Collector{
				Metrics: []Metric{
					NewCounter(testCounterOpts),
				},
				CollectCallback: customResourceCollectorCallback,
			}

			err := RegisterCollector(collector)
			Expect(err).NotTo(HaveOccurred())

			Expect(operatorRegistry.registeredCollectorMetrics).To(HaveLen(1))
			Expect(operatorRegistry.registeredCollectorMetrics).To(HaveKey(testCounterOpts.Name))
		})

		It("should replace metrics with the same name in different RegisterCollector call", func() {
			collector := Collector{
				Metrics: []Metric{
					NewCounter(testCounterOpts),
				},
				CollectCallback: customResourceCollectorCallback,
			}

			err := RegisterCollector(collector)
			Expect(err).NotTo(HaveOccurred())

			err = RegisterCollector(collector)
			Expect(err).NotTo(HaveOccurred())

			Expect(operatorRegistry.registeredCollectorMetrics).To(HaveLen(1))
			Expect(operatorRegistry.registeredCollectorMetrics).To(HaveKey(testCounterOpts.Name))
		})
	})

	Describe("CleanRegistry", func() {
		BeforeEach(func() {
			err := CleanRegistry()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should remove all metrics from the registry", func() {
			counter := NewCounter(testCounterOpts)
			gauge := NewGauge(testGaugeOpts)

			err := RegisterMetrics([]Metric{counter, gauge})
			Expect(err).NotTo(HaveOccurred())

			err = CleanRegistry()
			Expect(err).NotTo(HaveOccurred())

			Expect(operatorRegistry.registeredMetrics).To(HaveLen(0))
		})

		It("should remove all collectors from the registry", func() {
			collector := Collector{
				Metrics: []Metric{
					NewCounter(testCounterOpts),
				},
				CollectCallback: func() []CollectorResult {
					return []CollectorResult{}
				},
			}

			err := RegisterCollector(collector)
			Expect(err).NotTo(HaveOccurred())

			err = CleanRegistry()
			Expect(err).NotTo(HaveOccurred())

			Expect(operatorRegistry.registeredCollectors).To(HaveLen(0))
			Expect(operatorRegistry.registeredCollectorMetrics).To(HaveLen(0))
		})
	})
})
