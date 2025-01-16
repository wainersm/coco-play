/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package e2e

import (
	"context"
	"os"
	"testing"

	"github.com/wainersm/coco-play/cmd"
	"github.com/wainersm/coco-play/pkg/cluster"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

var (
	testEnv env.Environment
)

func TestMain(m *testing.M) {
	testEnv = env.New()

	// Run before all tests, indirectly testing the play-create command.
	testEnv.Setup(func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		cluster.ClusterName = "coco-play"
		// Run play-create command
		cmd.RootCmd.SetArgs([]string{"play-create"})
		err := cmd.RootCmd.Execute()

		return ctx, err
	})

	// Run after all tests, indirectly testing the play-delete command.
	testEnv.Finish(func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		cluster.ClusterName = "coco-play"
		// Run play-delete command
		cmd.RootCmd.SetArgs([]string{"play-delete"})
		err := cmd.RootCmd.Execute()

		return ctx, err
	})

	os.Exit(testEnv.Run(m))
}
