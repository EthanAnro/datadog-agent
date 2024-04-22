// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package examples

import (
	"context"
	"strings"
	"testing"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/e2e"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/environments"
	awskubernetes "github.com/DataDog/datadog-agent/test/new-e2e/pkg/environments/aws/kubernetes"

	"github.com/DataDog/test-infra-definitions/common/config"
	"github.com/DataDog/test-infra-definitions/components/datadog/apps/dogstatsd"
	"github.com/DataDog/test-infra-definitions/components/datadog/apps/nginx"
	compkube "github.com/DataDog/test-infra-definitions/components/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type myKindSuite struct {
	e2e.BaseSuite[environments.Kubernetes]
}

func TestMyKindSuite(t *testing.T) {
	e2e.Run(t, &myKindSuite{}, e2e.WithProvisioner(
		awskubernetes.Provisioner(
			awskubernetes.WithoutFakeIntake(),
			awskubernetes.WithWorkloadApp(func(e config.CommonEnvironment, kubeProvider *kubernetes.Provider) (*compkube.Workload, error) {
				return nginx.K8sAppDefinition(e, kubeProvider, "nginx", "", nil)
			}),
			awskubernetes.WithWorkloadApp(func(e config.CommonEnvironment, kubeProvider *kubernetes.Provider) (*compkube.Workload, error) {
				return dogstatsd.K8sAppDefinition(e, kubeProvider, "dogstatsd", 8125, "/var/run/datadog/dsd.socket")
			}),
		)))
}

func (v *myKindSuite) TestClusterAgentInstalled() {
	res, _ := v.Env().KubernetesCluster.Client().CoreV1().Pods("datadog").List(context.TODO(), v1.ListOptions{})
	containsClusterAgent := false
	for _, pod := range res.Items {
		if strings.Contains(pod.Name, "cluster-agent") {
			containsClusterAgent = true
			break
		}
	}
	assert.True(v.T(), containsClusterAgent, "Cluster Agent not found")
	assert.Equal(v.T(), v.Env().Agent.InstallNameLinux, "dda")
}