# Kubernetes Operator Observability Utilities

This repository contains a set of utilities for Kubernetes Operators to help
with observability.

The goal is to help developers of Kubernetes Operators follow the
[Operator SDK Observability Best Practices](https://sdk.operatorframework.io/docs/best-practices/observability-best-practices/)
while instrumenting their operators and avoiding common pitfalls and mistakes.

Check the [examples](examples) folder for a complete example of how to use
these utilities.

## Design

### Metrics

Operator developers can make use of the utilities provided here to uniformize
the way metrics are registered and their values set. In many projects,
developers register and set metrics in multiple ways and places. This makes it
hard to have a global view of the existing metrics, their values, and how they
are set.

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

var (
    operatorMetrics = []operatormetrics.Metric{
        reconcileCount,
    }

    reconcileCount = operatormetrics.NewCounter(
        operatormetrics.MetricOpts{
            Name: metricPrefix + "reconcile_count",
            Help: "Number of times the operator has executed the reconcile loop",
            ConstLabels: map[string]string{
                "controller": "guestbook",
            },
            ExtraFields: map[string]string{
                "StabilityLevel": "STABLE",
            },
        },
    )
)

func IncrementReconcileCountMetric() {
    reconcileCount.Inc()
}
```

All metrics should be registered, for example, a `SetupMetrics()` function:

```go
// metrics/metrics.go

func SetupMetrics() {
    // When using controller-runtime metrics, you must register the metrics
    // with the controller-runtime metrics registry
    operatormetrics.Register = operatormetrics.ControllerRuntimeRegister

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

var (
    customResourceCollector = operatormetrics.Collector{
        Metrics: []operatormetrics.CollectorMetric{
            {
                Metric: crCount,
                Labels: []string{"namespace"},
            },
        },
        CollectCallback: customResourceCollectorCallback,
    }

    crCount = operatormetrics.NewGauge(
        operatormetrics.MetricOpts{
            Name:           metricPrefix + "cr_count",
            Help:           "Number of existing guestbook custom resources",
            ConstLabels:    map[string]string{"controller": "guestbook"},
            ExtraFields: map[string]string{
                "StabilityLevel":    "DEPRECATED",
                "DeprecatedVersion": "1.14.0",
            },
        },
    )
)

func customResourceCollectorCallback() []operatormetrics.CollectionResult {
    result := unstructured.UnstructuredList{}
    err := collectorK8sClient.List(context.TODO(), &result, client.InNamespace("default"))
    ...

    return []operatormetrics.CollectorResult{
        {
            Metric: crCount,
            Labels: []string{"default"},
            Value:  float64(len(result.Items)),
        },
    }
}
```

### Prometheus Rules

This section describes how to create and manage Prometheus rules to be
reconciled by your Kubernetes Operator. Prometheus rules are a crucial part of
observability, enabling you to define alerts and record new time series based
on existing metric data.

#### Recording Rules

Recording rules allow you to precompute frequently needed or computationally
expensive expressions and save their result as a new set of time series.

The file holding the recording rules related to the operator workload might look
like this:

```go
// rules/operator_recording_rules.go

var operatorRecordingRules = []operatorrules.RecordingRule{
	...
    {
        MetricsOpts: operatormetrics.MetricOpts{
            Name:        recordingRulesPrefix + "number_of_ready_pods",
            Help:        "Number of ready guestbook operator pods in the cluster",
            ExtraFields: map[string]string{"StabilityLevel": "ALPHA"},
            ConstLabels: map[string]string{"controller": "guestbook"},
        },
        MetricType: operatormetrics.GaugeType,
        Expr:       intstr.FromString(fmt.Sprintf("sum(up{namespace='%s', pod=~'guestbook-operator-.*', ready='true'}) or vector(0)", namespace)),
    },
    ...
}
```

#### Alerts

Alerts notify you when specific conditions are met, such as when a metric value
exceeds a certain threshold or when a system component is unavailable. You can
configure alerts using Prometheus alerting rules.

The file holding the recording rules related to the operator workload might look
like this:

```go
// rules/operator_alerts.go


var operatorAlerts = []promv1.Rule{
	...
    {
        Alert: "GuestbookOperatorNotReady",
        Expr:  intstr.FromString(fmt.Sprintf("%snumber_of_ready_pods < %snumber_of_pods", recordingRulesPrefix, recordingRulesPrefix)),
        For:   "5m",
        Annotations: map[string]string{
            "summary":     "Guestbook operator is not ready",
            "description": "Guestbook operator is not ready for more than 5 minutes.",
        },
        Labels: map[string]string{
            "severity": "critical",
        },
    },
}

```

#### Setup

The rules should be registered in the `init()` function with:

```go
func SetupRules() *promv1.PrometheusRule {
    err := operatorrules.RegisterRecordingRules(recordingRules...)
    ...
    
    err = operatorrules.RegisterAlerts(alerts...)
    ...
    
    prometheusRuleObj, err := operatorrules.BuildPrometheusRule(
        "guestbook-operator-prometheus-rules",
        "default",
        map[string]string{"app": "guestbook-operator"},
    )
    ...
```

### Documentation

Having all resources in one place makes it easy to document them and track the
changes. The documentation of the metrics, recording rules, and alerts can be
generated from the code using docs utilities. The utilities will generate a
string with the documentation that you can later print or save to a file. The
documentation includes a default template that you can customize.

For metrics and recording rules:
```go
func main() {
    metrics.SetupMetrics()
    rules.SetupRules()

    docsString := docs.BuildMetricsDocs(metrics.ListMetrics(), rules.ListRecordingRules())
    fmt.Println(docsString)
}
```

For alerts:
```go
func main() {
    rules.SetupRules()
    docsString := docs.BuildAlertsDocs(alerts.ListAlerts())
    fmt.Println(docsString)
}
```

## Next Steps

- Add a declarative way to define metrics, recording rules, and alerts

- Create a set of macros to make it easier to define metrics, recording rules,
and alerts expressions

- Propose design changes to the Operator SDK examples

- Also add utils for Events and E2E tests

- Build a Kubebuilder/Operator SDK plugin to allow developers to effortlessly
add observability to their operators

- Add Kubebuilder/Operator SDK command line instructions to generate code for
new Metrics, Alerts, and Events
