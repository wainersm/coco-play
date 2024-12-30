/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wainersm/coco-play/pkg/versions"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the components versions",
	Long: `Print the versions of the components to be installed.

Example of output:

$ coco-play version
coco-play version: v0.10.0-1b05731
CoCo version: v0.10.0
KBS version: v0.10.1
`,
	Run: func(cmd *cobra.Command, args []string) {
		toolVersion := versions.CocoVersion + "-" + versions.GitCommit
		msg := `coco-play version: %s
CoCo version: %s
KBS version: %s
`
		fmt.Printf(msg, toolVersion, versions.CocoVersion, versions.KbsVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
