// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package main

import (
	"fmt"
	"os"

	"github.com/cilium/cilium-cli/cli"

	"github.com/cilium/cilium/pkg/cilium-cli/api"
	"github.com/cilium/cilium/pkg/cilium-cli/cmd"
)

func main() {
	command := cli.NewDefaultCiliumCommand()
	command.AddCommand(cmd.NewCmdConnectivity(&api.NopHooks{}))
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
