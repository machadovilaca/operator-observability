package operatorrules

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/util/intstr"

	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	"github.com/machadovilaca/operator-observability/pkg/operatormetrics"
)

var _ = Describe("Validators", func() {
	Context("RecordingRule Validation", func() {
		It("should validate recording rule with valid input", func() {
			recordingRule := &RecordingRule{
				MetricsOpts: operatormetrics.MetricOpts{
					Name: "ExampleRecordingRule",
				},
				Expr: intstr.FromString("sum(rate(http_requests_total[5m]))"),
			}
			err := recordingRuleValidator(recordingRule)
			Expect(err).To(BeNil())
		})

		It("should return error if recording rule name is empty", func() {
			recordingRule := &RecordingRule{
				MetricsOpts: operatormetrics.MetricOpts{},
				Expr:        intstr.FromString("sum(rate(http_requests_total[5m]))"),
			}
			err := recordingRuleValidator(recordingRule)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("recording rule must have a name"))
		})

		It("should return error if recording rule expression is empty", func() {
			recordingRule := &RecordingRule{
				MetricsOpts: operatormetrics.MetricOpts{
					Name: "ExampleRecordingRule",
				},
				Expr: intstr.FromString(""),
			}
			err := recordingRuleValidator(recordingRule)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("recording rule must have an expression"))
		})
	})

	Context("Alert Validation", func() {
		It("should validate alert with valid input", func() {
			alert := &promv1.Rule{
				Alert: "ExampleAlert",
				Expr:  intstr.FromString("sum(rate(http_requests_total[5m]))"),
				Labels: map[string]string{
					"severity": "critical",
				},
				Annotations: map[string]string{
					"summary":     "Example summary",
					"description": "Example description",
				},
			}
			err := alertValidator(alert)
			Expect(err).To(BeNil())
		})

		It("should return error if alert name is not in PascalCase format", func() {
			alert := &promv1.Rule{
				Alert: "example_alert",
				Expr:  intstr.FromString("sum(rate(http_requests_total[5m]))"),
				Labels: map[string]string{
					"severity": "critical",
				},
				Annotations: map[string]string{
					"summary":     "Example summary",
					"description": "Example description",
				},
			}
			err := alertValidator(alert)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("alert must have a name in PascalCase format"))
		})

		It("should return error if alert expression is empty", func() {
			alert := &promv1.Rule{
				Alert: "ExampleAlert",
				Expr:  intstr.FromString(""),
				Labels: map[string]string{
					"severity": "critical",
				},
				Annotations: map[string]string{
					"summary":     "Example summary",
					"description": "Example description",
				},
			}
			err := alertValidator(alert)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("alert must have an expression"))
		})

		It("should return error if severity label is missing or invalid", func() {
			alert := &promv1.Rule{
				Alert: "ExampleAlert",
				Expr:  intstr.FromString("sum(rate(http_requests_total[5m]))"),
				Labels: map[string]string{
					"severity": "invalid_severity",
				},
				Annotations: map[string]string{
					"summary":     "Example summary",
					"description": "Example description",
				},
			}
			err := alertValidator(alert)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("alert must have a severity label with value critical, warning, or info"))
		})

		It("should return error if summary annotation is missing", func() {
			alert := &promv1.Rule{
				Alert: "ExampleAlert",
				Expr:  intstr.FromString("sum(rate(http_requests_total[5m]))"),
				Labels: map[string]string{
					"severity": "critical",
				},
				Annotations: map[string]string{
					"description": "Example description",
				},
			}
			err := alertValidator(alert)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("alert must have a summary annotation"))
		})

		It("should return error if description annotation is missing", func() {
			alert := &promv1.Rule{
				Alert: "ExampleAlert",
				Expr:  intstr.FromString("sum(rate(http_requests_total[5m]))"),
				Labels: map[string]string{
					"severity": "critical",
				},
				Annotations: map[string]string{
					"summary": "Example summary",
				},
			}
			err := alertValidator(alert)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("alert must have a description annotation"))
		})
	})
})
