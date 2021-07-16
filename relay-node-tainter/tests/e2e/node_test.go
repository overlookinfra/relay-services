package e2e_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/ettle/strcase"
	"github.com/puppetlabs/leg/k8sutil/pkg/test/endtoend"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	"github.com/puppetlabs/relay-services/relay-node-tainter/pkg/controller/app"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestNodeTaintLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	WithEnvironmentInTest(t, func(e *endtoend.Environment) {
		// Create a condition type.
		name := petname.Generate(3, "-")
		cond := corev1.NodeConditionType(strcase.ToPascal(name))

		// Pick a node. Doesn't matter which one.
		var node *corev1.Node
		require.NoError(t, Wait(ctx, func(ctx context.Context) (bool, error) {
			nodes := &corev1.NodeList{}
			if err := e.ControllerClient.List(ctx, nodes); err != nil {
				return retry.Done(err)
			}

			if len(nodes.Items) == 0 {
				return retry.Repeat(fmt.Errorf("waiting for nodes"))
			}

			node = &nodes.Items[0]
			return retry.Done(nil)
		}))

		// Add the condition.
		require.NoError(t, MutateStatus(ctx, e, node, func() {
			node.Status.Conditions = append(node.Status.Conditions, corev1.NodeCondition{
				Type:   cond,
				Status: corev1.ConditionTrue,
				Reason: "InTest",
			})
		}))

		// Now wait until the node has the taint applied.
		require.NoError(t, Wait(ctx, func(ctx context.Context) (bool, error) {
			if err := e.ControllerClient.Get(ctx, client.ObjectKeyFromObject(node), node); err != nil {
				return retry.Done(err)
			}

			for _, taint := range node.Spec.Taints {
				if taint.Key == app.TaintKeyPrefix+name {
					return retry.Done(nil)
				}
			}

			return retry.Repeat(fmt.Errorf("waiting for node to be tainted for condition %s", cond))
		}))

		// Update the condition again and wait for the taint to be removed.
		require.NoError(t, MutateStatus(ctx, e, node, func() {
			for i := range node.Status.Conditions {
				candidate := &node.Status.Conditions[i]
				if candidate.Type != cond {
					continue
				}

				candidate.Status = corev1.ConditionFalse
				candidate.Reason = "NotInTest"
				break
			}
		}))

		// Wait until the node no longer has the taint.
		require.NoError(t, Wait(ctx, func(ctx context.Context) (bool, error) {
			if err := e.ControllerClient.Get(ctx, client.ObjectKeyFromObject(node), node); err != nil {
				return retry.Done(err)
			}

			for _, taint := range node.Spec.Taints {
				if taint.Key == app.TaintKeyPrefix+name {
					return retry.Repeat(fmt.Errorf("waiting for node taint to be removed for condition %s", cond))
				}
			}

			return retry.Done(nil)
		}))
	})
}
