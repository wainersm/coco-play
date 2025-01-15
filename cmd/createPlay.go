/*
Copyright Confidential Containers Contributors
SPDX-License-Identifier: Apache-2.0
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wainersm/coco-play/pkg/cluster"
	"github.com/wainersm/coco-play/pkg/coco"
	"github.com/wainersm/coco-play/pkg/kbs"
	"github.com/wainersm/coco-play/pkg/versions"
)

// createPlayCmd represents the play-create command
var createPlayCmd = &cobra.Command{
	Use:          "play-create",
	Short:        "Create a new playground",
	SilenceUsage: true,
	Long: `This command will create a new playground. The following are performed:

- Create a new Kind (https://kind.sigs.k8s.io) cluster
- Install the CoCo operator
- Install the KBS`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if err = cluster.CreateCluster(); err != nil {
			return err
		}

		if err = coco.Install(versions.CocoVersion); err != nil {
			return err
		}

		if err = kbs.InstallKbs(versions.KbsVersion); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(createPlayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createPlayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createPlayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
