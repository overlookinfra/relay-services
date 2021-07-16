package e2e_test

import (
	"context"
	"time"

	"github.com/puppetlabs/leg/k8sutil/pkg/test/endtoend"
	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var backoffFactory = backoff.Build(
	backoff.Exponential(250*time.Millisecond, 2.0),
	backoff.MaxBound(5*time.Second),
	backoff.FullJitter(),
	backoff.NonSliding,
)

func Wait(ctx context.Context, work retry.WorkFunc) error {
	return retry.Wait(ctx, work, retry.WithBackoffFactory(backoffFactory))
}

type getter interface {
	Get(ctx context.Context, key client.ObjectKey, obj client.Object) error
}

type updater interface {
	Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error
}

func mutate(ctx context.Context, getter getter, updater updater, obj client.Object, fn func()) error {
	key := client.ObjectKeyFromObject(obj)
	return retry.Wait(ctx, func(ctx context.Context) (bool, error) {
		// Mutation function.
		fn()

		if err := updater.Update(ctx, obj); errors.IsConflict(err) {
			if err := getter.Get(ctx, key, obj); err != nil {
				return retry.Done(err)
			}

			return retry.Repeat(err)
		} else if err != nil {
			return retry.Done(err)
		}

		return retry.Done(nil)
	})
}

func Mutate(ctx context.Context, e *endtoend.Environment, obj client.Object, fn func()) error {
	return mutate(ctx, e.ControllerClient, e.ControllerClient, obj, fn)
}

func MutateStatus(ctx context.Context, e *endtoend.Environment, obj client.Object, fn func()) error {
	return mutate(ctx, e.ControllerClient, e.ControllerClient.Status(), obj, fn)
}
