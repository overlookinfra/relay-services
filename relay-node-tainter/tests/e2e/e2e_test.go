package e2e_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/puppetlabs/leg/k8sutil/pkg/test/endtoend"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func init() {
	log.SetLogger(klogr.NewWithOptions(klogr.WithFormat(klogr.FormatKlog)))
}

func WithEnvironmentInTest(t *testing.T, fn func(e *endtoend.Environment)) {
	viper.SetEnvPrefix("relay_node_tainter_test_e2e")
	viper.AutomaticEnv()

	kubeconfigs := strings.TrimSpace(viper.GetString("kubeconfig"))
	if testing.Short() {
		t.Skip("not running end-to-end tests with -short")
	} else if kubeconfigs == "" {
		t.Skip("not running end-to-end tests without one or more Kubeconfigs specified by RELAY_NODE_TAINTER_TEST_E2E_KUBECONFIG")
	}

	opts := []endtoend.EnvironmentOption{
		endtoend.EnvironmentWithClientKubeconfigs(filepath.SplitList(kubeconfigs)),
		endtoend.EnvironmentWithClientContext(viper.GetString("context")),
	}
	require.NoError(t, endtoend.WithEnvironment(opts, fn))
}
