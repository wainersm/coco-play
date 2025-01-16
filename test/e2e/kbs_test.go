/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package e2e

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

// Test commands related with KBS
func TestKbs(t *testing.T) {
	infoCmd := features.New("Command kbs-info").
		WithLabel("cmd", "kbs-info").
		Assess("Should work", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			out, err := RunCmd("kbs-info")
			// For debugging sake
			fmt.Print(string(out))
			if err != nil {
				t.Fail()
			} else {
				expectRegex := regexp.MustCompile(`Status: Running\nService address: ([0-9]+(\.)?)+:[0-9]+`)
				if !expectRegex.Match(out) {
					t.Fail()
				}
			}

			return ctx
		}).
		Feature()

	setResourceCmd := features.New("Command kbs-set-resource").
		WithLabel("cmd", "kbs-set-resource").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			f, err := os.CreateTemp("", "coco-play-test-*.txt")
			if err != nil {
				ctx = context.WithValue(ctx, "setup-error", fmt.Errorf("Failed to create temporary file: %v", err))
			} else {
				ctx = context.WithValue(ctx, "secret-file", f)
				secret := "anothersecret"
				ctx = context.WithValue(ctx, "secret", secret)
				if _, err = f.WriteString(secret); err != nil {
					ctx = context.WithValue(ctx, "setup-error", fmt.Errorf("Failed to write secret to %s: %v", f.Name(), err))
				}
			}

			return ctx
		}).
		Assess("Should work", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			// Check whether Setup() failed or not
			err := ctx.Value("setup-error")
			if err != nil {
				t.Errorf("setup failed: %v", err)
				return ctx
			}

			f := ctx.Value("secret-file").(*os.File)
			out, err := RunCmd("kbs-set-resource", "default/tests/key", f.Name())
			if err != nil {
				t.Errorf("Failed to run kbs-set-resource command: %v", err)
				return ctx
			}

			// The resource value is echo'ed on command's output
			secret := ctx.Value("secret").(string)
			if !strings.Contains(string(out), secret) {
				t.Errorf("Failed to insert new secret")
			}
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			f := ctx.Value("secret-file").(*os.File)
			if f != nil {
				os.Remove(f.Name())
			}
			return ctx
		}).
		Feature()

	testEnv.Test(t, infoCmd, setResourceCmd)
}
