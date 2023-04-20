package docs

import (
	"bytes"
	"log"
	"sort"
	"strings"
	"text/template"

	"github.com/machadovilaca/operator-observability/pkg/operatormetrics"
)

const defaultTemplate = `# Operator Metrics

{{- range . }}

### {{.Name}}
{{.Help}}. Type: {{.Type}}.
{{- end }}

## Developing new metrics

All metrics documented here are auto-generated and reflect exactly what is being
exposed. After developing new metrics or changing old ones please regenerate
this document.
`

type metricDocs struct {
	Name        string
	Help        string
	Type        string
	ExtraFields map[string]string
}

// BuildDocsWithCustomTemplate returns a string with the documentation for the
// given metrics, using the given template.
func BuildDocsWithCustomTemplate(metrics []operatormetrics.Metric, tplString string) string {
	tpl, err := template.New("metrics").Parse(tplString)
	if err != nil {
		log.Fatalln(err)
	}

	buf := bytes.NewBufferString("")
	err = tpl.Execute(buf, buildMetricsDocs(metrics))
	if err != nil {
		log.Fatalln(err)
	}

	return buf.String()
}

// BuildDocs returns a string with the documentation for the given metrics.
func BuildDocs(metrics []operatormetrics.Metric) string {
	return BuildDocsWithCustomTemplate(metrics, defaultTemplate)
}

func buildMetricsDocs(metrics []operatormetrics.Metric) []metricDocs {
	metricsDocs := make([]metricDocs, len(metrics))
	for i, metric := range metrics {
		metricOpts := metric.GetOpts()
		metricsDocs[i] = metricDocs{
			Name:        metricOpts.Name,
			Help:        metricOpts.Help,
			Type:        getAndConvertMetricType(metric),
			ExtraFields: metricOpts.ExtraFields,
		}
	}
	sortMetricsDocs(metricsDocs)

	return metricsDocs
}

func sortMetricsDocs(metricsDocs []metricDocs) {
	sort.Slice(metricsDocs, func(i, j int) bool {
		return metricsDocs[i].Name < metricsDocs[j].Name
	})
}

func getAndConvertMetricType(metric operatormetrics.Metric) string {
	return strings.ReplaceAll(string(metric.GetType()), "Vec", "")
}
