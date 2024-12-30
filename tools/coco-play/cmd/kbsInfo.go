/*
Copyright © 2024 Confidential Containers Contributors
*/
package cmd

import (
	"fmt"

	"github.com/confidential-containers/confidential-containers/tools/coco-play/pkg/kbs"
	"github.com/spf13/cobra"
)

// kbsInfoCmd represents the kbs-info command
var kbsInfoCmd = &cobra.Command{
	Use:          "kbs-info",
	Short:        "Print information about the KBS in use",
	SilenceUsage: true,
	Long: `This command print the following information about the KBS:
- Status
- Service address`,
	Example: `
$ coco-play kbs-info
Status: Running
Service address: 172.18.0.2:30945`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if status, err := kbs.GetStatus(); err != nil {
			return err
		} else {
			fmt.Printf("Status: %s\n", status)
		}
		if addr, err := kbs.GetAddress(); err != nil {
			return err
		} else {
			fmt.Printf("Service address: %s\n", addr)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(kbsInfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kbsInfoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kbsInfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
