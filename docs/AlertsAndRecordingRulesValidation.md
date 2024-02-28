# Alert and Recording Rules Validation

The operator-observability toolkit provides a set of utilities to validate
Prometheus alerts and recording rules for Kubernetes Operators. This document
outlines the purpose and usage of these validators.

Check out the [_examples/tools/lint/](../_examples/tools/lint/) for an example
usage.

## Overview

Alerts and recording rules are essential components of observability in
Kubernetes Operators. However, to ensure consistency, reliability, and adherence
to best practices, it's crucial to validate these components before they're
applied. 

### Initialization

Use the `New()` function to create a new instance of the Linter.
This initializes the linter with default validations, where you can add
additional custom validations.

### Default Validations

The toolkit provides default validations for both recording rules and alerts:

**defaultRecordingRuleValidation:** Validates that the recording rule:
- has a name
- has an expression.

**defaultAlertValidation:** Validates that the alert:
- has a name in PascalCase format.
- has an expression.
- includes a severity label (critical, warning, or info).
- includes summary and description annotations.


### Adding Custom Validations

#### Add custom validation functions for recording rules.

```
type RecordRuleValidation = func(rr *operatorrules.RecordingRule) []Problem

func (linter *Linter) AddCustomRecordRuleValidations(validations ...RecordRuleValidation)
```

#### Add custom validation functions for alerts.

```
type AlertValidation = func(alert *promv1.Rule) []Problem

func (linter *Linter) AddCustomAlertValidations(validations ...AlertValidation)
```

You can define your own alert validation rules or use some custom validations exported and
available for usage in [pkg/testutil/alert_custom_validations.go](../pkg/testutil/alert_custom_validations.go).

### Linting

`LintRecordingRules(recordingRules []operatorrules.RecordingRule) []Problem`:
Lint a slice of recording rules and return a slice of problems found.

`LintRecordingRule(recordingRule operatorrules.RecordingRule) []Problem`: Lint a
single recording rule and return a slice of problems found.

`LintAlerts(alerts []promv1.Rule) []Problem`: Lint a slice of alerts and return
a slice of problems found.

`LintAlert(alert promv1.Rule) []Problem`: Lint a single alert and return a slice
of problems found.
