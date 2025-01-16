/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package e2e

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/wainersm/coco-play/cmd"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/e2e-framework/klient"
)

// GetPodLog returns the logs of a pod similarly to `kubectl logs pod NAME`
//
// Note: this is a copy of cloud-api-adaptor/src/cloud-api-adaptor/test/e2e/assessment_helpers.go#GetPodLog
func GetPodLog(ctx context.Context, client klient.Client, pod *v1.Pod) (string, error) {
	clientset, err := kubernetes.NewForConfig(client.RESTConfig())
	if err != nil {
		return "", err
	}

	req := clientset.CoreV1().Pods(pod.ObjectMeta.Namespace).GetLogs(pod.ObjectMeta.Name, &v1.PodLogOptions{})
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer podLogs.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

// RunCmd run a coco-play command and return the output captured
func RunCmd(name string, args ...string) ([]byte, error) {
	// Save the original stdout to restore later
	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
	}()

	r, w, err := os.Pipe()
	if err != nil {
		return []byte{}, err
	}
	os.Stdout = w

	stdoutBuff := make(chan []byte)
	go func() {
		var b bytes.Buffer
		if _, err := io.Copy(&b, r); err != nil {
			log.Println(err)
		}
		stdoutBuff <- b.Bytes()
	}()

	cmd.RootCmd.SetArgs(append([]string{name}, args...))
	err = cmd.RootCmd.Execute()
	os.Stdout.Close()

	return <-stdoutBuff, err
}
