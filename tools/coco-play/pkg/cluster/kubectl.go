/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package cluster

import (
	"os/exec"
	"strings"
)

func getInternalClusterName() string {
	return "kind-" + ClusterName
}

func Kubectl(arg ...string) (string, error) {
	var out strings.Builder
	newArgs := []string{"--kubeconfig", KubeConfig, "--cluster",
		getInternalClusterName()}
	newArgs = append(newArgs, arg...)

	cmd := exec.Command("kubectl", newArgs...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	return out.String(), err
}
