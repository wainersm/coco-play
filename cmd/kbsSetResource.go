/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wainersm/coco-play/pkg/kbs"
)

// kbsSetResourceCmd represents the kbs-set-resource command
var kbsSetResourceCmd = &cobra.Command{
	Use:     "kbs-set-resource path resource_file",
	Short:   "Set KBS confidential resource",
	Example: "$ coco-play default/signatures/key /path/to/file",
	Long: `Set a confidential resource in the KBS.

Examples of resources:
 - image decryption key
 - public key for image signature check
 - secret to be fetch from CDH by the workload
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return err
		}
		if len(strings.Split(args[0], "/")) != 3 {
			return fmt.Errorf("resource path (%s) should have 3 segments like respository/type/resource", args[0])
		}
		if _, err := os.Stat(args[1]); err != nil {
			return fmt.Errorf("resource file (%s) does not exist", args[1])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := kbs.SetResource(args[0], args[1]); err != nil {
			fmt.Printf("Failed to set resource %s from %s file: %v\n", args[0], args[1], err)
		}
	},
}

func init() {
	RootCmd.AddCommand(kbsSetResourceCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kbsSetResourceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kbsSetResourceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
