---
title: API Extensions
status: implementable
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

The overarching design for the extensions mechanism is founded on adding
extension fields to KGO's APIs and enabling extension developers to provide
replacement logic for the controllers for those APIs.

In the following sub-sections we will go into details about how the extensions
are meant to be implemented.

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
	Extensions []ExtensionRef `json:"extensions"`
}

// ExtensionRef corresponds to another resource in the Kubernetes cluster which
// defines extended behavior for a resource (e.g. DataPlane).
type ExtensionRef struct {
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

Replacement is up to the extension provider: watch predicates will be in place
that can be used to disable the standard controllers for a resource (e.g. for
`DataPlane`) and the extension provider can then load their own controller(s)
to replace that logic, with all of the standard logic being available to them
"off the shelf" via a standard library. Extensions can potentially replace a
single controller with one extension controller, or many.

> **Note**: a "replacement strategy" for controllers is the only mechanism
> we'll provide for now, though others were considered. See the [alternatives
> considered](#alternatives-considered) section for more details.

##### Loaded Extensions & Status

Any resource that implements extensions **MUST** implement our standard status
field:

```golang
type ExtensionStatus struct {
    Loaded []ExtentionRef
}
```

The `Loaded` field indicates to any controller that would act on the resource
that extensions were previously loaded by the resources referenced in the
field.

> **Note**: the standard implementation will refuse to operate on any resource
> which has extensions loaded on it, even if they've been removed in the
> specification. To migrate back to standard from an extended implementation,
> the standard implementation expects the extensions to be "drained" and the
> status to show an empty `Loaded` list.

##### Migrations

Extension implementations **MUST** provide and document a mechanism to migrate
back to the standard KGO implementation for resources such as `DataPlane`. For
some implementations this may be as simple as draining the extensions (e.g.
remove them from the spec, and wait for the status to show no more `Loaded`
extensions).

> **Note**: Extension implementations are responsible for providing testing
> to ensure that migrating back to the standard implementation works properly
> and maintain integrity when control is transferred back to standard mode.

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

For any resource which will have extension points made available (e.g.
`DataPlane`, `ControlPlane`) an example extension will be provided that
provides some modified behavior.

Integration tests in `test/integration` will be responsible for exercising
these example extensions with the following high level testing plan:

- verify running the controller manager _without_ the extension controller
  loaded
- verify creating several resources (e.g. `DataPlane`) with the standard
  implementation
- re-deploy the controller manager with the example extensions enabled
- verify that the existing resources with no extensions on them continue to
  exhibit appropriate behavior
- verify attaching extensions to some of the test resources
- verify that resources with extensions exhibit extended behavior
- verify that resources with no extensions continue to exhibit standard behavior

This testing will help us to reduce the risk of breakage or regressions in
future releases.

### Graduation Criteria

- [ ] extension points are added to `DataPlane`.
- [ ] extension points are added to `ControlPlane`.
- [ ] extension predicates are added to relevant controllers to filter out
      resources with extensions attached via specification, or loaded in the
      resource status.
- [ ] example extensions are created for `DataPlane`
- [ ] example extensions are created for `ControlPlane`
- [ ] integration tests are added for the example extensions

## Production Readiness

Production readiness for this feature will be marked by a combination of the
above graduation criteria being entirely resolved _and_ there being at least
one extension implementation actively maintained and in use for a period of
several months (likely this will be done within Kong itself, though the door
technically remains open for other downstream implementations).

## Alternatives

Some alternative mechanisms for extensions were considered by the KGO
maintainers, but ultimately were rejected:

### Controller Hooks

In brainstorming among contributors our first concept for enabling extensions
was hook points where extension logic would be loaded into the standard
controllers directly. This had both maintenance and composability problems
compared to the "replacement with libraries" strategy. For our immediate goals
and motivations this was not optimal, and there are no plans to revisit this
unless some future requirements arise that make it necessary.

### Forking

We could have suggested maintaining a fork as the "mechanism" by which
extended versions of the KGO could be developed, but this filled all of us with
rage and tears and so we noped out pretty hard on this one.
