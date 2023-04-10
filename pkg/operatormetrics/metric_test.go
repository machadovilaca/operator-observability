package operatormetrics

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

var _ = Describe("Metrics", func() {
	var (
		testCounterOpts = MetricOpts{
			Name: "test_counter",
			Help: "A test counter",
		}
		testCounterVecOpts = MetricOpts{
			Name: "test_counter_vec",
			Help: "A test counter vec",
		}
		testGaugeOpts = MetricOpts{
			Name: "test_gauge",
			Help: "A test gauge",
		}
		testGaugeVecOpts = MetricOpts{
			Name: "test_gauge_vec",
			Help: "A test gauge vec",
		}
		testHistogramOpts = MetricOpts{
			Name: "test_histogram",
			Help: "A test histogram",
		}
		testHistogramHistogramOpts = HistogramOpts{
			Buckets: prometheus.LinearBuckets(0, 10, 10),
		}
		testHistogramVecOpts = MetricOpts{
			Name: "test_histogram_vec",
			Help: "A test histogram vec",
		}
		testSummaryOpts = MetricOpts{
			Name: "test_summary",
			Help: "A test summary",
		}
		testSummarySummaryOpts = SummaryOpts{
			Objectives: map[float64]float64{0.1: 0.1, 0.2: 0.2, 0.3: 0.3, 0.4: 0.4, 0.5: 0.5},
		}
		testSummaryVecOpts = MetricOpts{
			Name: "test_summary_vec",
			Help: "A test summary vec",
		}
	)

	Describe("Metric Constructors", func() {
		It("should create a new Counter with the provided options", func() {
			counter := NewCounter(testCounterOpts)
			Expect(counter).NotTo(BeNil())
			Expect(counter.GetOpts()).To(Equal(testCounterOpts))
			Expect(counter.GetType()).To(Equal(CounterType))
		})

		It("should create a new CounterVec with the provided options and labels", func() {
			labels := []string{"label1", "label2"}
			counterVec := NewCounterVec(testCounterVecOpts, labels)
			Expect(counterVec).NotTo(BeNil())
			Expect(counterVec.GetOpts()).To(Equal(testCounterVecOpts))
			Expect(counterVec.GetType()).To(Equal(CounterVecType))
		})

		It("should create a new Gauge with the provided options", func() {
			gauge := NewGauge(testGaugeOpts)
			Expect(gauge).NotTo(BeNil())
			Expect(gauge.GetOpts()).To(Equal(testGaugeOpts))
			Expect(gauge.GetType()).To(Equal(GaugeType))
		})

		It("should create a new GaugeVec with the provided options and labels", func() {
			labels := []string{"label1", "label2"}
			gaugeVec := NewGaugeVec(testGaugeVecOpts, labels)
			Expect(gaugeVec).NotTo(BeNil())
			Expect(gaugeVec.GetOpts()).To(Equal(testGaugeVecOpts))
			Expect(gaugeVec.GetType()).To(Equal(GaugeVecType))
		})

		It("should create a new Histogram with the provided options", func() {
			histogram := NewHistogram(testHistogramOpts, testHistogramHistogramOpts)
			Expect(histogram).NotTo(BeNil())
			Expect(histogram.GetOpts()).To(Equal(testHistogramOpts))
			Expect(histogram.GetType()).To(Equal(HistogramType))
		})

		It("should create a new HistogramVec with the provided options and labels", func() {
			labels := []string{"label1", "label2"}
			histogramVec := NewHistogramVec(testHistogramVecOpts, testHistogramHistogramOpts, labels)
			Expect(histogramVec).NotTo(BeNil())
			Expect(histogramVec.GetOpts()).To(Equal(testHistogramVecOpts))
			Expect(histogramVec.GetType()).To(Equal(HistogramVecType))
		})

		It("should create a new Summary with the provided options", func() {
			summary := NewSummary(testSummaryOpts, testSummarySummaryOpts)
			Expect(summary).NotTo(BeNil())
			Expect(summary.GetOpts()).To(Equal(testSummaryOpts))
			Expect(summary.GetType()).To(Equal(SummaryType))
		})

		It("should create a new SummaryVec with the provided options and labels", func() {
			labels := []string{"label1", "label2"}
			summaryVec := NewSummaryVec(testSummaryVecOpts, testSummarySummaryOpts, labels)
			Expect(summaryVec).NotTo(BeNil())
			Expect(summaryVec.GetOpts()).To(Equal(testSummaryVecOpts))
			Expect(summaryVec.GetType()).To(Equal(SummaryVecType))
		})
	})

	Describe("Counter and CounterVec", func() {
		It("should increment the counter and counter with labels", func() {
			counter := NewCounter(testCounterOpts)
			counterVec := NewCounterVec(testCounterVecOpts, []string{"label1"})

			counter.Inc()
			counterVec.WithLabelValues("value1").Add(2)

			ch := make(chan prometheus.Metric, 2)
			counter.getCollector().Collect(ch)
			counterVec.getCollector().Collect(ch)

			metricCounter := <-ch
			metricCounterVec := <-ch

			Expect(metricCounter.Desc().String()).To(ContainSubstring(testCounterOpts.Name))
			Expect(metricCounterVec.Desc().String()).To(ContainSubstring(testCounterVecOpts.Name))

			dto := &io_prometheus_client.Metric{}

			err := metricCounter.Write(dto)
			Expect(err).NotTo(HaveOccurred())
			Expect(dto.Counter.GetValue()).To(BeEquivalentTo(1))

			err = metricCounterVec.Write(dto)
			Expect(err).NotTo(HaveOccurred())
			Expect(dto.Counter.GetValue()).To(BeEquivalentTo(2))
		})
	})

	Describe("Gauge and GaugeVec", func() {
		It("should set the gauge and gauge with labels", func() {
			gauge := NewGauge(testGaugeOpts)
			gaugeVec := NewGaugeVec(testGaugeVecOpts, []string{"label1"})

			gauge.Set(42)
			gaugeVec.WithLabelValues("value1").Set(43)

			ch := make(chan prometheus.Metric, 2)
			gauge.getCollector().Collect(ch)
			gaugeVec.getCollector().Collect(ch)

			metricGauge := <-ch
			metricGaugeVec := <-ch

			Expect(metricGauge.Desc().String()).To(ContainSubstring(testGaugeOpts.Name))
			Expect(metricGaugeVec.Desc().String()).To(ContainSubstring(testGaugeVecOpts.Name))

			dto := &io_prometheus_client.Metric{}

			err := metricGauge.Write(dto)
			Expect(err).NotTo(HaveOccurred())
			Expect(dto.Gauge.GetValue()).To(BeEquivalentTo(42))

			err = metricGaugeVec.Write(dto)
			Expect(err).NotTo(HaveOccurred())
			Expect(dto.Gauge.GetValue()).To(BeEquivalentTo(43))
		})
	})

	Describe("Histogram and HistogramVec", func() {
		It("should observe the histogram and histogram with labels", func() {
			histogram := NewHistogram(testHistogramOpts, testHistogramHistogramOpts)
			histogramVec := NewHistogramVec(testHistogramVecOpts, testHistogramHistogramOpts, []string{"label1"})

			histogram.Observe(42)
			histogramVec.WithLabelValues("value1").Observe(43)

			ch := make(chan prometheus.Metric, 2)
			histogram.getCollector().Collect(ch)
			histogramVec.getCollector().Collect(ch)

			metricHistogram := <-ch
			metricHistogramVec := <-ch

			Expect(metricHistogram.Desc().String()).To(ContainSubstring(testHistogramOpts.Name))
			Expect(metricHistogramVec.Desc().String()).To(ContainSubstring(testHistogramVecOpts.Name))

			dto := &io_prometheus_client.Metric{}

			err := metricHistogram.Write(dto)
			Expect(err).NotTo(HaveOccurred())
			Expect(dto.Histogram.GetSampleCount()).To(BeEquivalentTo(1))

			err = metricHistogramVec.Write(dto)
			Expect(err).NotTo(HaveOccurred())
			Expect(dto.Histogram.GetSampleCount()).To(BeEquivalentTo(1))
		})
	})

	Describe("Summary and SummaryVec", func() {
		It("should observe the summary and summary with labels", func() {
			summary := NewSummary(testSummaryOpts, testSummarySummaryOpts)
			summaryVec := NewSummaryVec(testSummaryVecOpts, testSummarySummaryOpts, []string{"label1"})

			summary.Observe(42)
			summaryVec.WithLabelValues("value1").Observe(43)

			ch := make(chan prometheus.Metric, 2)
			summary.getCollector().Collect(ch)
			summaryVec.getCollector().Collect(ch)

			metricSummary := <-ch
			metricSummaryVec := <-ch

			Expect(metricSummary.Desc().String()).To(ContainSubstring(testSummaryOpts.Name))
			Expect(metricSummaryVec.Desc().String()).To(ContainSubstring(testSummaryVecOpts.Name))

			dto := &io_prometheus_client.Metric{}

			err := metricSummary.Write(dto)
			Expect(err).NotTo(HaveOccurred())
			Expect(dto.Summary.GetSampleCount()).To(BeEquivalentTo(1))

			err = metricSummaryVec.Write(dto)
			Expect(err).NotTo(HaveOccurred())
			Expect(dto.Summary.GetSampleCount()).To(BeEquivalentTo(1))
		})
	})
})
