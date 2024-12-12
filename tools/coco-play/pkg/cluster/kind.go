/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package cluster

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"sigs.k8s.io/kind/pkg/cluster"
)

const DefaultClusterWaitForReadyTimeout = time.Minute * 5

var (
	ClusterName string
	KubeConfig  string
)

func CreateCluster() error {
	fmt.Printf("Creating Kind cluster (%s)...\n", ClusterName)
	provider := cluster.NewProvider(cluster.ProviderWithDocker())

	err := provider.Create(ClusterName, cluster.CreateWithWaitForReady(DefaultClusterWaitForReadyTimeout))
	if err != nil {
		return errors.Errorf("Failed to create cluster: %v\n", err)
	}

	if err = provider.ExportKubeConfig(ClusterName, KubeConfig, false); err != nil {
		return errors.Errorf("Failed to export kubeconfig to %s: %v", KubeConfig, err)
	}

	return nil
}

// DeleteCluster deletes the Kind cluster
func DeleteCluster() error {
	fmt.Printf("Deleting Kind cluster (%s)...\n", ClusterName)
	provider := cluster.NewProvider(cluster.ProviderWithDocker())

	if err := provider.Delete(ClusterName, KubeConfig); err != nil {
		return errors.Errorf("Failed to delete cluster: %v", err)
	}

	return nil
}
