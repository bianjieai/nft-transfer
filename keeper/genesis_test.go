package keeper_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/bianjieai/nft-transfer/types"
	"github.com/stretchr/testify/suite"
)

func TestKeeperTestSuite3(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestGenesis() {
	var (
		path   string
		traces types.Traces
	)

	for i := 0; i < 5; i++ {
		prefix := fmt.Sprintf("nft-transfer/channelToChain%d", i)
		if i == 0 {
			path = prefix
		} else {
			path = prefix + "/" + path
		}

		classTrace := types.ClassTrace{
			BaseClassId: "kitty",
			Path:        path,
		}
		traces = append(types.Traces{classTrace}, traces...)
		suite.chainA.GetSimApp().NFTTransferKeeper.SetClassTrace(suite.chainA.GetContext(), classTrace)
	}

	genesis := suite.chainA.GetSimApp().NFTTransferKeeper.ExportGenesis(suite.chainA.GetContext())

	//clone DefaultPorts
	var expPorts []string
	expPorts = append(expPorts, types.DefaultPorts...)

	sort.Strings(expPorts)
	suite.Require().EqualValues(expPorts, genesis.PortIds)
	suite.Require().Equal(traces.Sort(), genesis.Traces)

	suite.Require().NotPanics(func() {
		suite.chainA.GetSimApp().NFTTransferKeeper.InitGenesis(suite.chainA.GetContext(), *genesis)
	})
}
