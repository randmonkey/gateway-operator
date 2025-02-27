package kubernetes

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// ConditionType literal that defines the different types of condition
type ConditionType string

// CoditionReason literal to enumerate a specific condition reason
type ConditionReason string

const (

	// ReadyType indicates if the resource has all the dependent conditions Ready
	ReadyType ConditionType = "Ready"

	// ProgrammedType condition indicates whether a Gateway has generated some
	// configuration that is assumed to be ready soon in the underlying data
	// plane.
	ProgrammedType ConditionType = ConditionType(gatewayv1beta1.GatewayConditionProgrammed)

	// PendingReason is a Reason for Programmed condition.
	PendingReason ConditionReason = ConditionReason(gatewayv1beta1.GatewayReasonPending)

	// DependenciesNotReadyReason is a generic reason describing that the other Conditions are not true
	DependenciesNotReadyReason ConditionReason = "DependenciesNotReady"

	// ResourceReadyReason indicates the resource is ready
	ResourceReadyReason ConditionReason = ConditionReason("Ready")

	// WaitingToBecomeReadyReason generic message for dependent resources waiting to be ready
	WaitingToBecomeReadyReason ConditionReason = "WaitingToBecomeReady"

	// ResourceCreatedOrUpdatedReason generic message for missing or outdated resources
	ResourceCreatedOrUpdatedReason ConditionReason = "ResourceCreatedOrUpdated"

	// UnableToProvisionReason generic message for unexpected errors
	UnableToProvisionReason ConditionReason = "UnableToProvision"

	// DependenciesNotReadyMessage indicates the other conditions are not yet ready
	DependenciesNotReadyMessage = "There are other conditions that are not yet ready"

	// WaitingToBecomeReadyMessage indicates the target resource is not ready
	WaitingToBecomeReadyMessage = "Waiting for the resource to become ready"

	// ResourceCreatedMessage indicates a missing resource was provisioned
	ResourceCreatedMessage = "Resource has been created"

	// ResourceUpdatedMessage indicates a resource was updated
	ResourceUpdatedMessage = "Resource has been updated"
)

// ConditionsAware represents a CRD type that has been enabled with metav1.Conditions,
// it can then benefit of a series of utility methods.
type ConditionsAware interface {
	GetConditions() []metav1.Condition
	SetConditions(conditions []metav1.Condition)
}

// SetCondition sets a new condition to the provided resource
func SetCondition(condition metav1.Condition, resource ConditionsAware) {
	conditions := resource.GetConditions()
	newConditions := make([]metav1.Condition, 0, len(conditions))

	for i := 0; i < len(conditions); i++ {
		if conditions[i].Type != condition.Type {
			newConditions = append(newConditions, conditions[i])
		}
	}

	newConditions = append(newConditions, condition)
	resource.SetConditions(newConditions)
}

// GetCondition returns the condition with the given type, if it exists. If the condition does not exists it returns false.
func GetCondition(cType ConditionType, resource ConditionsAware) (metav1.Condition, bool) {
	for _, condition := range resource.GetConditions() {
		if condition.Type == string(cType) {
			return condition, true
		}
	}
	return metav1.Condition{}, false
}

// IsValidCondition returns a true value whether the condition is ConditionTrue, false otherwise
func IsValidCondition(cType ConditionType, resource ConditionsAware) bool {
	for _, condition := range resource.GetConditions() {
		if condition.Type == string(cType) {
			return condition.Status == metav1.ConditionTrue
		}
	}
	return false
}

// InitReady initializes the Ready status to False
func InitReady(resource ConditionsAware) {
	SetCondition(NewCondition(ReadyType, metav1.ConditionFalse, DependenciesNotReadyReason, DependenciesNotReadyMessage), resource)
}

// InitProgrammed initializes the Programmed status to False
func InitProgrammed(resource ConditionsAware) {
	SetCondition(NewCondition(ProgrammedType, metav1.ConditionFalse, PendingReason, DependenciesNotReadyMessage), resource)
}

// SetReady evaluates all the existing conditions and sets the Ready status accordingly
func SetReady(resource ConditionsAware, generation int64) {
	ready := metav1.Condition{
		Type:               string(ReadyType),
		LastTransitionTime: metav1.Now(),
		ObservedGeneration: generation,
	}

	if areAllConditionsHaveTrueStatus(resource) {
		ready.Status = metav1.ConditionTrue
		ready.Reason = string(ResourceReadyReason)
	} else {
		ready.Status = metav1.ConditionFalse
		ready.Reason = string(DependenciesNotReadyReason)
		ready.Message = DependenciesNotReadyMessage
	}
	SetCondition(ready, resource)
}

// SetProgrammed evaluates all the existing conditions and sets the Programmed status accordingly
func SetProgrammed(resource ConditionsAware, generation int64) {
	ready := metav1.Condition{
		Type:               string(gatewayv1beta1.GatewayConditionProgrammed),
		LastTransitionTime: metav1.Now(),
		ObservedGeneration: generation,
	}

	if areAllConditionsHaveTrueStatus(resource) {
		ready.Status = metav1.ConditionTrue
		ready.Reason = string(gatewayv1beta1.GatewayReasonProgrammed)
	} else {
		ready.Status = metav1.ConditionFalse
		ready.Reason = string(DependenciesNotReadyReason)
		ready.Message = DependenciesNotReadyMessage
	}
	SetCondition(ready, resource)
}

func areAllConditionsHaveTrueStatus(resource ConditionsAware) bool {
	for _, condition := range resource.GetConditions() {
		if condition.Type == string(gatewayv1beta1.GatewayConditionProgrammed) {
			continue
		}
		if condition.Type != string(ReadyType) && condition.Status != metav1.ConditionTrue {
			return false
		}
	}
	return true
}

// IsReady evaluates whether a resource is in Ready state, meaning
// that all its conditions are in the True state.
func IsReady(resource ConditionsAware) bool {
	for _, condition := range resource.GetConditions() {
		if condition.Type == string(ReadyType) {
			return condition.Status == metav1.ConditionTrue
		}
	}
	return false
}

// IsProgrammed evaluates whether a resource is in Programmed state.
func IsProgrammed(resource ConditionsAware) bool {
	for _, condition := range resource.GetConditions() {
		if condition.Type == string(ProgrammedType) {
			return condition.Status == metav1.ConditionTrue
		}
	}
	return false
}

// NewCondition convenience method for creating conditions
func NewCondition(cType ConditionType, status metav1.ConditionStatus, reason ConditionReason, message string) metav1.Condition {
	return metav1.Condition{
		Type:               string(cType),
		Reason:             string(reason),
		Message:            message,
		LastTransitionTime: metav1.Now(),
		Status:             status,
	}
}

// NewConditionWithGeneration convenience method for creating conditions with ObservedGeneration set.
func NewConditionWithGeneration(cType ConditionType, status metav1.ConditionStatus, reason ConditionReason, message string, observedGeneration int64) metav1.Condition {
	c := NewCondition(cType, status, reason, message)
	c.ObservedGeneration = observedGeneration
	return c
}

// NeedsUpdate retrieves the persisted state and compares all the conditions
// to decide whether the status must be updated or not
func NeedsUpdate(current, updated ConditionsAware) bool {
	if len(current.GetConditions()) != len(updated.GetConditions()) {
		return true
	}

	for _, c := range current.GetConditions() {
		u, exists := GetCondition(ConditionType(c.Type), updated)
		if !exists {
			return true
		}
		if u.Reason != c.Reason || u.Message != c.Message || u.Status != c.Status {
			return true
		}
	}
	return false
}
