package main

import (
	"fmt"

	shipwright "github.com/shipwright-io/build/pkg/apis/build/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func createImageRegistrySecret(username, password, email, server string) (*corev1.Secret, error) {

	dockerConfigJSONContent, err := handleDockerCfgJSONContent(username, password, email, server)
	if err != nil {
		return nil, fmt.Errorf("error creating docker config json content : %v", err)
	}

	imageRegistrySecret := &corev1.Secret{
		TypeMeta:   TypeMeta("Secret", "v1"),
		ObjectMeta: ObjectMeta(types.NamespacedName{Namespace: "default", Name: secretName}),
		Type:       corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: dockerConfigJSONContent,
		},
	}

	return imageRegistrySecret, nil
}

func createBuild(imageRegistryServer, gitHubOrg, repoURL, imageRegistryOrg, repoName, secretName, contextDir string) (*shipwright.Build, error) {

	imageName, err := generateUniqueImageName(imageRegistryServer, imageRegistryOrg, gitHubOrg, repoName)
	if err != nil {
		return nil, fmt.Errorf("Error generating UUID for image tag: %v", err)
	}

	build := &shipwright.Build{
		TypeMeta:   TypeMeta("Build", "shipwright.io/v1alpha1"),
		ObjectMeta: ObjectMeta(types.NamespacedName{Namespace: "", Name: fmt.Sprintf("%v-build", repoName)}),
		Spec: shipwright.BuildSpec{
			Source: shipwright.Source{
				URL:        repoURL,
				ContextDir: &contextDir,
			},
			Strategy: &shipwright.Strategy{
				Name: "buildpacks-v3",
				Kind: &strategyKind,
			},
			Output: shipwright.Image{
				Image: imageName,
				Credentials: &corev1.LocalObjectReference{
					Name: secretName,
				},
			},
		},
	}

	return build, nil
}

func createBuildRun(repoName string) *shipwright.BuildRun {

	buildRun := &shipwright.BuildRun{
		TypeMeta: TypeMeta("BuildRun", "shipwright.io/v1alpha1"),
		// ObjectMeta: v1.ObjectMeta{GenerateName: fmt.Sprintf("%v-buildrun-", repoName)},
		ObjectMeta: ObjectMeta(types.NamespacedName{Namespace: "default", Name: fmt.Sprintf("%v-buildrun", repoName)}),
		Spec: shipwright.BuildRunSpec{
			BuildRef: &shipwright.BuildRef{
				Name: fmt.Sprintf("%v-build", repoName),
			},
		},
	}

	return buildRun
}
