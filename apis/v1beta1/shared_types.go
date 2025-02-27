package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
)

// DeploymentOptions is a shared type used on objects to indicate that their
// configuration results in a Deployment which is managed by the Operator and
// includes options for managing Deployments such as the the number of replicas
// or pod options like container image and resource requirements.
// version, as well as Env variable overrides.
type DeploymentOptions struct {
	// Replicas describes the number of desired pods.
	// This is a pointer to distinguish between explicit zero and not specified.
	// This only affects the DataPlane deployments for now, for more details on
	// ControlPlane scaling please see https://github.com/Kong/gateway-operator/issues/736.
	//
	// +optional
	// +kubebuilder:default=1
	Replicas *int32 `json:"replicas,omitempty"`

	// PodTemplateSpec defines PodTemplateSpec for Deployment's pods.
	//
	// +optional
	PodTemplateSpec *corev1.PodTemplateSpec `json:"podTemplateSpec,omitempty"`
}

// Rollout defines options for rollouts.
type Rollout struct {
	// Strategy contains the deployment strategy for rollout.
	Strategy RolloutStrategy `json:"strategy"`
}

// RolloutStrategy holds the rollout strategy options.
type RolloutStrategy struct {
	// BlueGreen holds the options specific for Blue Green Deployments.
	//
	// +optional
	BlueGreen *BlueGreenStrategy `json:"blueGreen,omitempty"`
}

// BlueGreenStrategy defines the Blue Green deployment strategy.
type BlueGreenStrategy struct {
	// Promotion defines how the operator handles promotion of resources.
	Promotion Promotion `json:"promotion"`
}

type Promotion struct {
	// Strategy indicates how you want the operator to handle the promotion of
	// the preview (green) resources (Deployments and Services) after all workflows
	// and tests succeed, OR if you even want it to break before performing
	// the promotion to allow manual inspection.
	//
	// +kubebuilder:validation:Enum=AutomaticPromotion;BreakBeforePromotion
	// +kubebuilder:default=BreakBeforePromotion
	Strategy PromotionStrategy `json:"strategy"`
}

type PromotionStrategy string

const (
	// AutomaticPromotion indicates that once all workflows and tests have completed successfully,
	// the new resources should be promoted and replace the previous resources.
	AutomaticPromotion PromotionStrategy = "AutomaticPromotion"

	// BreakBeforePromotion is the same as AutomaticPromotion but with an added breakpoint
	// to enable manual inspection.
	// The user must indicate manually when they want the promotion to continue.
	// TODO: finalizer/annotation?
	BreakBeforePromotion PromotionStrategy = "BreakBeforePromotion"
)

// GatewayConfigurationTargetKind is an object kind that can be targeted for
// GatewayConfiguration attachment.
type GatewayConfigurationTargetKind string

const (
	// GatewayConfigurationTargetKindGateway is a target kind which indicates
	// that a Gateway resource is the target.
	GatewayConfigurationTargetKindGateway GatewayConfigurationTargetKind = "Gateway"

	// GatewayConfigurationTargetKindGatewayClass is a target kind which indicates
	// that a GatewayClass resource is the target.
	GatewayConfigurationTargetKindGatewayClass GatewayConfigurationTargetKind = "GatewayClass"
)
