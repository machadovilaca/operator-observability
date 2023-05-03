package rules

import (
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	"github.com/machadovilaca/operator-observability/pkg/operatorrules"
)

const (
	recordingRulesPrefix = "guestbook_operator_"
	namespace            = "guestbook-operator"
)

var (
	// Add your custom recording rules here
	recordingRules = [][]operatorrules.RecordingRule{
		operatorRecordingRules,
	}

	// Add your custom alerts here
	alerts = [][]promv1.Rule{
		operatorAlerts,
	}
)

func SetupRules() *promv1.PrometheusRule {
	err := operatorrules.RegisterRecordingRules(recordingRules...)
	if err != nil {
		panic(err)
	}

	err = operatorrules.RegisterAlerts(alerts...)
	if err != nil {
		panic(err)
	}

	rules, err := operatorrules.BuildPrometheusRule(
		"guestbook-operator-prometheus-rules",
		"default",
		map[string]string{"app": "guestbook-operator"},
	)
	if err != nil {
		panic(err)
	}

	return rules
}

func ListRecordingRules() []operatorrules.RecordingRule {
	return operatorrules.ListRecordingRules()
}

func ListAlerts() []promv1.Rule {
	return operatorrules.ListAlerts()
}
