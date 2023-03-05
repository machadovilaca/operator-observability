# Kubernetes Operator Observability Utilities

This repository contains a set of utilities for Kubernetes Operators to help
with observability.

The goal is to help developers of Kubernetes Operators follow the
[Operator SDK Observability Best Practices](https://sdk.operatorframework.io/docs/best-practices/observability-best-practices/)
while instrumenting their operators and avoiding common pitfalls and mistakes.

## Design

### Metrics

Operator developers can make use of the utilities provided here to uniformize
the way metrics are registered and their values set. Operator developers can use
the utilities provided here to uniformize how metrics are registered and their
values set. In many projects, developers register and set metrics in multiple
ways and places. This makes it hard to have a global view of the existing
metrics, their values, and how they are set.

#### Usage

Define different scopes for different metrics. Some metrics such as
`...reconcile_count` are related to the operator workload itself, while others
like `...out_of_band_modifications_count` are related to the custom resources
managed by the operator. These metrics can be grouped in different files so that
we have a clear separation of concerns.

The file holding the metrics related to the operator workload might look like
this:

```go
// metrics/operator_metrics.go

var operatorMetrics = []operatormetrics.Metric{
	{
		Name: reconcileCount,
		Help: "Number of times the operator has executed the reconcile loop",
		Type: operatormetrics.Counter,
		ConstLabels: map[string]string{
			"controller": "my_controller",
		},
		StabilityLevel: operatormetrics.GA,
	},
	...
}

func IncrementReconcileCountMetric() {
    m := operatormetrics.GetCounterMetric(reconcileCount)
    m.Inc()
}
```

All metrics would be registered in the `init()` function of a "central" file,
such as:

```go
// metrics/metrics.go

func init() {
	err := operatormetrics.RegisterMetrics(operatorMetrics, crMetrics, ...)
...
```

And the operator developer would use the `IncrementReconcileCountMetric()` to
increment the `...reconcile_count` metric in the reconcile loop. Remember that
for metrics that require more logic to set their values, we should still make an
effort to avoid adding monitoring logic code to the business logic of the
operator.

#### Collectors

If at any time you need to pull information from other systems, such as the
Kubernetes or Cloud Provider APIs, you can create a custom collector. The
collector should follow all the same principles as described above for the
metrics. It should also be registered in the `init()` function with:

```go
err = operatormetrics.RegisterCollector(customResourceCollector, ...)
...
```

The biggest difference from metrics is that the collectors we create will have a
callback function that will be called when the metrics are collected. This
callback function should be used to pull information from other systems and set
the values of the metrics.

```go
// metrics/custom_resource_collector.go

...
func SetupCustomResourceCollector(k8sClient *kubernetes.Clientset) {
	collectorK8sClient = k8sClient
}

var customResourceCollector = operatormetrics.Collector{
	Metrics: []operatormetrics.Metric{
		{
			Name:           customResourceCount,
			...
		},
	},
	CollectCallback: customResourceCollectorCallback,
}

func customResourceCollectorCallback() []operatormetrics.CollectionResult {
	result := unstructured.UnstructuredList{}
	err := collectorK8sClient.RESTClient().Get().Resource("MyCR").Do(context.Background()).Into(&result)
	...

	return []operatormetrics.CollectionResult{
		{
			Name:   customResourceCount,
			Labels: []string{"default"},
			Value: float64(len(result.Items)),
		},
	}
}
```

#### Documentation

Having all metrics in one place makes it easy to document them and track the
changes. The documentation of the metrics can be generated from the code using
docs utility. The utility will generate a string with the documentation of the
metrics that you can later print or save to a file. The documentation will
include the name, help, type, labels, and stability level. It includes a default
header, metric, and footer template that you can customize.

```go
func main() {
	docs.SetHeader("# My Custom Operator Metrics\n\n")
	docsString := docs.BuildDocs(metrics.ListMetrics())
	fmt.Println(docsString)
}
```

## Next Steps

- Propose design changes to the Operator SDK examples

- Also add utils for Alerts, Events, and E2E tests

- Build a Kubebuilder/Operator SDK plugin to allow developers to effortlessly
add observability to their operators

- Add Kubebuilder/Operator SDK command line instructions to generate code for
new Metrics, Alerts, and Events
