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

		It("should return an error for invalid recording rules", func() {
			invalidRecordingRules := []RecordingRule{
				{
					Expr: intstr.FromString("sum(rate(http_requests_total[5m]))"),
				},
			}

			err := RegisterRecordingRules(invalidRecordingRules)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("invalid recording rule"))
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

		It("should return an error for invalid alerts", func() {
			invalidAlerts := []promv1.Rule{
				{
					Expr: intstr.FromString("sum(rate(http_requests_total[1m])) > 100"),
					Labels: map[string]string{
						"severity": "critical",
					},
					Annotations: map[string]string{
						"summary":     "High request rate",
						"description": "The request rate is too high.",
					},
				},
			}

			err := RegisterAlerts(invalidAlerts)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("invalid alert"))
		})
	})
})
