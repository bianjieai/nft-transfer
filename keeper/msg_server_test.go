package keeper_test

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bianjieai/nft-transfer/types"
)

var govAcc = authtypes.NewEmptyModuleAccount(govtypes.ModuleName, authtypes.Minter)

func (suite *KeeperTestSuite) TestMsgUpdateParams() {
	// default params
	params := types.DefaultParams()
	nftTransferKeeper := suite.GetSimApp(suite.chainA).NFTTransferKeeper

	testCases := []struct {
		name      string
		input     *types.MsgUpdateParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid authority",
			input: &types.MsgUpdateParams{
				Authority: "invalid",
				Params:    params,
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "send enabled param",
			input: &types.MsgUpdateParams{
				Authority: nftTransferKeeper.GetAuthority(),
				Params: types.Params{
					SendEnabled:    true,
					ReceiveEnabled: false,
				},
			},
			expErr: false,
		},
		{
			name: "receive enabled param",
			input: &types.MsgUpdateParams{
				Authority: nftTransferKeeper.GetAuthority(),
				Params: types.Params{
					SendEnabled:    false,
					ReceiveEnabled: true,
				},
			},
			expErr: false,
		},
		{
			name: "all enabled",
			input: &types.MsgUpdateParams{
				Authority: nftTransferKeeper.GetAuthority(),
				Params: types.Params{
					SendEnabled:    true,
					ReceiveEnabled: true,
				},
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			_, err := nftTransferKeeper.UpdateParams(suite.chainA.GetContext(), tc.input)

			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
				actParams := nftTransferKeeper.GetParams(suite.chainA.GetContext())
				suite.Require().EqualValues(tc.input.Params, actParams, "not equal")
			}
		})
	}
}
