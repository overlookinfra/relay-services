package reconciler

import (
	"context"
	"regexp"
	"time"

	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/relay-services/relay-node-tainter/pkg/controller/app"
	"github.com/puppetlabs/relay-services/relay-node-tainter/pkg/opt"
	"golang.org/x/time/rate"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;update

type NodeReconciler struct {
	cl      client.Client
	pattern *regexp.Regexp
}

var _ reconcile.Reconciler = &NodeReconciler{}

func (nr *NodeReconciler) Reconcile(ctx context.Context, req reconcile.Request) (r reconcile.Result, err error) {
	klog.InfoS("node reconciler: starting reconcile for node", "node", req.Name)
	defer klog.InfoS("node reconciler: ending reconcile for node", "node", req.Name)
	defer func() {
		if err != nil {
			klog.ErrorS(err, "node reconciler: failed to reconcile node", "node", req.Name)
		}
	}()

	node := corev1obj.NewNode(req.Name)
	if ok, err := node.Load(ctx, nr.cl); err != nil || !ok {
		return reconcile.Result{}, err
	}

	app.ConfigureNodeTaints(node, nr.pattern)

	err = node.Persist(ctx, nr.cl)
	return
}

func NewNodeReconciler(cl client.Client, pattern *regexp.Regexp) *NodeReconciler {
	return &NodeReconciler{
		cl:      cl,
		pattern: pattern,
	}
}

func AddNodeReconcilerToManager(mgr manager.Manager, cfg *opt.Config) error {
	rl := workqueue.NewMaxOfRateLimiter(
		workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, cfg.ControllerMaxReconcileBackoffDuration),
		&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(10), 100)},
	)

	r := NewNodeReconciler(mgr.GetClient(), cfg.ControllerConditionTypeInclusionPattern)

	return builder.ControllerManagedBy(mgr).
		For(&corev1.Node{}).
		WithOptions(controller.Options{RateLimiter: rl}).
		Complete(r)
}
