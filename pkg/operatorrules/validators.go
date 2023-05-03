package operatorrules

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/grafana/regexp"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

var (
	recordingRuleValidator = defaultRecordingRuleValidation
	alertValidator         = defaultAlertValidation
)

// SetRecordingRuleValidator sets the validator for recording rules.
func SetRecordingRuleValidator(validator func(recordingRule *RecordingRule) error) {
	recordingRuleValidator = validator
}

// SetAlertValidator sets the validator for alerts.
func SetAlertValidator(validator func(alert *promv1.Rule) error) {
	alertValidator = validator
}

func defaultRecordingRuleValidation(recordingRule *RecordingRule) error {
	if recordingRule.MetricsOpts.Name == "" {
		return fmt.Errorf("recording rule must have a name")
	}

	if recordingRule.Expr == intstr.FromString("") {
		return fmt.Errorf("recording rule must have an expression")
	}

	return nil
}

// based on https://sdk.operatorframework.io/docs/best-practices/observability-best-practices/#alerts-style-guide
func defaultAlertValidation(alert *promv1.Rule) error {
	if alert.Alert == "" || !isPascalCase(alert.Alert) {
		return fmt.Errorf("alert must have a name in PascalCase format")
	}

	if alert.Expr == intstr.FromString("") {
		return fmt.Errorf("alert must have an expression")
	}

	// Alerts MUST include a severity label indicating the alert’s urgency.
	if err := validateLabels(alert); err != nil {
		return err
	}

	// Alerts MUST include summary and description annotations.
	err := validateAnnotations(alert)

	return err
}

func isPascalCase(s string) bool {
	pascalCasePattern := `^[A-Z][a-z]*(?:[A-Z][a-z]*)*$`
	pascalCaseRegex := regexp.MustCompile(pascalCasePattern)
	return pascalCaseRegex.MatchString(s)
}

func validateLabels(alert *promv1.Rule) error {
	severity := alert.Labels["severity"]

	if severity == "" {
		return fmt.Errorf("alert must have a severity label")
	}

	if severity != "critical" && severity != "warning" && severity != "info" {
		return fmt.Errorf("alert severity must be one of critical, warning, info")
	}

	return nil
}

func validateAnnotations(alert *promv1.Rule) error {
	summary := alert.Annotations["summary"]
	if summary == "" {
		return fmt.Errorf("alert must have a summary annotation")
	}

	description := alert.Annotations["description"]
	if description == "" {
		return fmt.Errorf("alert must have a description annotation")
	}

	return nil
}
