package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"regexp"
	"time"

	"github.com/pkg/errors"
	"sigs.k8s.io/kind/pkg/cluster"
)

const DefaultClusterWaitForReadyTimeout = time.Minute * 5

// Path to kubeconfig file
var kubeConfig string

// Cluster name as user sees it
var clusterName string

func getInternalClusterName() string {
	return "kind-" + clusterName
}

func kubectl(arg ...string) (string, error) {
	var out strings.Builder
	newArgs := []string{"--kubeconfig", kubeConfig, "--cluster",
		getInternalClusterName()}
	newArgs = append(newArgs, arg...)

	cmd := exec.Command("kubectl", newArgs...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	return out.String(), err
}

func createCluster() error {
	fmt.Printf("Creating Kind cluster (%s)...\n", clusterName)
	provider := cluster.NewProvider(cluster.ProviderWithDocker())

	err := provider.Create(clusterName, cluster.CreateWithWaitForReady(DefaultClusterWaitForReadyTimeout))
	if err != nil {
		return errors.Errorf("Failed to create cluster: %v\n", err)
	}

	if err = provider.ExportKubeConfig(clusterName, kubeConfig, false); err != nil {
		return errors.Errorf("Failed to export kubeconfig to %s: %v", kubeConfig, err)
	}

	return nil
}

func installCoco(version string) error {
	namespace := "confidential-containers-system"

	fmt.Println("Labeling worker...")

	out, err := kubectl("label", "node", clusterName+"-control-plane", "node.kubernetes.io/worker=")
	if err != nil {
		return errors.Errorf("Failed to label node: %s", out)
	}
	fmt.Println(out)

	fmt.Println("Creating CoCo controller...")
	if out, err = kubectl("apply", "-k", "github.com/confidential-containers/operator/config/release?ref="+version); err != nil {
		return errors.Errorf("Failed to create controller: %v", err)
	}
	fmt.Println(out)

	if out, err = kubectl("rollout", "status", "-w", "deployment/cc-operator-controller-manager", "-n", namespace); err != nil {
		return errors.Errorf("Controller is not ready: %v", err)
	}
	fmt.Println(out)

	fmt.Println("Creating CoCo ccruntime CRD...")

	if out, err = kubectl("apply", "-k", "github.com/confidential-containers/operator/config/samples/ccruntime/default?ref="+version); err != nil {
		errors.Errorf("Failed to create ccruntime: %v\n", err)
	}
	fmt.Println(out)

	pattern, _ := regexp.Compile("^cc-operator-daemon-install.*Running.*")
	// TODO: make timeout configurable
	for i := range "1..10" {
		time.Sleep(time.Second * 30)
		if out, err = kubectl("get", "pod", "-n", namespace); err != nil {
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

func main() {
	var err error

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	kubeConfig = path.Join(home, ".kube", "config")
	clusterName = "coco-play"

	if err = createCluster(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	cocoVersion := "v0.10.0"
	if err = installCoco(cocoVersion); err != nil {
		fmt.Printf("Failed to install CoCo: %v", err)
		os.Exit(1)
	}

	if err = installKbs("v0.10.1"); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

}
