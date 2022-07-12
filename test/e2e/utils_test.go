//go:build e2e_tests
// +build e2e_tests

package e2e

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

// namespace creates a namespace for a test case and cleans it up after the test finishes.
func namespace(t *testing.T) (*corev1.Namespace, func()) {
	namespaceForTestCase, err := clusters.GenerateNamespace(ctx, env.Cluster(), t.Name())
	exitOnErr(err)

	cleanup := func() {
		assert.NoError(t, clusters.CleanupGeneratedResources(ctx, env.Cluster(), t.Name()))
	}

	return namespaceForTestCase, cleanup
}

const gatewayOperatorImageKustomizationContents = `
images:
- name: ghcr.io/kong/gateway-operator:main
  newName: %v
  newTag: '%v'
`

func setOperatorImage() error {
	var image string
	if imageLoad != "" {
		image = imageLoad
	} else {
		image = imageOverride
	}

	if image == "" {
		fmt.Println("INFO: use default image")
		return nil
	}

	// TODO: deal with image names in format <host>:<port>/<repo>/<name>:[tag]
	// e.g localhost:32000/kong/gateway-operator:xxx
	parts := strings.Split(image, ":")
	if len(parts) != 2 {
		fmt.Printf("could not parse override image '%s', use default image\n", image)
		return nil
	}
	imageName := parts[0]
	imageTag := parts[1]

	fmt.Println("INFO: use custom image", image)

	// TODO: write back the kustomization file after test finished?
	buf, err := os.ReadFile("../../config/default/kustomization.yaml")
	if err != nil {
		return err
	}

	// append image contents to replace image
	fmt.Println("INFO: replacing image in kustomization file")
	appendImageKustomizationContents := []byte(fmt.Sprintf(gatewayOperatorImageKustomizationContents, imageName, imageTag))
	newBuf := append(buf, appendImageKustomizationContents...)
	return os.WriteFile("../../config/default/kustomization.yaml", newBuf, os.ModeAppend)
}
