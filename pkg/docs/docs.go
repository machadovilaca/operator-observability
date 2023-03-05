package docs

import (
	"fmt"
	"sort"
	"strings"

	"github.com/machadovilaca/operator-observability/pkg/operatormetrics"
)

var (
	header = "# Operator Metrics\n\n"

	metricEntryTemplate = "### %s\n%s%s. Type: %s.\n\n"

	footer = `## Developing new metrics

All metrics documented here are auto-generated and reflect exactly what is being
exposed. After developing new metrics or changing old ones please regenerate
this document.
`
)

// BuildDocs returns a string with the documentation for the given metrics.
func BuildDocs(metrics []operatormetrics.Metric) string {
	sb := &strings.Builder{}

	sb.WriteString(header)

	sortMetrics(metrics)
	writeMetrics(metrics, sb)

	sb.WriteString(footer)

	return sb.String()
}

// SetHeader sets the header of the documentation.
func SetHeader(newHeader string) {
	header = newHeader
}

// SetMetricEntryTemplate sets the template used to generate the metric entry.
// The template must have 4 string format tags:
// 1. metric name
// 2. metric stability
// 3. metric help
// 4. metric type
func SetMetricEntryTemplate(newMetricEntryTemplate string) {
	metricEntryTemplate = newMetricEntryTemplate
}

// SetFooter sets the footer of the documentation.
func SetFooter(newFooter string) {
	footer = newFooter
}

func sortMetrics(metricsList []operatormetrics.Metric) {
	sort.Slice(metricsList, func(i, j int) bool {
		return metricsList[i].Name < metricsList[j].Name
	})
}

func writeMetrics(metricsList []operatormetrics.Metric, sb *strings.Builder) {
	for _, metric := range metricsList {
		writeMetric(metric, sb)
	}
}

func writeMetric(metric operatormetrics.Metric, sb *strings.Builder) {
	stability := ""
	if metric.StabilityLevel != operatormetrics.GA {
		stability = fmt.Sprintf("[%s] ", metric.StabilityLevel)
	}

	sb.WriteString(fmt.Sprintf(metricEntryTemplate, metric.Name, stability, metric.Help, metric.Type))
}
