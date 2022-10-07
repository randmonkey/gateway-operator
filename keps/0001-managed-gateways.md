---
title: Managed Gateways
status: implementable
---

# Managed Gateways

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [Design](#design)
    - [API Resources](#api-resources)
    - [Controllers](#controllers)
    - [Development Environment](#development-environment)
  - [Test plan](#test-plan)
  - [Graduation Criteria](#graduation-criteria)
- [Production Readiness](#production-readiness)
- [Alternatives](#alternatives)
<!-- /toc -->

## Summary

Historically the [Kong Kubernetes Ingress Controller (KIC)][kic] was used to
manage routing traffic to a back-end [Kong Gateway][kong] and both were deployed
using [Helm Chart][charts].

The purpose of this proposal is to suggest an alternative deployment mechanism
founded on the [operator pattern][operators] which would allow Kong Gateways to
be provisioned in a dynamic and Kubernetes-native way, as well as enabling
automation of Kong cluster operations and management of the Gateway lifecycle.

[kic]:https://github.com/kong/kubernetes-ingress-controller
[kong]:https://github.com/kong/kong
[charts]:https://github.com/kong/charts
[operators]:https://kubernetes.io/docs/concepts/extend-kubernetes/operator/

## Motivation

- streamline deployment and operations of Kong on Kubernetes
- configure and manage Kong on Kubernetes using CRDs and the Kubernetes API (as opposed to Helm templating)
- easily manage and view multiple deployments of the Gateway in a single cluster
- automate historically manual cluster operations for Kong (such as upgrades and scaling)

### Goals

- create a foundational [golang-based][gosdk] operator for Kong
- enable deploying Kong with Kubernetes [Gateway][gwapis] resources
- enable deploying the KIC configured to manage any number `Gateways` (resolves historical [KIC#702][kic702])
- enable automated canary upgrades of Kong (and KIC)
- stay in spec with [Operator Framework][ofrm] standards so that we can be published on [operatorhub][ohub]
- provide easy defaults for deployment while also enabling power users

[gosdk]:https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/
[gwapis]:https://github.com/kubernetes-sigs/gateway-api
[kic702]:https://github.com/Kong/kubernetes-ingress-controller/issues/702
[ofrm]:https://operatorframework.io/
[ohub]:https://operatorhub.io/

## Proposal

### Design

The operator will introduce some new [Custom Resource Definitions (CRDs)][crds]
in addition to building upon resources already defined in [Gateway API][gwapi]:

- `DataPlane`
- `ControlPlane`
- `Gateway` (upstream)
- `GatewayClass` (upstream)
- `GatewayConfiguration`

The `DataPlane` resource will create and manage the lifecycle of `Deployments`,
`Secrets` and `Services` relevant to the [Kong Gateway][kong]. By creating a
`DataPlane` resource on the cluster you can deploy and configure an edge proxy
for ingress traffic. It is expected that this API will be used predominantly as
an implementation detail for `Gateway` resources, and not be used directly by
end-users (though we'll make it possible to do so).

The `ControlPlane` resource will create and manage the lifecycle of `Deployments`,
`Secrets` and `Services` relevant to the [Kong Kubernetes Ingress Controller][kic].
Creating a `ControlPlane` enables the use of `Ingress` and `HTTPRoute` APIs for
the any number of `DataPlanes` the `ControlPlane` is configured to managed.
Similar to the `DataPlane` resource above this is not designed to be an end-user
API.

The upstream `Gateway` and `GatewayClass` resources can be used to create and
manage the lifecycle of both `DataPlanes` and `ControlPlanes`. By default the
creation of a `Gateway` (which is configured with the relevant `GatewayClass`)
results in a default `ControlPlane` and `DataPlane`, with a single proxy
instance listening for ingress traffic and configurable via Kubernetes APIs
(e.g. `Ingress`, `HTTPRoute`, `TCPRoute`, `UDPRoute`, e.t.c.).

Finally, the `GatewayConfiguration` resource allows for implementation specific
configuration of `Gateways`. Some things like listener configuration can be
defined in the `Gateway` resource and that would work across all the
implementations outside of Kong, but things like configuring the max number of
nginx worker processes may not. `GatewayConfiguration` resources can be applied
to `Gateway` resources by way of the `GatewayClass.Spec.ParametersReference`.

[crds]:https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/
[gwapi]:https://github.com/kubernetes-sigs/gateway-api
[kong]:https://github.com/kong/kong
[kic]:https://github.com/kong/kubernetes-ingress-controller

#### Development Environment

We will fundamentally rely on the [Operator SDK][osdk] as the tooling for the
project using the [Golang Operator][gopr] paradigm which will allow us to
quickly and easily ship package updates for [Operator Hub][ohub].

In previous projects such as the KIC we had used [Kubebuilder][kb] which is a
subset of Operator SDK, but the Operator SDK gives us additional features on
top of this mainly for packaging and distributing our operator via the Operator
Hub which is an explicit goal of the project.

At the time of writing the Operator SDK was under [CNCF][cncf] governance as an
[incubating][prjx] project gave us confidence that we can rely on it in the
long term.

[osdk]:https://github.com/operator-framework/operator-sdk
[gopr]:https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/
[ohub]:https://operatorhub.io/
[prjx]:https://www.cncf.io/projects/

#### Supported Scaling Modes & Network Topologies

Providing managed `Gateways` (or headless `DataPlanes`) has implications for
network topology: A `Gateway` may have one or many underlying "instances"\
(think `Pods`, system processes, e.t.c.) and may have a control plane that is
network adjacent, or running somewhere outside of the local cluster. `Gateways`,
`ControlPlanes` and `DataPlanes` can all potentially be vertically or
horizontally scaled resulting in changes to network topology for application and
control plane traffic. The purpose of this section is to document the specific
scaling modes that we intend to support, and highlight the network topologies
they represent according to application and control plane traffic flow. Some
scaling modes can be used in combination with others.

##### Vertical DataPlane Scalability Mode

Vertical DataPlane Scalability Mode is the default and most basic: instances of
the `ControlPlane` and `DataPlane` can be scaled to support more workloads.
This scaling is done "vertically" by changing their CPU or Memory resource
availability and allowing the components to consume more cluster resources
The `Gateway` and `GatewayConfiguration` APIs are used to implement this. This
mode can be combined with other scaling modes.

###### User Interface

Since this is a vertical scaling mode the most important tunables are the CPU
and Memory of the underlying `DataPlane`'s `Pod`. Sane defaults will be applied
to `Pods` to avoid unbounded consumption, and the `GatewayConfiguration` will
provide fields for managing the maximum CPU and memory utilization in order to
increase the ceiling when limits are starting to be hit.

###### Network Topology

As vertical scaling mode can be applied within other topologies it doesn't
inherently imply any change in the network topology in an of itself. However
for the purposes of illustration we will focus on how things look when a single
`Gateway` is used, which is effectively the "default" deployment for users. In
this example ingress traffic passes through a single `DataPlane`, and increases
in traffic requirements can potentially be resolved by increasing CPU/Mem.

![vert-scaling-diag](https://user-images.githubusercontent.com/5332524/187459461-0cbdcd47-d36a-4735-84dc-6ace2e2c0a73.png)

###### Relevant Issues:

- Vertical Scaling: [#233](https://github.com/Kong/gateway-operator/issues/233)
- Performance Testing: [#235](https://github.com/Kong/gateway-operator/issues/235)

##### Horizontal DataPlane Scalability Mode

Horizontal DataPlane Scalability mode implies that there will be 2 or more
`Pods` provisioned for a `DataPlane` resource. In this mode traffic is split
between the backend `Pods` according to the `Service` (generally, implemented
with a `LoadBalancer` strategy and managed with a simple balancing algorithm
such as "round-robin").

With multiple `Pods` in action the Gateway Operator is responsible for
deploying a single `ControlPlane` instance which will serve equivalent
configuration to each of them. Regular tests need to be run to determine the
upper bounds of this integration and verify how many `DataPlanes` a single
`ControlPlane` can effectively serve with fast turnaround.

###### User Interface

The `Gateway` API must be used to employ this strategy, as the gateway
controller will be responsible for responding to scaling requests and
provisioning the relevant `DataPlanes`, as well as reconfiguring the
relevant `ControlPlanes` to become connected to them and synchronized.

Scaling requests in this mode will be handled using `GatewayScalingPolicy`,
which is an implementation of a [Policy Attachment][pol] API. The policy
includes a field in the spec for an explicit number of replicas to provision,
as well as options for automatically scaling the `Gateway` according to
[provided resource metrics or custom metrics][hpa]. By default all `Pods` will
be provisioned with node anti-affinity for each other, though this will also be
configurable through the scaling policy.

###### Network Topology

In this setup application traffic is load-balanced between any number of `Pods`
all of which are opportunistically placed on separate `Nodes`. `ControlPlane`
traffic is transmitted via the `Pod` network to the `DataPlane` from a network
adjacent location.

![horz-scaling-diag](https://user-images.githubusercontent.com/5332524/187459519-ccf132db-293f-493b-8fcc-4262ac698636.png)

###### Relevant Issues

- Horizontal DataPlane Scaling: [#8](https://github.com/Kong/gateway-operator/issues/8)
- Automatic Horizontal DataPlane Scaling: [#171](https://github.com/Kong/gateway-operator/issues/171)
- Performance Testing: [#235](https://github.com/Kong/gateway-operator/issues/235)

[pol]:https://gateway-api.sigs.k8s.io/references/policy-attachment/
[hpa]:https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/

##### Hybrid DataPlane Scalability Mode

Hybrid DataPlane Scalability Mode is similar to the above mode in that it
provides a means to horizontally scale using `DataPlanes` but in this case the
`ControlPlane` is missing, as the control-plane is expected to be provided
elsewhere (in the future, we may add support for Hybrid ControlPlane deployment
but that's not in scope for the moment).

###### User Interface

Unlike the above modes, this is instrumented entirely by using the `DataPlane`
API directly, and wont be usable with `Gateway`. The caller needs to be able to
configure both the ENV and the mounts of the `Pods` which are created for the
`DataPlane` so as to provide the network and certificate configurations needed
so that the underlying Kong Gateway can communicate back with the Hybrid
control plane.

###### Network Topology

In this setup application traffic is load-balanced between any number of `Pods`
all of which are opportunistically placed on separate `Nodes`. `ControlPlane`
traffic is transmitted via webhook initialized by the `DataPlane` from a
control plane unmanaged by the Gateway Operator, that traffic can either be
intra-cluster OR extra cluster (across the internet, e.t.c.).

![hybrid-scaling-diag](https://user-images.githubusercontent.com/5332524/187458947-83c4f5a0-4ba3-4f8f-a3da-e3aa52d081a8.png)

###### Relevant Issues

- Hybrid DataPlane Mode Support: [#229](https://github.com/Kong/gateway-operator/issues/229)

### Test Plan

Testing for this new operator will be performed similarly to what's already
performed in the [KIC][kic] which will include:

- unit tests for all Go packages using [go tests][gotest]
- integration tests using the [Kong Kubernetes Testing Framework (KTF)][ktf]
- e2e tests using [KTF][ktf]

All tests should be able to run locally using `go test` or for integration and
e2e tests using a local system Kubernetes deployment like [Kubernetes in Docker
(KIND)][kind].

[kic]:https://github.com/kong/kubernetes-ingress-controller
[gotest]:https://pkg.go.dev/testing
[ktf]:https://github.com/kong/kubernetes-testing-framework
[kind]:https://github.com/kubernetes-sigs/kind

### Graduation Criteria

The following milestones are considered prerequisites for a generally available
`v1` release of the Gateway Operator:

#### Milestone - Basic Install

This milestone covers the basic functionality and automation needed for a
simple, traditional edge proxy deployment with Kubernetes API support which is
instrumented with `Gateway` resources (and optionally `GatewayConfigurations`).

This milestone corresponds with [Operator Capabilities Level 1 "Basic
Install][ocap].

View this milestone and all its issues on Github [here][gom1].

[ocap]:https://operatorframework.io/operator-capabilities/
[gom1]:https://github.com/Kong/gateway-operator/milestone/1

#### Milestone - Manual Canary Upgrades

This milestone covers the functionality neededed to perform a canary upgrade
of a `Gateway`.

This milestone corresponds with [Operator Capabilities Level 2 "Seamless
Upgrades"][ocap].

View this milestone and all its issues on Github [here][gom8].

[ocap]:https://operatorframework.io/operator-capabilities/
[gom8]:https://github.com/Kong/gateway-operator/milestone/8

#### Milestone - Automated Upgrades

This milestone covers automating canary upgrades so that users can define
a strategy for automatic upgrades and testing automatically completes the
upgrade.

This milestone corresponds with [Operator Capabilities Level 3 "Full
Lifecycle"][ocap].

View this milestone and all its issues on Github [here][gom9].

[ocap]:https://operatorframework.io/operator-capabilities/
[gom9]:https://github.com/Kong/gateway-operator/milestone/9

#### Milestone - Backup & Restore

This milestone covers backing up and restoring the state of a `Gateway` so that
it can be restored to a new cluster, or a previous state can be restored.

This milestone corresponds with [Operator Capabilities Level 3 "Full
Lifecycle"][ocap].

View this milestone and all its issues on Github [here][gom12].

[ocap]:https://operatorframework.io/operator-capabilities/
[gom12]:https://github.com/Kong/gateway-operator/milestone/12

#### Milestone - Monitoring

Integration with [Prometheus][prometheus] to provide monitoring and insights
and reporting of failures (such as upgrade or backup failures).

This milestone corresponds with [Operator Capabilities Level 4 "Deep
Insights"][ocap].

View this milestone and all its issues on Github [here][gom14].

[prometheus]:https://prometheus.io/
[ocap]:https://operatorframework.io/operator-capabilities/
[gom14]:https://github.com/Kong/gateway-operator/milestone/14

#### Milestone - Autoscaling

This milestone covers automating `Gateway` scaling to scale the number of
underlying instances to support growing traffic dynamically.

This milestone corresponds with [Operator Capabilities Level 5 "Auto
Pilot"][ocap].

View this milestone and all its issues on Github [here][gom15].

[ocap]:https://operatorframework.io/operator-capabilities/
[gom15]:https://github.com/Kong/gateway-operator/milestone/15

#### Milestone - Documentation

This milestone covers getting full and published documentation for the Gateway
Operator published on the [Kong Documentation Site][kongdocs] under it's own
product listing.

Notes about documentation: the original vision for the documentation for this
project was for documentation to be considered one of the most (if not the
most) important components of the software. A "guide first" approach is desired
based on the expected common use cases and paths, with advanced configurations
and functionality nested under some caveats which clearly define the actual
situations under which they would be required.

Additionally the following key features are desired:

- working examples provided for all configuration options (within reason)
- documentation is _testable_ and includes testing CI workflows: CI literally
  can test the examples and ensure they continue working over time (this may
  require the help of the maintainers of docs.konghq.com to help us add the
  relevant provisions)
- a compatibility matrix for cloud providers AND Kubernetes versions that is
  automatically updated based on testing in the operator repository (e.g. when
  new version support is added, CI automatically puts in a docs PR to update the
  version matrix)
- upfront and complete guides for any cloud providers we support: (for example,
  (but not decided at the time of writing): GKE, EKS, AKS, Openshift)

View this milestone and all its issues on Github [here][gom16].

[kongdocs]:https://docs.konghq.com
[gom16]:https://github.com/Kong/gateway-operator/milestone/16

## Production Readiness

Production readiness of this operator is marked by the following requirements:

- [ ] All milestones of the above `Graduation Criteria` have been completed
- [ ] Unit, integration and E2E tests are present at a high level of coverage

## Alternatives

### Blue/Green Upgrades

Originally we considered automated blue/green upgrades for our `v1` release of
the operator, however we decided that canary upgrades would be sufficient for
the initial release as this would reduce scope while also focusing on what we
expected to be the more commonly used strategy. We will be able to revisit
whether to add blue/green upgrades as part of a later iteration and separate KEP.

### PostGreSQL Mode

For the first iteration we are _only_ supporting DBLESS mode deployments of
Kong as this is the preferred operational mode on Kubernetes and PostGreSQL
mode adds a lot of burden and complexity.
