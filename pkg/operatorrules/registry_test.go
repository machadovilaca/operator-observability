package operatorrules

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/util/intstr"

	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	"github.com/machadovilaca/operator-observability/pkg/operatormetrics"
)

var _ = Describe("OperatorRules", func() {
	BeforeEach(func() {
		operatorRegistry = newRegistry()
	})

	Context("RecordingRule Registration", func() {
		It("should register recording rules without error", func() {
			recordingRules := []RecordingRule{
				{
					MetricsOpts: operatormetrics.MetricOpts{Name: "ExampleRecordingRule1"},
					Expr:        intstr.FromString("sum(rate(http_requests_total[5m]))"),
				},
				{
					MetricsOpts: operatormetrics.MetricOpts{Name: "ExampleRecordingRule2"},
					Expr:        intstr.FromString("sum(rate(http_requests_total[5m]))"),
				},
			}

			err := RegisterRecordingRules(recordingRules)
			Expect(err).To(BeNil())

			registeredRules := ListRecordingRules()
			Expect(registeredRules).To(ConsistOf(recordingRules))
		})

		It("should replace recording rule with the same name and expression", func() {
			recordingRules := []RecordingRule{
				{
					MetricsOpts: operatormetrics.MetricOpts{Name: "ExampleRecordingRule1"},
					Expr:        intstr.FromString("sum(rate(http_requests_total[5m]))"),
				},
				{
					MetricsOpts: operatormetrics.MetricOpts{Name: "ExampleRecordingRule1"},
					Expr:        intstr.FromString("sum(rate(http_requests_total[5m]))"),
				},
			}

			err := RegisterRecordingRules(recordingRules)
			Expect(err).To(BeNil())

			registeredRules := ListRecordingRules()
			Expect(registeredRules).To(HaveLen(1))
			Expect(registeredRules[0].Expr.String()).To(Equal("sum(rate(http_requests_total[5m]))"))
		})

		It("should create 2 recording rules when registered with the same name but different expressions", func() {
			recordingRules := []RecordingRule{
				{
					MetricsOpts: operatormetrics.MetricOpts{Name: "ExampleRecordingRule1"},
					Expr:        intstr.FromString("sum(rate(http_requests_total[5m]))"),
				},
			}

			err := RegisterRecordingRules(recordingRules)
			Expect(err).To(BeNil())

			recordingRules = []RecordingRule{
				{
					MetricsOpts: operatormetrics.MetricOpts{Name: "ExampleRecordingRule1"},
					Expr:        intstr.FromString("sum(rate(http_requests_total[10m]))"),
				},
			}

			err = RegisterRecordingRules(recordingRules)
			Expect(err).To(BeNil())

			registeredRules := ListRecordingRules()
			Expect(registeredRules).To(HaveLen(2))
			Expect(registeredRules[0].Expr.String()).To(Equal("sum(rate(http_requests_total[10m]))"))
			Expect(registeredRules[1].Expr.String()).To(Equal("sum(rate(http_requests_total[5m]))"))
		})
	})

	Context("Alert Registration", func() {
		It("should register alerts without error", func() {
			alerts := []promv1.Rule{
				{
					Alert: "ExampleAlert1",
					Expr:  intstr.FromString("sum(rate(http_requests_total[1m])) > 100"),
					Labels: map[string]string{
						"severity": "critical",
					},
					Annotations: map[string]string{
						"summary":     "High request rate",
						"description": "The request rate is too high.",
					},
				},
				{
					Alert: "ExampleAlert2",
					Expr:  intstr.FromString("sum(rate(http_requests_total[5m])) > 100"),
					Labels: map[string]string{
						"severity": "warning",
					},
					Annotations: map[string]string{
						"summary":     "Moderate request rate",
						"description": "The request rate is moderately high.",
					},
				},
			}

			err := RegisterAlerts(alerts)
			Expect(err).To(BeNil())

			registeredAlerts := ListAlerts()
			Expect(registeredAlerts).To(ConsistOf(alerts))
		})

		It("should replace alerts with the same name in the same RegisterAlerts call", func() {
			alerts := []promv1.Rule{
				{
					Alert: "ExampleAlert1",
					Expr:  intstr.FromString("sum(rate(http_requests_total[1m])) > 100"),
					Labels: map[string]string{
						"severity": "critical",
					},
					Annotations: map[string]string{
						"summary":     "High request rate",
						"description": "The request rate is too high.",
					},
				},
				{
					Alert: "ExampleAlert1",
					Expr:  intstr.FromString("sum(rate(http_requests_total[1m])) > 200"),
					Labels: map[string]string{
						"severity": "critical",
					},
					Annotations: map[string]string{
						"summary":     "High request rate",
						"description": "The request rate is too high.",
					},
				},
			}

			err := RegisterAlerts(alerts)
			Expect(err).To(BeNil())

			registeredAlerts := ListAlerts()
			Expect(registeredAlerts).To(HaveLen(1))
			Expect(registeredAlerts[0].Expr.String()).To(Equal("sum(rate(http_requests_total[1m])) > 200"))
		})

		It("should replace alerts with the same name in different RegisterAlerts calls", func() {
			alerts := []promv1.Rule{
				{
					Alert: "ExampleAlert1",
					Expr:  intstr.FromString("sum(rate(http_requests_total[1m])) > 100"),
					Labels: map[string]string{
						"severity": "critical",
					},
					Annotations: map[string]string{
						"summary":     "High request rate",
						"description": "The request rate is too high.",
					},
				},
			}

			err := RegisterAlerts(alerts)
			Expect(err).To(BeNil())

			alerts = []promv1.Rule{
				{
					Alert: "ExampleAlert1",
					Expr:  intstr.FromString("sum(rate(http_requests_total[1m])) > 200"),
					Labels: map[string]string{
						"severity": "critical",
					},
					Annotations: map[string]string{
						"summary":     "High request rate",
						"description": "The request rate is too high.",
					},
				},
			}

			err = RegisterAlerts(alerts)
			Expect(err).To(BeNil())

			registeredAlerts := ListAlerts()
			Expect(registeredAlerts).To(HaveLen(1))
			Expect(registeredAlerts[0].Expr.String()).To(Equal("sum(rate(http_requests_total[1m])) > 200"))
		})
	})

	Context("Clean Registry", func() {
		It("should clean registry without error", func() {
			recordingRules := []RecordingRule{
				{
					MetricsOpts: operatormetrics.MetricOpts{Name: "ExampleRecordingRule1"},
					Expr:        intstr.FromString("sum(rate(http_requests_total[5m]))"),
				},
			}

			alerts := []promv1.Rule{
				{
					Alert: "ExampleAlert1",
					Expr:  intstr.FromString("sum(rate(http_requests_total[1m])) > 100"),
				},
			}

			err := RegisterRecordingRules(recordingRules)
			Expect(err).To(BeNil())
			registeredRules := ListRecordingRules()
			Expect(registeredRules).To(ConsistOf(recordingRules))

			err = RegisterAlerts(alerts)
			Expect(err).To(BeNil())
			registeredAlerts := ListAlerts()
			Expect(registeredAlerts).To(ConsistOf(alerts))

			err = CleanRegistry()
			Expect(err).To(BeNil())

			registeredRules = ListRecordingRules()
			Expect(registeredRules).To(BeEmpty())
			registeredAlerts = ListAlerts()
			Expect(registeredAlerts).To(BeEmpty())
		})
	})
})
