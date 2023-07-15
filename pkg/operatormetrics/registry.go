package operatormetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type RegistryFunc func(c prometheus.Collector) error

// Register is the function used to register metrics and collectors by this package.
var Register = PrometheusRegister

// Common Register functions that can be used as Register.
var (
	PrometheusRegister        RegistryFunc = prometheus.Register
	ControllerRuntimeRegister RegistryFunc = metrics.Registry.Register
)
