package keeper_test

import (
	"fmt"

	"github.com/bianjieai/nft-transfer/types"
)

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
		suite.GetSimApp(suite.chainA).NFTTransferKeeper.SetClassTrace(suite.chainA.GetContext(), classTrace)
	}

	genesis := suite.GetSimApp(suite.chainA).NFTTransferKeeper.ExportGenesis(suite.chainA.GetContext())

	suite.Require().Equal(types.PortID, genesis.PortId)
	suite.Require().Equal(traces.Sort(), genesis.Traces)

	suite.Require().NotPanics(func() {
		suite.GetSimApp(suite.chainA).NFTTransferKeeper.InitGenesis(suite.chainA.GetContext(), *genesis)
	})
}
