/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package coco

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/wainersm/coco-play/pkg/cluster"
)

func Install(version string) error {
	namespace := "confidential-containers-system"

	fmt.Println("Labeling worker...")

	out, err := cluster.Kubectl("label", "node", cluster.ClusterName+"-control-plane", "node.kubernetes.io/worker=")
	if err != nil {
		return errors.Errorf("Failed to label node: %s", out)
	}
	fmt.Println(out)

	fmt.Println("Creating CoCo controller...")
	if out, err = cluster.Kubectl("apply", "-k", "github.com/confidential-containers/operator/config/release?ref="+version); err != nil {
		return errors.Errorf("Failed to create controller: %v", err)
	}
	fmt.Println(out)

	if out, err = cluster.Kubectl("rollout", "status", "-w", "deployment/cc-operator-controller-manager", "-n", namespace); err != nil {
		return errors.Errorf("Controller is not ready: %v", err)
	}
	fmt.Println(out)

	fmt.Println("Creating CoCo ccruntime CRD...")

	if out, err = cluster.Kubectl("apply", "-k", "github.com/confidential-containers/operator/config/samples/ccruntime/default?ref="+version); err != nil {
		errors.Errorf("Failed to create ccruntime: %v\n", err)
	}
	fmt.Println(out)

	pattern, _ := regexp.Compile("^cc-operator-daemon-install.*Running.*")
	// TODO: make timeout configurable
	for i := range "1..10" {
		time.Sleep(time.Second * 30)
		if out, err = cluster.Kubectl("get", "pod", "-n", namespace); err != nil {
			continue
		}
		found := false
		for _, line := range strings.Split(out, "\n") {
			if pattern.Match([]byte(line)) {
				found = true
				fmt.Println(line)
				break
			}
		}
		if found {
			break
		} else if i == 9 {
			return errors.Errorf("ccruntime creation is not done after timeout")
		}
	}

	return nil
}
