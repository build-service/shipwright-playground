package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	shipwright "github.com/shipwright-io/build/pkg/apis/build/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"

	goyaml "github.com/go-yaml/yaml"
)

var (
	strategyKind shipwright.BuildStrategyKind = shipwright.ClusterBuildStrategyKind
)

func TypeMeta(kind, apiVersion string) v1.TypeMeta {
	return v1.TypeMeta{
		Kind:       kind,
		APIVersion: apiVersion,
	}
}

func ObjectMeta(n types.NamespacedName, opts ...objectMetaFunc) v1.ObjectMeta {
	om := v1.ObjectMeta{
		Namespace: n.Namespace,
		Name:      n.Name,
	}
	for _, o := range opts {
		o(&om)
	}
	return om

}

// handleDockerCfgJSONContent serializes a ~/.docker/config.json file
func handleDockerCfgJSONContent(username, password, email, server string) ([]byte, error) {
	dockerConfigAuth := DockerConfigEntry{
		Username: username,
		Password: password,
		Email:    email,
		Auth:     encodeDockerConfigFieldAuth(username, password),
	}
	dockerConfigJSON := DockerConfigJSON{
		Auths: map[string]DockerConfigEntry{server: dockerConfigAuth},
	}

	return json.Marshal(dockerConfigJSON)
}

// encodeDockerConfigFieldAuth returns base64 encoding of the username and password string
func encodeDockerConfigFieldAuth(username, password string) string {
	fieldValue := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(fieldValue))
}

func SplitYAML(resources []byte) ([][]byte, error) {
	dec := goyaml.NewDecoder(bytes.NewReader(resources))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := goyaml.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
}

func testClusterConnection(k8sClient *kubernetes.Clientset) (bool, error) {
	_, err := k8sClient.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return false, err
	}

	return true, nil

}

func generateUniqueImageName(quayServer, imageRegistryOrg, gitHubOrg, repoName string) (string, error) {
	uniqueID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("Could not generate UUID: %v", err)
	}

	imageName := fmt.Sprintf("%v/%v/%s-%s:%v", quayServer, imageRegistryOrg, gitHubOrg, repoName, uniqueID)
	return imageName, nil

}

func parseGitRepoURL(repoURL string) (string, string) {
	components := strings.Split(repoURL, "/")
	return components[len(components)-2], components[len(components)-1]
}
