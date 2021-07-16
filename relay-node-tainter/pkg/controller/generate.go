//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen rbac:roleName=relay-node-tainter-controller paths=./... output:artifacts:config=../../manifests/controller/generated

package controller
