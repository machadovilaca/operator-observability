package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/machadovilaca/operator-observability/examples/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	metrics.SetupMetrics()

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	fmt.Println("Server started on port 2112")
	for {
		metrics.IncrementReconcileCountMetric()
		metrics.IncrementReconcileActionMetric("sleep")
		time.Sleep(10 * time.Second)
	}
}
