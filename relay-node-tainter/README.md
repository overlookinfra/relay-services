# Node Tainter [![Node Tainter CI](https://github.com/puppetlabs/relay-services/actions/workflows/relay-node-tainter-ci.yaml/badge.svg)](https://github.com/puppetlabs/relay-services/actions/workflows/relay-node-tainter-ci.yaml)

The node tainter is a small Kubernetes controller that watches for changes in node conditions and propagates them to node taints. It is functionally similar to the Taint Nodes by Condition feature of the control plane's node controller, but uses a regular expression to match conditions to include as taints.

Conditions are propagated to taints when their status changes to `True`, and taints are automatically removed when their corresponding condition status changes to any other value. (Note that the Ready condition is always ignored as it is handled specially anyway.)

By default, the matching regular expression is empty, so all conditions will match it.

Condition types are automatically transformed to a suitable format for taints, so for example, the condition type `KernelDeadlock` would correspond to the taint `node-tainter.services.relay.sh/condition.kernel-deadlock`. Conditions are given the value `"true"` and the effect `NoSchedule`.

This controller works best with the [Kubernetes descheduler](https://github.com/kubernetes-sigs/descheduler) to automatically move pods off a node when a taint is applied.
