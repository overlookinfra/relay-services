package app_test

import (
	"regexp"
	"testing"

	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/k8sutil/pkg/norm"
	"github.com/puppetlabs/relay-services/relay-node-tainter/pkg/controller/app"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestConfigureNodeTaints(t *testing.T) {
	tests := []struct {
		Name           string
		Taints         []corev1.Taint
		Conditions     []corev1.NodeCondition
		Pattern        string
		ExpectedTaints []corev1.Taint
	}{
		{
			Name: "Basic",
			Conditions: []corev1.NodeCondition{
				{
					Type:   corev1.NodeDiskPressure,
					Status: corev1.ConditionTrue,
				},
			},
			ExpectedTaints: []corev1.Taint{
				{
					Key:    app.TaintKeyPrefix + "disk-pressure",
					Value:  "true",
					Effect: corev1.TaintEffectNoSchedule,
				},
			},
		},
		{
			Name: "Leaves unmanaged taints",
			Taints: []corev1.Taint{
				{
					Key:    "sandbox.gke.io/runtime",
					Value:  "gvisor",
					Effect: corev1.TaintEffectNoSchedule,
				},
			},
			Conditions: []corev1.NodeCondition{
				{
					Type:   corev1.NodeDiskPressure,
					Status: corev1.ConditionTrue,
				},
			},
			ExpectedTaints: []corev1.Taint{
				{
					Key:    "sandbox.gke.io/runtime",
					Value:  "gvisor",
					Effect: corev1.TaintEffectNoSchedule,
				},
				{
					Key:    app.TaintKeyPrefix + "disk-pressure",
					Value:  "true",
					Effect: corev1.TaintEffectNoSchedule,
				},
			},
		},
		{
			Name: "Removes and replaces taints that no longer apply",
			Taints: []corev1.Taint{
				{
					Key:    app.TaintKeyPrefix + "memory-pressure",
					Value:  "true",
					Effect: corev1.TaintEffectNoSchedule,
				},
				{
					Key:    app.TaintKeyPrefix + "disk-pressure",
					Value:  "true",
					Effect: corev1.TaintEffectNoSchedule,
				},
			},
			Conditions: []corev1.NodeCondition{
				{
					Type:   corev1.NodeReady,
					Status: corev1.ConditionTrue,
				},
				{
					Type:   corev1.NodeDiskPressure,
					Status: corev1.ConditionTrue,
				},
				{
					Type:   corev1.NodePIDPressure,
					Status: corev1.ConditionTrue,
				},
			},
			ExpectedTaints: []corev1.Taint{
				{
					Key:    app.TaintKeyPrefix + "disk-pressure",
					Value:  "true",
					Effect: corev1.TaintEffectNoSchedule,
				},
				{
					Key:    app.TaintKeyPrefix + "pid-pressure",
					Value:  "true",
					Effect: corev1.TaintEffectNoSchedule,
				},
			},
		},
		{
			Name: "Only taints using conditions that match pattern",
			Conditions: []corev1.NodeCondition{
				{
					Type:   corev1.NodeDiskPressure,
					Status: corev1.ConditionTrue,
				},
				{
					Type:   corev1.NodePIDPressure,
					Status: corev1.ConditionTrue,
				},
			},
			Pattern: "PID",
			ExpectedTaints: []corev1.Taint{
				{
					Key:    app.TaintKeyPrefix + "pid-pressure",
					Value:  "true",
					Effect: corev1.TaintEffectNoSchedule,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			node := corev1obj.NewNode(norm.MetaName(test.Name))
			node.Object.Spec.Taints = test.Taints
			node.Object.Status.Conditions = test.Conditions

			pattern, err := regexp.Compile(test.Pattern)
			require.NoError(t, err)

			app.ConfigureNodeTaints(node, pattern)
			require.Equal(t, test.ExpectedTaints, node.Object.Spec.Taints, "execution #1")

			// Run twice: this function should be idempotent.
			app.ConfigureNodeTaints(node, pattern)
			require.Equal(t, test.ExpectedTaints, node.Object.Spec.Taints, "execution #2")
		})
	}
}
