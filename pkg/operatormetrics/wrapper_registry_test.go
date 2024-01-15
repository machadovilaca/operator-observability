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
			CleanRegistry()
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

		It("should return an error when registering a duplicate metric", func() {
			counter := NewCounter(testCounterOpts)

			err := RegisterMetrics([]Metric{counter})
			Expect(err).NotTo(HaveOccurred())

			err = RegisterMetrics([]Metric{counter})
			Expect(err).To(HaveOccurred())

			Expect(operatorRegistry.registeredMetrics).To(HaveLen(1))
		})
	})

	Describe("ListMetrics", func() {
		BeforeEach(func() {
			CleanRegistry()
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
})
