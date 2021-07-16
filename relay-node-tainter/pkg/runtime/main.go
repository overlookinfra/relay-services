package runtime

import (
	"context"
	"flag"
	"fmt"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/eventctx"
	"github.com/puppetlabs/leg/mainutil"
	"github.com/puppetlabs/relay-services/relay-node-tainter/pkg/opt"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func Main(cfg *opt.Config, opts manager.Options, transforms ...func(mgr manager.Manager) error) int {
	return mainutil.TrapAndWait(context.Background(), func(ctx context.Context) error {
		defer klog.Flush()

		flag.Parse()

		kfs := flag.NewFlagSet("klog", flag.ExitOnError)
		klog.InitFlags(kfs)

		if cfg.Debug {
			_ = kfs.Set("v", "5")
		}

		log.SetLogger(klogr.NewWithOptions(klogr.WithFormat(klogr.FormatKlog)))

		if opts.Scheme == nil {
			opts.Scheme = scheme.Scheme
		}

		// These are always overridden by us, but turning on leader election is
		// up to the caller.
		opts.LeaderElectionResourceLock = resourcelock.LeasesResourceLock
		opts.LeaderElectionID = fmt.Sprintf("%s.lease.node-tainter.services.relay.sh", cfg.Name)
		opts.LeaderElectionNamespace = cfg.Namespace
		opts.LeaderElectionReleaseOnCancel = true

		mgr, err := manager.New(config.GetConfigOrDie(), opts)
		if err != nil {
			return fmt.Errorf("failed to create manager: %w", err)
		}

		for i, transform := range transforms {
			if err := transform(mgr); err != nil {
				return fmt.Errorf("failed to apply manager transform #%d: %w", i, err)
			}
		}

		return mgr.Start(eventctx.WithEventRecorder(ctx, mgr, cfg.Name))
	})
}
