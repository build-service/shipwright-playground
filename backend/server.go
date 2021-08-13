package main

import (
	"container/heap"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	shipwright "github.com/shipwright-io/build/pkg/apis/build/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	dockerServer        = "https://index.docker.io/v1/"
	quayServer          = "quay.io"
	secretName          = "image-registry-secret"
	contextDir          string
	imageRegistryServer string
	serverPort          = 8085
	clusterPool         PriorityQueue
)

func buildRequestHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	imageRegistry := r.FormValue("image-registry")
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	repoURL := r.FormValue("repo-url")
	contextDir = r.FormValue("context-dir")
	repoName := repoURL[strings.LastIndex(repoURL, "/")+1:]

	if imageRegistry == "docker" {
		imageRegistryServer = dockerServer
	} else {
		imageRegistryServer = quayServer
	}

	k8sClient, _ := getk8sClient(currentClusterRestConfig)
	_, err := k8sClient.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		_, err := getNewClusterRestConfig()
		if err != nil {
			fmt.Fprintf(w, "Unable to establish connection with an active cluster: %v", err)
			log.Fatalf("Unable to establish connection with an active cluster: %v", err)
		} else {
			k8sClient, _ = getk8sClient(currentClusterRestConfig)
		}
	}
	buildClient, _ := getShipwrightClient(currentClusterRestConfig)

	if _, err := k8sClient.CoreV1().Secrets("default").Get(context.TODO(), secretName, v1.GetOptions{}); err == nil {
		err := k8sClient.CoreV1().Secrets("default").Delete(context.TODO(), secretName, v1.DeleteOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
	dockerSecret, err := createDockerSecret(username, password, email, imageRegistryServer)
	_, err = k8sClient.CoreV1().Secrets("default").Create(context.TODO(), dockerSecret, v1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}

	if _, err := buildClient.ShipwrightV1alpha1().Builds("default").Get(context.TODO(), fmt.Sprintf("%v-build", repoName), v1.GetOptions{}); err == nil {
		err := buildClient.ShipwrightV1alpha1().Builds("default").Delete(context.TODO(), fmt.Sprintf("%v-build", repoName), v1.DeleteOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
	buildObj := createBuild(imageRegistry, repoURL, username, repoName, secretName, contextDir)
	_, err = buildClient.ShipwrightV1alpha1().Builds("default").Create(context.TODO(), buildObj, v1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}

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

	kubeConfigContents := r.FormValue("kubeconfig-contents")
	expirationHours := r.FormValue("expires")
	expires, _ := strconv.Atoi(expirationHours)

	newCluster := cluster{
		KubeConfigContents: kubeConfigContents,
		Expires:            time.Now().Add(time.Hour * time.Duration(expires)),
	}

	heap.Push(&clusterPool, newCluster)

	fmt.Fprintf(w, "New Cluster added successfully!")
}

func main() {

	instantiateClusterPool()

	newClusterRestConfig, err := getNewClusterRestConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = applyManifestsToCluster(manifestURLs, newClusterRestConfig)
	if err != nil {
		log.Fatalf(" Aplication of required manifests failed: %v", err)
	}

	http.HandleFunc("/request-build", buildRequestHandler)
	http.HandleFunc("/add-cluster", clusterAdditionHandler)

	fmt.Printf("Starting server at port %d\n", serverPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil); err != nil {
		log.Fatal(err)
	}

}
