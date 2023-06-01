---
title: API Extensions
status: provisional
---

# API Extensions

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [Design](#design)
  - [Test plan](#test-plan)
  - [Graduation Criteria](#graduation-criteria)
- [Production Readiness](#production-readiness)
- [Alternatives](#alternatives)
<!-- /toc -->

## Summary

We want to be able to add extensions to the Kong Gateway Operator (KGO) so that
we can provide additional or alternative behaviors to APIs such as `DataPlanes`.

## Motivation

- enable custom builds of the KGO
- enable feature extension of the KGO for use cases out of scope

### Goals

- develop a standard for how APIs in the KGO can be extended
- provide common library functionality to support extensions
- provide extension support for the `DataPlane` resource
- provide extension support for the `ControlPlane` resource
- provide extension support for the `Gateway` resource

## Proposal

Our APIs will be extendable by way of "extension fields" within each resource
which can point to another Kubernetes resource that defines the extension
behavior.

### Design

TODO: there's some content here that will need to be expanded, but this
feels like a good place to stop and check in with other contributors to make
sure the premise so far sounds good.

#### Extensions Specification

The APIs for resources which can be extended (`DataPlane`, `ControlPlane`,
`GatewayConfiguration`) will include an **optional** extensions field in their
specification, e.g.:

```golang
// DataPlaneSpec defines the desired state of DataPlane
type DataPlaneSpec struct {
	DataPlaneOptions `json:",inline"`

	// Extensions provide additional or replacement features for the DataPlane
	// resources to influence or enhance functionality.
	//
	// +optional
	Extensions []Extension `json:"extensions"`
}

// Extension corresponds to another resource in the Kubernetes cluster which
// defines extended behavior for a resource (e.g. DataPlane).
type Extension struct {
	// Group is the group of the extension resource.
	Group Group `json:"group"`

	// Kind is kind of the extension resource.
	Kind Kind `json:"kind"`

	// Name is the name of the extension resource.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name"`

	// Namespace is the namespace of the extension resource.
	//
	// For namespace-scoped resources if no Namespace is provided then the
	// namespace of the parent object MUST be used.
	//
	// This field MUST not be set when referring to cluster-scoped resources.
	//
	// +optional
	Namespace *Namespace `json:"namespace,omitempty"`
}
```

> **Note**: The extensions field to attach to other resources was chosen
> explicitly over having extension resources attach _themselves_ to the
> effected resource to make extensions more declarative, more discoverable and
> generally easier for end-users to reason about. This was influenced by our
> learnings from the [Gateway API][gwapi] project.

[gwapi]:https://github.com/kubernetes-sigs/gateway-api

#### Extension Logic

The logic for extensions will be loaded as **controller-runtime controllers**
which can replace standard (e.g. lives in the core OSS repository)
reconciliation logic. The standard controllers will include watch predicates
which identify controllers providing the extension and will disable themselves
to hand-off responsibility to the extension controllers.

> **Note**: a "replacement strategy" for controllers is the only mechanism
> we'll provide for now as it provides simplicitly.
>
> In the future we may consider other explicit "replacement strategies", e.g.:
>
>  * "hook" style attachment points

#### Extension Conflicts

Conflicts could arise from multiple extension attachments as the replacement
strategy means extension implementations compete for control. For now we'll keep
things simple by only allowing a single `Group` of APIs to be attached at one
time (e.g. `extensions.gateway-operator.konghq.com`).

Extension implementations **MUST** package all their relevant controllers
together.

> **Note**: in the future we can expand to some kind of cooperative mode
> where multiple extension groups are supported on a single resource, but at
> the time of writing there were no requirements driving this for the first
> iteration of extensions.

#### Extension API Guidelines

Extension resources are open-ended and can be _practically whatever the
developers want_ as it's up to them to provide the relevant logic. Despite the
custom nature of extension resources the following guidelines should be adhered
to during their development:

- To reduce confusion for users extensions should only be one level deep: don't
  add nested extension fields to extension resources.
- Extension APIs should be developed with multiple attachment in mind: expect
  that more than one resource may extend themselves by attaching the extension.
- Extension APIs should include status information to indicate which and how
  many resources are utilizing a particular extension. A cluster admin should
  be able to see at a glance
- Custom extension logic should include finalizers for Extension APIs. The
  premise being that it's more annoying for an admin to accidentally knock out
  production traffic due to a mishandled delete of an extension than it is to
  have them have to investigate and resolve the reason why the finalizer isn't
  letting them delete it.

> **Note**: Specification and status for Extension APIs is up to the developers
> ultimately, but the status sub-information which indicates how many and which
> resources are utilizing the extension will be provided by a common library.

### Test Plan

TODO

### Graduation Criteria

TODO

## Production Readiness

TODO

## Alternatives

TODO
