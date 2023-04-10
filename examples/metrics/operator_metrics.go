package metrics

import "github.com/machadovilaca/operator-observability/pkg/operatormetrics"

var (
	operatorMetrics = []operatormetrics.Metric{
		reconcileCount,
		reconcileAction,
	}

	reconcileCount = operatormetrics.NewCounter(
		operatormetrics.MetricOpts{
			Name: metricPrefix + "reconcile_count",
			Help: "Number of times the operator has executed the reconcile loop",
			ConstLabels: map[string]string{
				"controller": "guestbook",
			},
			StabilityLevel: operatormetrics.Stable,
		},
	)

	reconcileAction = operatormetrics.NewCounterVec(
		operatormetrics.MetricOpts{
			Name:           metricPrefix + "reconcile_action_count",
			Help:           "Number of times the operator has executed the reconcile loop with a given action",
			StabilityLevel: operatormetrics.Alpha,
		},
		[]string{"action"},
	)
)

func IncrementReconcileCountMetric() {
	reconcileCount.Inc()
}

func IncrementReconcileActionMetric(action string) {
	reconcileAction.WithLabelValues(action).Inc()
}
