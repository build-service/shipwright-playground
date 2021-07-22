// Copyright The Shipwright Contributors
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/shipwright-io/build/pkg/apis/build/v1alpha1"
	scheme "github.com/shipwright-io/build/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ClusterBuildStrategiesGetter has a method to return a ClusterBuildStrategyInterface.
// A group's client should implement this interface.
type ClusterBuildStrategiesGetter interface {
	ClusterBuildStrategies() ClusterBuildStrategyInterface
}

// ClusterBuildStrategyInterface has methods to work with ClusterBuildStrategy resources.
type ClusterBuildStrategyInterface interface {
	Create(ctx context.Context, clusterBuildStrategy *v1alpha1.ClusterBuildStrategy, opts v1.CreateOptions) (*v1alpha1.ClusterBuildStrategy, error)
	Update(ctx context.Context, clusterBuildStrategy *v1alpha1.ClusterBuildStrategy, opts v1.UpdateOptions) (*v1alpha1.ClusterBuildStrategy, error)
	UpdateStatus(ctx context.Context, clusterBuildStrategy *v1alpha1.ClusterBuildStrategy, opts v1.UpdateOptions) (*v1alpha1.ClusterBuildStrategy, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.ClusterBuildStrategy, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.ClusterBuildStrategyList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ClusterBuildStrategy, err error)
	ClusterBuildStrategyExpansion
}

// clusterBuildStrategies implements ClusterBuildStrategyInterface
type clusterBuildStrategies struct {
	client rest.Interface
}

// newClusterBuildStrategies returns a ClusterBuildStrategies
func newClusterBuildStrategies(c *ShipwrightV1alpha1Client) *clusterBuildStrategies {
	return &clusterBuildStrategies{
		client: c.RESTClient(),
	}
}

// Get takes name of the clusterBuildStrategy, and returns the corresponding clusterBuildStrategy object, and an error if there is any.
func (c *clusterBuildStrategies) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ClusterBuildStrategy, err error) {
	result = &v1alpha1.ClusterBuildStrategy{}
	err = c.client.Get().
		Resource("clusterbuildstrategies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterBuildStrategies that match those selectors.
func (c *clusterBuildStrategies) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ClusterBuildStrategyList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ClusterBuildStrategyList{}
	err = c.client.Get().
		Resource("clusterbuildstrategies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested clusterBuildStrategies.
func (c *clusterBuildStrategies) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("clusterbuildstrategies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a clusterBuildStrategy and creates it.  Returns the server's representation of the clusterBuildStrategy, and an error, if there is any.
func (c *clusterBuildStrategies) Create(ctx context.Context, clusterBuildStrategy *v1alpha1.ClusterBuildStrategy, opts v1.CreateOptions) (result *v1alpha1.ClusterBuildStrategy, err error) {
	result = &v1alpha1.ClusterBuildStrategy{}
	err = c.client.Post().
		Resource("clusterbuildstrategies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterBuildStrategy).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a clusterBuildStrategy and updates it. Returns the server's representation of the clusterBuildStrategy, and an error, if there is any.
func (c *clusterBuildStrategies) Update(ctx context.Context, clusterBuildStrategy *v1alpha1.ClusterBuildStrategy, opts v1.UpdateOptions) (result *v1alpha1.ClusterBuildStrategy, err error) {
	result = &v1alpha1.ClusterBuildStrategy{}
	err = c.client.Put().
		Resource("clusterbuildstrategies").
		Name(clusterBuildStrategy.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterBuildStrategy).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *clusterBuildStrategies) UpdateStatus(ctx context.Context, clusterBuildStrategy *v1alpha1.ClusterBuildStrategy, opts v1.UpdateOptions) (result *v1alpha1.ClusterBuildStrategy, err error) {
	result = &v1alpha1.ClusterBuildStrategy{}
	err = c.client.Put().
		Resource("clusterbuildstrategies").
		Name(clusterBuildStrategy.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterBuildStrategy).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the clusterBuildStrategy and deletes it. Returns an error if one occurs.
func (c *clusterBuildStrategies) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("clusterbuildstrategies").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *clusterBuildStrategies) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("clusterbuildstrategies").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched clusterBuildStrategy.
func (c *clusterBuildStrategies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ClusterBuildStrategy, err error) {
	result = &v1alpha1.ClusterBuildStrategy{}
	err = c.client.Patch(pt).
		Resource("clusterbuildstrategies").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
