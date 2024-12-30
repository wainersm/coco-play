package main

import (
	"fmt"
	"os"
	"path"

	"github.com/wainersm/coco-play/cmd"
	"github.com/wainersm/coco-play/pkg/cluster"
)

func main() {
	var err error

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	cluster.KubeConfig = path.Join(home, ".kube", "config")
	cluster.ClusterName = "coco-play"

	cmd.Execute()
}
