// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build kubeapiserver

package autoinstrumentation

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	workloadmeta "github.com/DataDog/datadog-agent/comp/core/workloadmeta/def"
	"github.com/DataDog/datadog-agent/pkg/util/kubernetes"
)

// getOwnerNameAndKind returns the name and kind of the first owner of the pod if it exists
// if the first owner is a replicaset, it returns the name
func getOwnerNameAndKind(pod *corev1.Pod) (string, string, bool) {
	owners := pod.GetOwnerReferences()

	if len(owners) == 0 {
		return "", "", false
	}

	owner := owners[0]
	ownerName, ownerKind := owner.Name, owner.Kind

	if ownerKind == "ReplicaSet" {
		deploymentName := kubernetes.ParseDeploymentForReplicaSet(ownerName)
		if deploymentName != "" {
			ownerKind = "Deployment"
			ownerName = deploymentName
		}
	}

	return ownerName, ownerKind, true
}

func getLibListFromDeploymentAnnotations(store workloadmeta.Component, deploymentName, ns, registry string) []libInfo {
	// populate libInfoList using the languages found in workloadmeta
	id := fmt.Sprintf("%s/%s", ns, deploymentName)
	deployment, err := store.GetKubernetesDeployment(id)
	if err != nil {
		return nil
	}

	var libList []libInfo
	for container, languages := range deployment.InjectableLanguages {
		for lang := range languages {
			// There's a mismatch between language detection and auto-instrumentation.
			// The Node language is a js lib.
			if lang == "node" {
				lang = "js"
			}

			l := language(lang)
			libList = append(libList, l.defaultLibInfo(registry, container.Name))
		}
	}

	return libList
}
