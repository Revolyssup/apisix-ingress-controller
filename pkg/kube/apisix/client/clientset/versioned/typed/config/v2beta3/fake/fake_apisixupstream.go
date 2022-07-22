// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v2beta3 "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2beta3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeApisixUpstreams implements ApisixUpstreamInterface
type FakeApisixUpstreams struct {
	Fake *FakeApisixV2beta3
	ns   string
}

var apisixupstreamsResource = schema.GroupVersionResource{Group: "apisix.apache.org", Version: "v2beta3", Resource: "apisixupstreams"}

var apisixupstreamsKind = schema.GroupVersionKind{Group: "apisix.apache.org", Version: "v2beta3", Kind: "ApisixUpstream"}

// Get takes name of the apisixUpstream, and returns the corresponding apisixUpstream object, and an error if there is any.
func (c *FakeApisixUpstreams) Get(ctx context.Context, name string, options v1.GetOptions) (result *v2beta3.ApisixUpstream, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(apisixupstreamsResource, c.ns, name), &v2beta3.ApisixUpstream{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2beta3.ApisixUpstream), err
}

// List takes label and field selectors, and returns the list of ApisixUpstreams that match those selectors.
func (c *FakeApisixUpstreams) List(ctx context.Context, opts v1.ListOptions) (result *v2beta3.ApisixUpstreamList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(apisixupstreamsResource, apisixupstreamsKind, c.ns, opts), &v2beta3.ApisixUpstreamList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v2beta3.ApisixUpstreamList{ListMeta: obj.(*v2beta3.ApisixUpstreamList).ListMeta}
	for _, item := range obj.(*v2beta3.ApisixUpstreamList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested apisixUpstreams.
func (c *FakeApisixUpstreams) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(apisixupstreamsResource, c.ns, opts))

}

// Create takes the representation of a apisixUpstream and creates it.  Returns the server's representation of the apisixUpstream, and an error, if there is any.
func (c *FakeApisixUpstreams) Create(ctx context.Context, apisixUpstream *v2beta3.ApisixUpstream, opts v1.CreateOptions) (result *v2beta3.ApisixUpstream, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(apisixupstreamsResource, c.ns, apisixUpstream), &v2beta3.ApisixUpstream{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2beta3.ApisixUpstream), err
}

// Update takes the representation of a apisixUpstream and updates it. Returns the server's representation of the apisixUpstream, and an error, if there is any.
func (c *FakeApisixUpstreams) Update(ctx context.Context, apisixUpstream *v2beta3.ApisixUpstream, opts v1.UpdateOptions) (result *v2beta3.ApisixUpstream, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(apisixupstreamsResource, c.ns, apisixUpstream), &v2beta3.ApisixUpstream{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2beta3.ApisixUpstream), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeApisixUpstreams) UpdateStatus(ctx context.Context, apisixUpstream *v2beta3.ApisixUpstream, opts v1.UpdateOptions) (*v2beta3.ApisixUpstream, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(apisixupstreamsResource, "status", c.ns, apisixUpstream), &v2beta3.ApisixUpstream{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2beta3.ApisixUpstream), err
}

// Delete takes name of the apisixUpstream and deletes it. Returns an error if one occurs.
func (c *FakeApisixUpstreams) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(apisixupstreamsResource, c.ns, name, opts), &v2beta3.ApisixUpstream{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeApisixUpstreams) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(apisixupstreamsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v2beta3.ApisixUpstreamList{})
	return err
}

// Patch applies the patch and returns the patched apisixUpstream.
func (c *FakeApisixUpstreams) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2beta3.ApisixUpstream, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(apisixupstreamsResource, c.ns, name, pt, data, subresources...), &v2beta3.ApisixUpstream{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2beta3.ApisixUpstream), err
}
