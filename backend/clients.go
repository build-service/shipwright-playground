package main

import (
	"container/heap"
	"fmt"
	"log"

	shipwrightClient "github.com/shipwright-io/build/pkg/client/clientset/versioned"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getRestConfigFromBytes(kubeConfigBytes []byte) (*rest.Config, error) {
	kubeClientConfig, err := clientcmd.NewClientConfigFromBytes(kubeConfigBytes)
	if err != nil {
		return nil, fmt.Errorf("Unable to create client config from kubeconfig bytes, error: %v", err)
	}

	restConfig, err := kubeClientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Unable to create rest config, error: %v", err)
	}

	return restConfig, nil
}

func getk8sClient(restConfig *rest.Config) (*kubernetes.Clientset, error) {
	k8sClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("Unable to create k8s client object, error: %v", err)
	}

	return k8sClient, nil
}

func getShipwrightClient(restConfig *rest.Config) (*shipwrightClient.Clientset, error) {
	shipwrightClient, err := shipwrightClient.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("Unable to create shipwright client object, error: %v", err)
	}

	return shipwrightClient, nil
}

func getDiscoveryClient(restConfig *rest.Config) (*discovery.DiscoveryClient, error) {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("Unable to create discovery client object, error: %v", err)
	}

	return discoveryClient, nil
}

func getDynamicClient(restConfig *rest.Config) (dynamic.Interface, error) {
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("Unable to create dynamic client object, error: %v", err)
	}

	return dynamicClient, nil
}

func getNewClusterRestConfig() (*rest.Config, error) {
	if len(clusterPool) == 0 {
		return nil, fmt.Errorf("Unable to create k8s client, no clusters available in cluster pool")
	}

	var newClusterRestConfig *rest.Config
	for len(clusterPool) > 0 {
		newCluster := heap.Pop(&clusterPool).(*cluster)
		newClusterRestConfig, err := getRestConfigFromBytes(newCluster.KubeConfigContents)
		if err != nil {
			log.Printf("error retrieving new cluster rest config: %v, Trying different cluster", err.Error())
			continue
		} else if newClusterRestConfig != nil {
			config = newClusterRestConfig
		}

	}

	if newClusterRestConfig == nil {
		return nil, fmt.Errorf("Unable to create k8s client, no valid clusters available in cluster pool")
	}

	return newClusterRestConfig, nil
}
