package versions

import (
	"fmt"

	"github.com/kong/gateway-operator/internal/consts"
)

const (
	// LatestClusterRoleVersion is the version of the ClusterRole that will be used for unversioned KIC
	LatestClusterRoleVersion = "2.9.3"
)

// RoleVersionsForKICVersions is a map that explicitly sets which ClusterRole version to use upon the KIC
// version. It is used by /hack/generators/kic-role-generator to generate the roles to be used by KIC.
// The file /internal/utils/kubernetes/resources/zz_generated_clusterrole_helpers.go is generated out of this map.
// This data follows the semver constraint syntax (see https://github.com/Masterminds/semver#basic-comparisons)
// to set the range of KIC versions to be associated with a specific role version.
//
// e.g., for KIC with a version lower than "2.4", but greater or equal to "2.3", version "2.3" of the role is used.
//
// Whenever the KIC ClusterRole is updated and released, that change should be reflected in this map,
// and the generator should be run again to produce the new ClusterRole files.
//
// e.g., when in the future new permissions will be granted to KIC, and that update will be included in
// the release 5.0, a new entry '">=5.0": "5.0"' should be added to this map, and the previous most
// updated entry should be limited to "<5.0".
var RoleVersionsForKICVersions = map[string]string{
	">=2.9":       "2.9.3",
	"<2.9, >=2.7": "2.7", // TODO: https://github.com/Kong/gateway-operator/issues/86
}

// supportedControlPlaneImages is the list of the supported ControlPlane images
var supportedControlPlaneImages = map[string]struct{}{
	fmt.Sprintf("%s:2.9.3", consts.DefaultControlPlaneBaseImage): {},
	fmt.Sprintf("%s:2.9", consts.DefaultControlPlaneBaseImage):   {},
	fmt.Sprintf("%s:2.7", consts.DefaultControlPlaneBaseImage):   {},
}

// IsControlPlaneSupported is a helper intended to validate the ControlPlane
// image support.
// The image is expected to follow the format <container-image-name>:<tag>
func IsControlPlaneSupported(image string) bool {
	_, ok := supportedControlPlaneImages[image]
	return ok
}
