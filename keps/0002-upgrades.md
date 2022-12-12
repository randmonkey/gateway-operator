---
title: Manual upgrades
status: provisional
---

# Manual upgrades

<!-- toc -->
- [Manual upgrades](#manual-upgrades)
  - [Summary](#summary)
  - [Goals](#goals)
  - [Proposal](#proposal)
    - [How to upgrade/downgrade](#how-to-upgradedowngrade)
    - [DataPlane upgrades](#dataplane-upgrades)
      - [ControlPlane upstream changes](#controlplane-upstream-changes)
      - [DataPlane upstream changes](#dataplane-upstream-changes)
      - [Gateway operator changes](#gateway-operator-changes)
      - [Alternative solution](#alternative-solution)
    - [ControlPlane upgrades](#controlplane-upgrades)
  - [Further improvements](#further-improvements)
<!-- /toc -->

## Summary

As a user, I want to upgrade or downgrade my `DataPlane`s and `ControlPlane`s
without downtime, i.e., the operation should not affect the traffic routed by the
`DataPlane` in any way, both after and in the meanwhile. Therefore, this document
aims at proposing a mechanism to implement the said user story in a way that:

- No route that was previously working should return 404 during a `DataPlane` upgrade/downgrade.
- No route that was previously working should return 404 during a `ControlPlane`
upgrade/downgrade.

## Goals

- upgrade/downgrade `DataPlane`s without traffic interruption
- upgrade/downgrade `ControlPlane`s without traffic interruption

## Proposal

> **Note**: In this proposal, upgrades and downgrades have the same meaning, i.e.,
> rolling out a component from version x to version y, no matter if x>y or y>x.

### How to upgrade/downgrade

`ControlPlane` and `DataPlane` resources can be managed, i.e., created by the Gateway
controller starting from a `Gateway` resource, or unmanaged, i.e., directly
created by the user.
Depending on the resource type (managed/unmanaged), the upgrade/downgrade should
be performed in two different ways:

- For managed resources, change the following field in the `GatewayConfiguration`
bound to the `Gateway` that owns the resource:
  - `.spec.dataPlaneDeploymentOptions.version` for `DataPlane`s
  - `.spec.controlPlaneDeploymentOptions.version` for `ControlPlane`s
- For unmanaged resources, directly change the `.spec.version` field in the resource.

### DataPlane upgrades

During `DataPlane`s version upgrades, no previously working routes should return
404. Since the `DataPlane` is the component exposed by the proxy `LoadBalancer`
Service and is the point where all the traffic flows and is routed, to ensure
seamless upgrades, there should always be exactly one DataPlane instance able to
route traffic. The condition to have a `DataPlane` instance correctly routing
traffic is that the `ControlPlane` has pushed the configuration to the `DataPlane`
instance, and the `DataPlane` has loaded such configuration.

#### ControlPlane upstream changes

The `ControlPlane` (KIC) must be able to [discover new `DataPlane` instances and
interact with them][KIC-702]. For this purpose, the `ControlPlane` should be able
to accept multiple `DataPlane` endpoints instead of a single one, and instantiate
a client with each `DataPlane` pod to push the configuration.
The following logic needs to be added to the `ControlPlane`:

1. At startup time, it parses the endpoints given as env var or argument.
2. It adds a new `DataPlane` client to the [synchronizer][synchronizer] to
interact with all the `DataPlane` pods and push the configuration to them.

[KIC-702]:https://github.com/kong/kubernetes-ingress-controller/issues/702
[synchronizer]:https://github.com/Kong/kubernetes-ingress-controller/blob/0243a95453a266c712a006a03fd0c763e4181ebb/internal/dataplane/synchronizer.go#L35-L50

#### DataPlane upstream changes

The `DataPlane` (Kong) needs to set its readiness only after the ControlPlane
correctly pushes the configuration. The [FTI-3081][FTI-3081] tracks this high-level
objective and is planned for Kong 3.2. As an interim workaround, the Gateway operator
will add a delay into the `DataPlane` readiness probe to have a certain level of
confidence that once the `DataPlane` gets ready, the `ControlPlane` has already
pushed the configuration.

[FTI-3081]:https://konghq.atlassian.net/browse/FTI-3081

#### Gateway operator changes

The operator must introduce the following changes:

1. Whenever the `DataPlane` version gets updated (through one of the mechanisms
specified [above](#how-to-upgradedowngrade)), the `DataPlane` controller enforces
such a version in the `DataPlane` deployment.
1. The `DataPlane` controller enforces the `.spec.strategy` field in the `DataPlane`
deployment so that the old `DataPlane` pod gets deleted only once the new `DataPlane`
pod is running.
1. A new controller watches for all the `DataPlane` pods and patches the controlPlane
with all the new `DataPlane` pod addresses.
1. Add a delay to the `DataPlane` readiness probe to have a certain level of
confidence that once the `DataPlane` gets ready, the `ControlPlane` has already
pushed the configuration. (workaround for FTI-3081).
1. Remove the static Kong endpoint configuration from the `ControlPlane`.
1. Remove the admin port from the Proxy `LoadBalancer` service, as the `ControlPlane`
communicates with the `DataPlane` through the pod network.

#### Alternative solution

The issue in statically passing the endpoints to the `ControlPlane` is binding
the `ControlPlane` pod lifecycle to the `DataPlane` pod lifecycle (i.e., every
time the `DataPlane` changes, the `ControlPlane` must be updated as well, leading
to the rollout of the `ControlPlane` pods).
For this reason, a dynamicity layer can be added through a new `DataPlaneInstance`
API, structured as follows:

```go
// DataplaneInstanceSpec defines the desired state of DataplaneInstance
type DataPlaneInstanceSpec struct {
  // Address is the address through which the DataPlane instance is exposed
  //
  // +kubebuilder:validation:Required
  // +kubebuilder:validation:Pattern=^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$
  Address string `json:"address"`
}

// DataplaneInstanceStatus defines the observed state of DataplaneInstance
type DataplaneInstanceStatus struct {
  // Conditions describe the current conditions of the DataplaneInstance.
  //
  // +optional
  // +listType=map
  // +listMapKey=type
  // +kubebuilder:validation:MaxItems=8
  // +kubebuilder:default={{type: "Scheduled", status: "Unknown", reason:"NotReconciled", message:"Waiting for controller", lastTransitionTime: "1970-01-01T00:00:00Z"}}
  Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DataplaneInstance is the Schema for the dataplaneinstances API
type DataplaneInstance struct {
  metav1.TypeMeta   `json:",inline"`
  metav1.ObjectMeta `json:"metadata,omitempty"`

  Spec   DataplaneInstanceSpec   `json:"spec,omitempty"`
  Status DataplaneInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DataplaneInstanceList contains a list of DataplaneInstance
type DataplaneInstanceList struct {
  metav1.TypeMeta `json:",inline"`
  metav1.ListMeta `json:"metadata,omitempty"`
  Items           []DataplaneInstance `json:"items"`
}
```

Such a resource has a 1:1 relationship with the `DataPlane` Pods and is reconciled
by a new controller, which performs the following operations:

1. For each new `DataPlaneInstance`, it adds a new `DataPlane` client to the
[synchronizer][synchronizer] to interact with the new `DataPlane` pod and push
the configuration.
1. For each `DataPlaneInstance` deleted, it removes the corresponding `DataPlane`
client from the synchronizer to stop pushing the configuration to the `DataPlane`.
1. Update each `DataPlaneInstance` status with the proper condition to give feedback
about the configuration push.

The `DataPlaneInstance` resource exposes the following condition to give a proper
indication of the resource lifecycle:

```go
  // This condition indicates whether the DataPlaneInstance has been
  // configured by the ControlPlane.
  //
  // Possible reasons for this condition to be true are:
  //
  // * "Ready"
  //
  // Possible reasons for this condition to be False are:
  //
  // * "Invalid"
  // * "Pending"
  //
  DataPlaneInstanceConditionReady DataPlaneInstanceConditionType = "Ready"

  // This reason is used with the "Ready" condition when the condition is
  // true.
  DataPlaneInstanceReasonReady DataPlaneInstanceConditionReason = "Ready"

  // This reason is used with the "Ready" condition when the
  // configuration contains some errors, and the push to the
  // DataPlane failed
  DataPlaneInstanceReasonInvalid DataPlaneInstanceConditionReason = "ConfigurationError"

  // This reason is used with the "Ready" condition when the
  // DataPlane is not ready and the configuration is not pushed yet
  DataPlaneInstanceReasonPending DataPlaneInstanceConditionReason = "Pending"
```

The utilization of such an API and the dynamic instantiation of sync loops is
exposed through a feature gate. With that feature gate, the user can use the
static Kong endpoint configuration (that will continue to be used by KIC as a
standalone), or the dynamic discovery, which should be enabled by default when KIC
is created by the Gateway Operator.

The above API is created by the Gateway operator starting from the `DataPlane`
pods (each new `DataPlane` pod implies the creation of a new `DataPlaneInstance`
resource) and reconciled by the `ControlPlane`, to allow dynamic customization
of `DataPlane` endpoints in the `ControlPlane`.

### ControlPlane upgrades

The `ControlPlane` upgrades should not require ad-hoc changes, as they are not as
sensitive as the `DataPlane` ones. Once the user upgrades the `ControlPlane` version
(using one of the methods described [above](#how-to-upgradedowngrade)), the
`ControlPlane` controller is in charge of enforcing such a version into the `ControlPlane`
Deployment, and the deployment rolls out. No traffic disruption is expected to happen,
as the configuration is already pushed into the `DataPlane` by the previous `ControlPlane`
instance, and everything works as usual. The level of disruption that can happen
is having a new K8s resource (`Gateway`, `xxxRoute`, ...) or an update on an
existing one during the `ControlPlane` upgrades; in such a scenario, the affected
resource is not taken into consideration (i.e., pushed into the `DataPlane`) until
the new `ControlPlane` pod gets ready.

## Further improvements

The natural improvement for the manual upgrades is the canary upgrades, which will
require an interface to let the user perform canary upgrades for `ControlPlane`s
and `DataPlane`s.
