package rules

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/intstr"

	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

var operatorAlerts = []promv1.Rule{
	{
		Alert: "GuestbookOperatorDown",
		Expr:  intstr.FromString(fmt.Sprintf("%snumber_of_pods == 0", recordingRulesPrefix)),
		Annotations: map[string]string{
			"summary":     "Guestbook operator is down",
			"description": "Guestbook operator is down for more than 5 minutes.",
		},
	},
	{
		Alert: "GuestbookOperatorNotReady",
		Expr:  intstr.FromString(fmt.Sprintf("%snumber_of_ready_pods < %snumber_of_pods", recordingRulesPrefix, recordingRulesPrefix)),
		For:   "5m",
		Annotations: map[string]string{
			"summary":     "Guestbook operator is not ready",
			"description": "Guestbook operator is not ready for more than 5 minutes.",
		},
	},
}
