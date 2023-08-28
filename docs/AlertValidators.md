# Alert and Recording Rules Validators

The operator-observability toolkit provides a set of utilities to validate
Prometheus alerts and recording rules for Kubernetes Operators. This document
outlines the purpose and usage of these validators.

## Overview

Alerts and recording rules are essential components of observability in
Kubernetes Operators. However, to ensure consistency, reliability, and adherence
to best practices, it's crucial to validate these components before they're
applied. 

### Initialization

Use the `New()` function to create a new instance of the Linter.
This initializes the linter with default validations, where you can add
additional custom validations.

### Adding Custom Validations

`AddCustomAlertValidations(validations ...AlertValidation)`: Add custom
validation functions for alerts.

`AddCustomRecordRuleValidations(validations ...RecordRuleValidation)`: Add
custom validation functions for recording rules.

### Linting Functions

`LintAlerts(alerts []promv1.Rule) []Problem`: Lint a slice of alerts and return
a slice of problems found.

`LintAlert(alert promv1.Rule) []Problem`: Lint a single alert and return a slice
of problems found.

`LintRecordingRules(recordingRules []operatorrules.RecordingRule) []Problem`:
Lint a slice of recording rules and return a slice of problems found.

`LintRecordingRule(recordingRule operatorrules.RecordingRule) []Problem`: Lint a
single recording rule and return a slice of problems found.

### Default Validators

The toolkit provides default validators for both recording rules and alerts:

**defaultRecordingRuleValidation:** Validates that a recording rule has a name and an expression.

**defaultAlertValidation:** Validates that:
- An alert has a name in PascalCase format.
- An alert has an expression.
- The alert includes a severity label (critical, warning, or info).
- The alert includes summary and description annotations.
