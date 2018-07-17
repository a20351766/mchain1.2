// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/hyperledger/mchain/common/localmsp"
	"github.com/hyperledger/mchain/common/tools/configtxgen/encoder"
	genesisconfig "github.com/hyperledger/mchain/common/tools/configtxgen/localconfig"
	cb "github.com/hyperledger/mchain/protos/common"
)

func newChainRequest(consensusType, creationPolicy, newChannelID string) *cb.Envelope {
	env, err := encoder.MakeChannelCreationTransaction(newChannelID, localmsp.NewSigner(), nil, genesisconfig.Load(genesisconfig.SampleSingleMSPChannelProfile))
	if err != nil {
		panic(err)
	}
	return env
}
