package main

import (
	"log"
	"os"

	"github.com/puppetlabs/relay-services/relay-node-tainter/pkg/controller/reconciler"
	"github.com/puppetlabs/relay-services/relay-node-tainter/pkg/opt"
	"github.com/puppetlabs/relay-services/relay-node-tainter/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {
	cfg, err := opt.NewConfig("relay-node-tainter-controller")
	if err != nil {
		log.Fatalf("configuration error: %s", err)
	}

	os.Exit(runtime.Main(
		cfg,
		manager.Options{
			LeaderElection: true,
		},
		func(mgr manager.Manager) error {
			return reconciler.AddNodeReconcilerToManager(mgr, cfg)
		},
	))
}
