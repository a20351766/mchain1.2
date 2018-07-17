// +build pluginsenabled,cgo
// +build darwin,go1.10 linux,go1.10 linux,go1.9,!ppc64le

/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc

import (
	"github.com/hyperledger/mchain/core/aclmgmt"
	"github.com/hyperledger/mchain/core/common/ccprovider"
)

// CreateSysCCs creates all of the system chaincodes which are compiled into mchain
// as well as those which are loaded by plugin
func CreateSysCCs(ccp ccprovider.ChaincodeProvider, p *Provider, aclProvider aclmgmt.ACLProvider) []*SystemChaincode {
	return append(builtInSystemChaincodes(ccp, p, aclProvider), loadSysCCs(p)...)
}
