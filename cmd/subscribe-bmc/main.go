// Copyright 2023 NJWS Inc.

package main

import (
	"fmt"
	"os"

	"github.com/foliagecp/inventory-bmc-app/pkg/subscribe"
)

func main() {
	if err := subscribe.CLI.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
