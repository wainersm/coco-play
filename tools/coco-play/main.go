package main

import (
	"fmt"
	"os"
	"path"

	"github.com/confidential-containers/confidential-containers/tools/coco-play/pkg/cluster"
	"github.com/confidential-containers/confidential-containers/tools/coco-play/pkg/coco"
	"github.com/confidential-containers/confidential-containers/tools/coco-play/pkg/kbs"
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

	if err = cluster.CreateCluster(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	cocoVersion := "v0.10.0"
	if err = coco.Install(cocoVersion); err != nil {
		fmt.Printf("Failed to install CoCo: %v", err)
		os.Exit(1)
	}

	if err = kbs.InstallKbs("v0.10.1"); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
