# Use a Different Registry for metrics and collectors

Sometime packages like
[kubernetes-sigs/controller-runtime](https://github.com/kubernetes-sigs/controller-runtime/tree/main/pkg/metrics)
implement their own `Registry` for metrics and collectors or overwrite the
`DefaultRegisterer` from 
[prometheus/client_golang](https://github.com/prometheus/client_golang/blob/main/prometheus/registry.go#L56).

This means that if your application was previously using one package that uses
that pattern, you will need to set the `Register` function to the one from the
package you are using.

## Example

```go
// To use the Prometheus Registerer from client_golang (Default - so can be omitted)
operatormetrics.Register = operatormetrics.PrometheusRegister

// To use the Prometheus Registerer from controller-runtime
operatormetrics.Register = operatormetrics.ControllerRuntimeRegister
// or
import "sigs.k8s.io/controller-runtime/pkg/metrics"
operatormetrics.Register = metrics.Registry.Register

// To use an entirely different Registerer
import "github.com/prometheus/client_golang/prometheus"
operatormetrics.Register = func(collector prometheus.Collector) error {
    // Register the collector
    return nil
}
```
