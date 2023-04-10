package metrics

import "github.com/machadovilaca/operator-observability/pkg/operatormetrics"

const metricPrefix = "guestbook_operator_"

func SetupMetrics() {
	err := operatormetrics.RegisterMetrics(operatorMetrics)
	if err != nil {
		panic(err)
	}

	err = operatormetrics.RegisterCollector(customResourceCollector)
	if err != nil {
		panic(err)
	}
}

// ListMetrics returns a list of all metrics exposed by the operator
func ListMetrics() []operatormetrics.Metric {
	return operatormetrics.ListMetrics()
}
