/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package e2e

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func NewPod(name, namespace string) (*v1.Pod, error) {
	runtimeClass := "kata-qemu-coco-dev"

	kbsInfo, err := RunCmd("kbs-info")
	if err != nil {
		return &v1.Pod{}, err
	}
	re := regexp.MustCompile(`Service address: ([0-9]+(\.)?)+:[0-9]+`)
	line := re.FindString(string(kbsInfo))
	if line == "" {
		return &v1.Pod{}, fmt.Errorf("Failed to find KBS address")
	}
	kbsAddr := strings.TrimSpace(strings.Replace(line, "Service address: ", "", 1))

	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Annotations: map[string]string{
				"io.containerd.cri.runtime-handler":                 "kata-qemu-coco-dev",
				"io.katacontainers.config.hypervisor.kernel_params": "agent.aa_kbc_params=cc_kbc::http://" + kbsAddr,
			},
		},
		Spec: v1.PodSpec{
			RuntimeClassName: &runtimeClass,
			Containers: []v1.Container{
				{
					Name:            "busybox",
					Image:           "quay.io/prometheus/busybox:latest",
					ImagePullPolicy: v1.PullAlways,
					Command:         []string{"sleep", "infinity"},
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}, nil
}

func TestCreateSimplePod(t *testing.T) {
	f := features.New("Create simple pod").
		WithLabel("type", "e2e").
		Assess("Should get a secret", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client, err := cfg.NewClient()
			if err != nil {
				t.Errorf("Failed to get client: %v", err)
				return ctx
			}

			pod, _ := NewPod("coco-test1", "default")
			pod.Spec.Containers[0].Command = []string{
				"sh", "-c",
				"wget -O- http://127.0.0.1:8006/cdh/resource/reponame/workload_key/key.bin; sleep infinity"}

			if err = client.Resources().Create(ctx, pod); err != nil {
				t.Errorf("Failed to create pod: %v", err)
				return ctx
			}

			ctx = context.WithValue(ctx, "test-pod", pod)

			if err = wait.For(conditions.New(client.Resources()).PodRunning(pod), wait.WithTimeout(time.Second*30)); err != nil {
				t.Errorf("Wait pod failed: %v", err)
				return ctx
			}

			// Wait the container command finish
			time.Sleep(time.Second * 10)

			logs, err := GetPodLog(ctx, client, pod)
			if err != nil {
				t.Errorf("Failed to get pod logs: %v", err)
				return ctx
			}
			fmt.Println(logs)

			if !strings.Contains(logs, "somesecret") {
				t.Errorf("Secret string not found on the pod logs")
			}

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			pod := ctx.Value("test-pod").(*v1.Pod)
			if pod == nil {
				return ctx
			}

			client, err := cfg.NewClient()
			if err != nil {
				t.Errorf("Failed to get client: %v", err)
				return ctx
			}

			if err = client.Resources().Delete(ctx, pod); err != nil {
				t.Errorf("Failed to delete pod: %v", err)
				return ctx
			}

			return ctx
		}).
		Feature()

	testEnv.Test(t, f)
}
