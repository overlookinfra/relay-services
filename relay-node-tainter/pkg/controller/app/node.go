package app

import (
	"regexp"
	"strings"

	"github.com/ettle/strcase"
	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/k8sutil/pkg/norm"
	corev1 "k8s.io/api/core/v1"
)

const TaintKeyPrefix = "node-tainter.services.relay.sh/condition."

func ConfigureNodeTaints(node *corev1obj.Node, pattern *regexp.Regexp) {
	// Clear taints ahead of recheck from conditions.
	var keep int
	for _, taint := range node.Object.Spec.Taints {
		if strings.HasPrefix(taint.Key, TaintKeyPrefix) {
			continue
		}

		node.Object.Spec.Taints[keep] = taint
		keep++
	}

	node.Object.Spec.Taints = node.Object.Spec.Taints[:keep]

	// Add new taints for relevant conditions.
	for _, cond := range node.Object.Status.Conditions {
		if cond.Type == corev1.NodeReady {
			// Ignore this special condition handled by the control plane.
			continue
		}

		if cond.Status != corev1.ConditionTrue {
			// Only add conditions that are currently active.
			continue
		}

		if !pattern.MatchString(string(cond.Type)) {
			// Not matched by the inclusion pattern.
			continue
		}

		taint := corev1.Taint{
			Key:    norm.AnyQualifiedName(TaintKeyPrefix + strcase.ToKebab(string(cond.Type))),
			Value:  "true",
			Effect: corev1.TaintEffectNoSchedule,
		}
		node.Object.Spec.Taints = append(node.Object.Spec.Taints, taint)
	}
}
