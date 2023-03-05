package metrics

import "github.com/machadovilaca/operator-observability/pkg/operatormetrics"

const (
	reconcileCount  = metricPrefix + "reconcile_count"
	reconcileAction = metricPrefix + "reconcile_action_count"
)

var operatorMetrics = []operatormetrics.Metric{
	{
		Name: reconcileCount,
		Help: "Number of times the operator has executed the reconcile loop",
		Type: operatormetrics.Counter,
		ConstLabels: map[string]string{
			"controller": "guestbook",
		},
		StabilityLevel: operatormetrics.GA,
	},
	{
		Name:           reconcileAction,
		Help:           "Number of times the operator has executed the reconcile loop with a given action",
		Type:           operatormetrics.Counter,
		Labels:         []string{"action"},
		StabilityLevel: operatormetrics.Alpha,
	},
}

func IncrementReconcileCountMetric() {
	m := operatormetrics.GetCounterMetric(reconcileCount)
	m.Inc()
}

func IncrementReconcileActionMetric(action string) {
	m := operatormetrics.GetCounterMetricWithLabels(reconcileAction, action)
	m.Inc()
}
