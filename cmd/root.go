/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "coco-play",
	Short: "Confidential Containers playground tool",
	Long: `The coco-play tool provides an environment where you can play with
Confidential Containers (CoCo) in your preferred workstation without
necessarily having confidential hardware (TEE) available, thanks to some
components mocking the remote attestation procedures.

Also it is a handy tool for developing and validating your workloads before
deploying on a production environments.

The playground consists of a Kubernetes-in-Docker (Kind) cluster plus the minimal
CoCo installation with remote attestation.

To use it you must have installed in your local system:

- Docker
- Qemu/KVM
- Kubectl
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.coco-play.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Disable the completion command
	RootCmd.CompletionOptions.DisableDefaultCmd = true
}
