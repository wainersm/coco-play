package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const KbsNamespace string = "coco-tenant"

func installKbs(version string) error {
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
	out, err := kubectl("rollout", "status", "-w", "deployment/kbs", "-n", KbsNamespace)
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

// getKbsAddress returns the host:port address of KBS
func getKbsAddress() (string, error) {
	host, err := kubectl("get", "nodes", "-o", "jsonpath='{.items[0].status.addresses[?(@.type==\"InternalIP\")].address}'", "-n", KbsNamespace)
	if err != nil {
		return "", err
	}
	host = strings.Trim(host, "'")

	port, err := kubectl("get", "svc", "kbs", "-n", KbsNamespace, "-o", "jsonpath='{.spec.ports[0].nodePort}'")
	if err != nil {
		return "", err
	}
	port = strings.Trim(port, "'")

	return fmt.Sprintf("%s:%s", host, port), nil
}

func setKbsResource(path, resourcefile string) error {
	//kubectl exec deploy/kbs -- mkdir -p "/opt/confidential-containers/kbs/repository/$(dirname "$KEY_PATH")"
	//cat "$KEY_FILE" | kubectl exec -i deploy/kbs -- tee "/opt/confidential-containers/kbs/repository/${KEY_PATH}" > /dev/null

	return nil
}
