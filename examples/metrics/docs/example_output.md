# Operator Metrics

### guestbook_operator_cr_count
[Beta] Number of existing guestbook custom resources. Type: Gauge.

### guestbook_operator_reconcile_action_count
[Alpha] Number of times the operator has executed the reconcile loop with a given action. Type: Counter.

### guestbook_operator_reconcile_count
[Stable] Number of times the operator has executed the reconcile loop. Type: Counter.

## Developing new metrics

All metrics documented here are auto-generated and reflect exactly what is being
exposed. After developing new metrics or changing old ones please regenerate
this document.
