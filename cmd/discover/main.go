/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"os"

	"github.com/hyperledger/mchain/bccsp/factory"
	"github.com/hyperledger/mchain/cmd/common"
	"github.com/hyperledger/mchain/discovery/cmd"
)

func main() {
	factory.InitFactories(nil)
	cli := common.NewCLI("discover", "Command line client for mchain discovery service")
	discovery.AddCommands(cli)
	cli.Run(os.Args[1:])
}
