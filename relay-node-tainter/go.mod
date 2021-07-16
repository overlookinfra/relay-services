module github.com/puppetlabs/relay-services/relay-node-tainter

go 1.16

require (
	github.com/dustinkirkland/golang-petname v0.0.0-20191129215211-8e5a1ed0cff0
	github.com/ettle/strcase v0.1.1
	github.com/golangci/golangci-lint v1.41.1
	github.com/google/uuid v1.2.0 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/puppetlabs/leg/k8sutil v0.6.0
	github.com/puppetlabs/leg/mainutil v0.1.2
	github.com/puppetlabs/leg/timeutil v0.4.2
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/time v0.0.0-20210611083556-38a9dc6acbc6
	gotest.tools/gotestsum v1.6.4
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	k8s.io/klog/v2 v2.10.0
	sigs.k8s.io/controller-runtime v0.9.2
	sigs.k8s.io/controller-tools v0.4.1
	sigs.k8s.io/kustomize/kustomize/v3 v3.10.0
)
