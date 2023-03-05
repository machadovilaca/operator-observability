package metrics

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"

	"github.com/machadovilaca/operator-observability/pkg/operatormetrics"
)

const customResourceCount = metricPrefix + "cr_count"

var (
	collectorK8sClient *kubernetes.Clientset
)

func SetupCustomResourceCollector(k8sClient *kubernetes.Clientset) {
	collectorK8sClient = k8sClient
}

var customResourceCollector = operatormetrics.Collector{
	Metrics: []operatormetrics.Metric{
		{
			Name:           customResourceCount,
			Help:           "Number of existing guestbook custom resources",
			Type:           operatormetrics.Gauge,
			ConstLabels:    map[string]string{"controller": "guestbook"},
			Labels:         []string{"namespace"},
			StabilityLevel: operatormetrics.Beta,
		},
	},
	CollectCallback: customResourceCollectorCallback,
}

func customResourceCollectorCallback() []operatormetrics.CollectionResult {
	result := unstructured.UnstructuredList{}
	err := collectorK8sClient.RESTClient().Get().Resource("MyCR").Do(context.Background()).Into(&result)
	if err != nil {
		return []operatormetrics.CollectionResult{}
	}

	return []operatormetrics.CollectionResult{
		{
			Name:   customResourceCount,
			Labels: []string{"default"},
			Value:  float64(len(result.Items)),
		},
	}
}
