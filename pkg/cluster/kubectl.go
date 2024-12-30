/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package cluster

import (
	"os"
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

func Kubectl_exec(resource string, namespace string, stdinFromFile string, arg ...string) (string, error) {
	var out strings.Builder
	newArgs := []string{"--kubeconfig", KubeConfig, "--cluster",
		getInternalClusterName(), "exec", resource}

	if namespace != "" {
		newArgs = append(newArgs, "-n", namespace)
	}

	if stdinFromFile != "" {
		newArgs = append(newArgs, "-i")
	}
	newArgs = append(newArgs, "--")
	newArgs = append(newArgs, arg...)

	cmd := exec.Command("kubectl", newArgs...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	if stdinFromFile != "" {
		file, err := os.Open(stdinFromFile)
		if err != nil {
			return "", err
		}
		defer file.Close()
		cmd.Stdin = file
	}
	err := cmd.Run()

	return out.String(), err
}
