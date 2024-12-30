/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package kbs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
	"github.com/wainersm/coco-play/pkg/cluster"
)

const KbsNamespace string = "coco-tenant"

func InstallKbs(version string) error {
	var err error

	fmt.Println("Install KBS...")

	trusteeDir, err := os.MkdirTemp("", "trustee-")
	if err != nil {
		return errors.Errorf("%v", err)
	}

	if _, err = git.PlainClone(trusteeDir, false, &git.CloneOptions{
		URL:           "https://github.com/confidential-containers/trustee",
		SingleBranch:  true,
		ReferenceName: plumbing.NewTagReferenceName(version),
		Progress:      os.Stdout,
	}); err != nil {
		return errors.Errorf("Failed to clone repository: %v", err)
	}

	defer func() {
		os.RemoveAll(trusteeDir)
	}()

	if err = os.WriteFile(trusteeDir+"/kbs/config/kubernetes/overlays/x86_64/key.bin", []byte("somesecret\n"), 0666); err != nil {
		return errors.Errorf("Faile to write key secret: %v", err)
	}

	/*
	 * Deploy KBS with allow_all policy already set.
	 */
	fRead, err := os.ReadFile(trusteeDir + "/kbs/sample_policies/allow_all.rego")
	if err != nil {
		return errors.Errorf("Failed to open sample policy file: %v", err)
	}
	err = os.WriteFile(trusteeDir+"/kbs/config/kubernetes/base/policy.rego", fRead, 0644)
	if err != nil {
		return errors.Errorf("Failed to write new policy file: %v", err)
	}

	cmd := exec.Command("./deploy-kbs.sh")
	cmd.Dir = trusteeDir + "/kbs/config/kubernetes"
	cmd.Env = append(cmd.Environ(), "DEPLOYMENT_DIR=nodeport")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Errorf("Failed to deploy KBS: %s", stdoutStderr)
	}

	fmt.Printf(string(stdoutStderr))

	/*
	 * Wait kbs is ready
	 */
	out, err := cluster.Kubectl("rollout", "status", "-w", "deployment/kbs", "-n", KbsNamespace)
	if err != nil {
		return errors.Errorf("Deployment is not ready after timeout: %v", err)
	}
	fmt.Printf(out)

	/* ???
	# copy to ~/.coco
	#"trustee/kbs/config/kubernetes/base/kbs.key"
	*/

	return nil
}

// GetAddress returns the host:port address of KBS
func GetAddress() (string, error) {
	host, err := cluster.Kubectl("get", "nodes", "-o", "jsonpath='{.items[0].status.addresses[?(@.type==\"InternalIP\")].address}'", "-n", KbsNamespace)
	if err != nil {
		return "", errors.Errorf("Failed to get host address: %v", err)
	}
	host = strings.Trim(host, "'")

	port, err := cluster.Kubectl("get", "svc", "kbs", "-n", KbsNamespace, "-o", "jsonpath='{.spec.ports[0].nodePort}'")
	if err != nil {
		return "", errors.Errorf("Failed to get port address: %v", err)
	}
	port = strings.Trim(port, "'")

	return fmt.Sprintf("%s:%s", host, port), nil
}

// GetStatus returns the status of KBS, i.e, raw status of the kbs pod
func GetStatus() (string, error) {
	status, err := cluster.Kubectl("get", "pods", "-l", "app=kbs", "-n", KbsNamespace, "-o", "jsonpath='{.items[0].status.phase}'")
	status = strings.Trim(status, "'")
	if err != nil {
		return "", errors.Errorf("Failed to get KBS status: %v", err)
	}

	return status, nil
}

func SetResource(path, resourcefile string) error {
	basedir := path[:strings.LastIndex(path, "/")]
	out, err := cluster.Kubectl_exec("deploy/kbs", KbsNamespace, "", "mkdir", "-p",
		"/opt/confidential-containers/kbs/repository/"+basedir)
	if err != nil {
		fmt.Print(out)
		return errors.Errorf("Failed to set resource (%s) in KBS: %v", path, err)
	}
	fmt.Print(out)

	if out, err = cluster.Kubectl_exec("deploy/kbs", KbsNamespace, resourcefile, "tee",
		"/opt/confidential-containers/kbs/repository/"+path); err != nil {
		return errors.Errorf("Failed to set resource (%s) from %s in KBS: %v", path, resourcefile, err)
	}
	fmt.Print(out)

	return nil
}
