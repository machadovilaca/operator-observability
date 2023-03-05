package main

import (
	"fmt"

	"github.com/machadovilaca/operator-observability/pkg/docs"

	"github.com/machadovilaca/operator-observability/examples/metrics"
)

func main() {
	docs.SetHeader("# My Custom Operator Metrics\n\n")
	docsString := docs.BuildDocs(metrics.ListMetrics())
	fmt.Println(docsString)
}
