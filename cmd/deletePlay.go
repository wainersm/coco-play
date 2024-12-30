/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wainersm/coco-play/pkg/cluster"
)

// deletePlayCmd represents the play-delete command
var deletePlayCmd = &cobra.Command{
	Use:   "play-delete",
	Short: "Delete a playground",
	Long:  `This command will delete the playground by deleting the cluster entirely`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cluster.DeleteCluster(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deletePlayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deletePlayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deletePlayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
