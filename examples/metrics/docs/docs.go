package main

import (
	"fmt"

	"github.com/machadovilaca/operator-observability/examples/metrics"
	"github.com/machadovilaca/operator-observability/pkg/docs"
)

func main() {
	docsString := docs.BuildDocs(metrics.ListMetrics())
	fmt.Println(docsString)
}
