package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	operatorv1alpha1 "github.com/kong/gateway-operator/apis/v1alpha1"
	"github.com/kong/gateway-operator/internal/consts"
	gatewayutils "github.com/kong/gateway-operator/internal/utils/gateway"
	k8sutils "github.com/kong/gateway-operator/internal/utils/kubernetes"
)

// MustListDataPlaneDeployments is a helper function for tests that
// conveniently lists all deployments managed by a given dataplane.
func MustListDataPlaneDeployments(t *testing.T, ctx context.Context, dataplane *operatorv1alpha1.DataPlane, clients K8sClients) []appsv1.Deployment {
	deployments, err := k8sutils.ListDeploymentsForOwner(
		ctx,
		clients.MgrClient,
		consts.GatewayOperatorControlledLabel,
		consts.DataPlaneManagedLabelValue,
		dataplane.Namespace,
		dataplane.UID,
	)
	require.NoError(t, err)
	return deployments
}

// MustListControlPlaneDeployments is a helper function for tests that
// conveniently lists all deployments managed by a given controlplane.
func MustListControlPlaneDeployments(t *testing.T, ctx context.Context, controlplane *operatorv1alpha1.ControlPlane, clients K8sClients) []appsv1.Deployment {
	deployments, err := k8sutils.ListDeploymentsForOwner(
		ctx,
		clients.MgrClient,
		consts.GatewayOperatorControlledLabel,
		consts.ControlPlaneManagedLabelValue,
		controlplane.Namespace,
		controlplane.UID,
	)
	require.NoError(t, err)
	return deployments
}

// MustListControlPlaneClusterRoles is a helper function for tests that
// conveniently lists all clusterroles owned by a given controlplane.
func MustListControlPlaneClusterRoles(t *testing.T, ctx context.Context, controlplane *operatorv1alpha1.ControlPlane, clients K8sClients) []rbacv1.ClusterRole {
	clusterRoles, err := k8sutils.ListClusterRolesForOwner(
		ctx,
		clients.MgrClient,
		consts.GatewayOperatorControlledLabel,
		consts.ControlPlaneManagedLabelValue,
		controlplane.UID,
	)
	require.NoError(t, err)
	return clusterRoles
}

// MustListControlPlaneClusterRoleBindings is a helper function for tests that
// conveniently lists all clusterrolebindings owned by a given controlplane.
func MustListControlPlaneClusterRoleBindings(t *testing.T, ctx context.Context, controlplane *operatorv1alpha1.ControlPlane, clients K8sClients) []rbacv1.ClusterRoleBinding {
	clusterRoleBindings, err := k8sutils.ListClusterRoleBindingsForOwner(
		ctx,
		clients.MgrClient,
		consts.GatewayOperatorControlledLabel,
		consts.ControlPlaneManagedLabelValue,
		controlplane.UID,
	)
	require.NoError(t, err)
	return clusterRoleBindings
}

// MustListControlPlanesForGateway is a helper function for tests that
// conveniently lists all controlplanes managed by a given gateway.
func MustListControlPlanesForGateway(t *testing.T, ctx context.Context, gateway *gatewayv1alpha2.Gateway, clients K8sClients) []operatorv1alpha1.ControlPlane {
	controlPlanes, err := gatewayutils.ListControlPlanesForGateway(ctx, clients.MgrClient, gateway)
	require.NoError(t, err)
	return controlPlanes
}

// MustListNetworkPoliciesForGateway is a helper function for tests that
// conveniently lists all NetworkPolicies managed by a given gateway.
func MustListNetworkPoliciesForGateway(t *testing.T, ctx context.Context, gateway *gatewayv1alpha2.Gateway, clients K8sClients) []networkingv1.NetworkPolicy {
	networkPolicies, err := gatewayutils.ListNetworkPoliciesForGateway(ctx, clients.MgrClient, gateway)
	require.NoError(t, err)
	return networkPolicies
}

// MustListServices is a helper function for tests that
// conveniently lists all services managed by a given dataplane.
func MustListDataPlaneServices(t *testing.T, ctx context.Context, dataplane *operatorv1alpha1.DataPlane, mgrClient client.Client) []corev1.Service {
	services, err := k8sutils.ListServicesForOwner(
		ctx,
		mgrClient,
		consts.GatewayOperatorControlledLabel,
		consts.DataPlaneManagedLabelValue,
		dataplane.Namespace,
		dataplane.UID,
	)
	require.NoError(t, err)
	return services
}

// MustListDataPlanesForGateway is a helper function for tests that
// conveniently lists all dataplanes managed by a given gateway.
func MustListDataPlanesForGateway(t *testing.T, ctx context.Context, gateway *gatewayv1alpha2.Gateway, clients K8sClients) []operatorv1alpha1.DataPlane {
	dataplanes, err := gatewayutils.ListDataPlanesForGateway(ctx, clients.MgrClient, gateway)
	require.NoError(t, err)
	return dataplanes
}

// MustGetGateway is a helper function for tests that conveniently gets a gateway by name.
// It will fail the test if getting the gateway fails.
func MustGetGateway(t *testing.T, ctx context.Context, gatewayNSN types.NamespacedName, clients K8sClients) *gatewayv1alpha2.Gateway {
	gateways := clients.GatewayClient.GatewayV1alpha2().Gateways(gatewayNSN.Namespace)
	gateway, err := gateways.Get(ctx, gatewayNSN.Name, metav1.GetOptions{})
	require.NoError(t, err)
	return gateway
}
