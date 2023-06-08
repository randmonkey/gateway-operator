// This file is generated by /hack/generators/kic-clusterrole-generator. DO NOT EDIT.

package resources

import (
	"fmt"

	"github.com/Masterminds/semver"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/kong/gateway-operator/internal/consts"
	"github.com/kong/gateway-operator/internal/utils/kubernetes/resources/clusterroles"
	"github.com/kong/gateway-operator/internal/versions"
)

// -----------------------------------------------------------------------------
// ClusterRole generator helper
// -----------------------------------------------------------------------------

// GenerateNewClusterRoleForControlPlane is a helper function that extract
// the version from the tag, and returns the ClusterRole with all the needed
// permissions.
func GenerateNewClusterRoleForControlPlane(controlplaneName string, image, tag *string) (*rbacv1.ClusterRole, error) {
	versionToUse := consts.DefaultControlPlaneTag
	imageToUse := consts.DefaultControlPlaneImage
	var constraint *semver.Constraints

	if image != nil && *image != "" && tag != nil && *tag != "" {
		askedImage := fmt.Sprintf("%s:%s", *image, *tag)
		supported, err := versions.IsControlPlaneImageVersionSupported(askedImage)
		if err != nil {
			return nil, err
		}
		if supported {
			imageToUse = askedImage
			versionToUse = *tag
		}
	}

	semVersion, err := semver.NewVersion(versionToUse)
	if err != nil {
		return nil, err
	}

	constraint, err = semver.NewConstraint("<2.9, >=2.7")
	if err != nil {
		return nil, err
	}
	if constraint.Check(semVersion) {
		return clusterroles.GenerateNewClusterRoleForControlPlane_lt2_9_ge2_7(controlplaneName), nil
	}

	constraint, err = semver.NewConstraint(">=2.9")
	if err != nil {
		return nil, err
	}
	if constraint.Check(semVersion) {
		return clusterroles.GenerateNewClusterRoleForControlPlane_ge2_9(controlplaneName), nil
	}

	return nil, fmt.Errorf("version %s not supported", imageToUse)
}
