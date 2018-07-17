/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig_test

import (
	"testing"

	newchannelconfig "github.com/hyperledger/mchain/common/channelconfig"
	"github.com/hyperledger/mchain/common/tools/configtxgen/configtxgentest"
	"github.com/hyperledger/mchain/common/tools/configtxgen/encoder"
	genesisconfig "github.com/hyperledger/mchain/common/tools/configtxgen/localconfig"
	"github.com/hyperledger/mchain/protos/utils"

	"github.com/stretchr/testify/assert"
)

func TestWithRealConfigtx(t *testing.T) {
	conf := configtxgentest.Load(genesisconfig.SampleSingleMSPSoloProfile)

	// None of the sample profiles define an application config section
	// in a genesis block (as this is a bad idea), but we combine them
	// here to better exercise the code.
	conf.Application = &genesisconfig.Application{
		Organizations: []*genesisconfig.Organization{
			conf.Orderer.Organizations[0],
		},
	}
	conf.Application.Organizations[0].AnchorPeers = []*genesisconfig.AnchorPeer{
		{
			Host: "foo",
			Port: 7,
		},
	}
	gb := encoder.New(conf).GenesisBlockForChannel("foo")
	env := utils.ExtractEnvelopeOrPanic(gb, 0)
	_, err := newchannelconfig.NewBundleFromEnvelope(env)
	assert.NoError(t, err)
}
