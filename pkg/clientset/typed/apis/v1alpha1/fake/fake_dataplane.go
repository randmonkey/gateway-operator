/*
Copyright 2022 Kong Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/kong/gateway-operator/apis/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDataPlanes implements DataPlaneInterface
type FakeDataPlanes struct {
	Fake *FakeApisV1alpha1
	ns   string
}

var dataplanesResource = v1alpha1.SchemeGroupVersion.WithResource("dataplanes")

var dataplanesKind = v1alpha1.SchemeGroupVersion.WithKind("DataPlane")

// Get takes name of the dataPlane, and returns the corresponding dataPlane object, and an error if there is any.
func (c *FakeDataPlanes) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.DataPlane, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(dataplanesResource, c.ns, name), &v1alpha1.DataPlane{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DataPlane), err
}

// List takes label and field selectors, and returns the list of DataPlanes that match those selectors.
func (c *FakeDataPlanes) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.DataPlaneList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(dataplanesResource, dataplanesKind, c.ns, opts), &v1alpha1.DataPlaneList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.DataPlaneList{ListMeta: obj.(*v1alpha1.DataPlaneList).ListMeta}
	for _, item := range obj.(*v1alpha1.DataPlaneList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested dataPlanes.
func (c *FakeDataPlanes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(dataplanesResource, c.ns, opts))

}

// Create takes the representation of a dataPlane and creates it.  Returns the server's representation of the dataPlane, and an error, if there is any.
func (c *FakeDataPlanes) Create(ctx context.Context, dataPlane *v1alpha1.DataPlane, opts v1.CreateOptions) (result *v1alpha1.DataPlane, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(dataplanesResource, c.ns, dataPlane), &v1alpha1.DataPlane{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DataPlane), err
}

// Update takes the representation of a dataPlane and updates it. Returns the server's representation of the dataPlane, and an error, if there is any.
func (c *FakeDataPlanes) Update(ctx context.Context, dataPlane *v1alpha1.DataPlane, opts v1.UpdateOptions) (result *v1alpha1.DataPlane, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(dataplanesResource, c.ns, dataPlane), &v1alpha1.DataPlane{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DataPlane), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeDataPlanes) UpdateStatus(ctx context.Context, dataPlane *v1alpha1.DataPlane, opts v1.UpdateOptions) (*v1alpha1.DataPlane, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(dataplanesResource, "status", c.ns, dataPlane), &v1alpha1.DataPlane{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DataPlane), err
}

// Delete takes name of the dataPlane and deletes it. Returns an error if one occurs.
func (c *FakeDataPlanes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(dataplanesResource, c.ns, name, opts), &v1alpha1.DataPlane{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDataPlanes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(dataplanesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.DataPlaneList{})
	return err
}

// Patch applies the patch and returns the patched dataPlane.
func (c *FakeDataPlanes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.DataPlane, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(dataplanesResource, c.ns, name, pt, data, subresources...), &v1alpha1.DataPlane{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DataPlane), err
}
