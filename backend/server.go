package main

import (
	"container/heap"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	shipwright "github.com/shipwright-io/build/pkg/apis/build/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	contextDir  string
	clusterPool PriorityQueue
)

const (
	username         = "jrao"
	password         = "Quayjrao1!"
	quayServer       = "quay.io"
	email            = "jrao@redhat.com"
	secretName       = "image-registry-secret"
	serverPort       = 8085
	imageRegistryOrg = "shipwrightplayground"
)

func buildRequestHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	repoURL := r.FormValue("repo-url")
	contextDir = r.FormValue("context-dir")
	gitHubOrg, repoName := parseGitRepoURL(repoURL)

	k8sClient, _ := getk8sClient(currentClusterRestConfig)
	ok, err := testClusterConnection(k8sClient)
	if !ok || k8sClient == nil {
		_, err := getNewClusterRestConfig()
		if err != nil {
			fmt.Fprintf(w, "Unable to establish connection with an active cluster: %v", err)
			log.Fatalf("Unable to establish connection with an active cluster: %v", err)
		}
	}

	k8sClient, _ = getk8sClient(currentClusterRestConfig)
	buildClient, _ := getShipwrightClient(currentClusterRestConfig)
	err = prepareCluster(k8sClient)
	if err != nil {
		fmt.Fprintf(w, "Error preparing cluster with pre-requisite resources: %v", err)
		log.Fatalf("Error preparing cluster with pre-requisite resources: %v", err)
	}

	// Create Build on cluster
	if _, err := buildClient.ShipwrightV1alpha1().Builds("default").Get(context.TODO(), fmt.Sprintf("%v-build", repoName), v1.GetOptions{}); err == nil {
		err := buildClient.ShipwrightV1alpha1().Builds("default").Delete(context.TODO(), fmt.Sprintf("%v-build", repoName), v1.DeleteOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
	buildObj, err := createBuild(quayServer, gitHubOrg, repoURL, imageRegistryOrg, repoName, secretName, contextDir)
	_, err = buildClient.ShipwrightV1alpha1().Builds("default").Create(context.TODO(), buildObj, v1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Create BuildRun on cluster
	if _, err := buildClient.ShipwrightV1alpha1().BuildRuns("default").Get(context.TODO(), fmt.Sprintf("%v-buildrun", repoName), v1.GetOptions{}); err == nil {
		err := buildClient.ShipwrightV1alpha1().BuildRuns("default").Delete(context.TODO(), fmt.Sprintf("%v-buildrun", repoName), v1.DeleteOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
	buildRunObj := createBuildRun(repoName)
	_, err = buildClient.ShipwrightV1alpha1().BuildRuns("default").Create(context.TODO(), buildRunObj, v1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}

	err = wait.Poll(time.Second*4, time.Second*180, func() (bool, error) {
		existingBuildRun, err := buildClient.ShipwrightV1alpha1().BuildRuns("default").Get(context.TODO(), buildRunObj.Name, v1.GetOptions{})
		if err != nil {
			// log.Fatalf("Error retrieving buildrun: %v", err)
			return false, err
		}
		for _, condition := range existingBuildRun.Status.Conditions {
			if condition.Type == shipwright.Succeeded && condition.Status == corev1.ConditionUnknown {
				fmt.Fprintf(w, "Building...")
				return false, nil
			} else if condition.Type == shipwright.Succeeded && condition.Status == corev1.ConditionFalse {
				fmt.Fprintf(w, "Build failed")
				return true, nil
			} else if condition.Type == shipwright.Succeeded && condition.Status == corev1.ConditionTrue {
				fmt.Fprintf(w, "Build successful!")
				break
			}
		}

		return true, nil
	})

}

func clusterAdditionHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	kubeConfigURL := r.FormValue("kubeconfigURL")
	expirationHours := r.FormValue("expires")
	expires, _ := strconv.Atoi(expirationHours)

	kubeConfigFile, err := http.Get(kubeConfigURL)
	if err != nil {
		fmt.Printf("Unable to retrieve manifest from URL: %v", err)
	}
	defer kubeConfigFile.Body.Close()

	kubeConfigBytes, err := ioutil.ReadAll(kubeConfigFile.Body)
	if err != nil {
		fmt.Printf("Unable to retrieve manifest bytes: %v", err)
	}

	newCluster := cluster{
		KubeConfigContents: kubeConfigBytes,
		Expires:            time.Now().Add(time.Hour * time.Duration(expires)),
	}

	heap.Push(&clusterPool, newCluster)

	fmt.Fprintf(w, "New Cluster added successfully!")
}

func main() {

	instantiateClusterPool()

	http.HandleFunc("/request-build", buildRequestHandler)
	http.HandleFunc("/add-cluster", clusterAdditionHandler)

	fmt.Printf("Starting server at port %d\n", serverPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil); err != nil {
		log.Fatal(err)
	}

}
