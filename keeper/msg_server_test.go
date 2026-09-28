package keeper_test

import (
	"cosmossdk.io/x/nft"
	ibctesting "github.com/bianjieai/nft-transfer/testing"
	"github.com/bianjieai/nft-transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
)

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

// IBC v10 removed capability checks. A user-supplied port must still be owned
// by NFT transfer, even when another IBC application's channel is open.
func (suite *KeeperTestSuite) TestTransferRejectsForeignPort() {
	path := ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)
	ctx := suite.chainA.GetContext()
	app := suite.GetSimApp(suite.chainA)
	sender := suite.chainA.SenderAccount.GetAddress()
	const classID, tokenID = "port-check", "nft"
	suite.Require().NoError(app.NFTKeeper.SaveClass(ctx, nft.Class{Id: classID}))
	suite.Require().NoError(app.NFTKeeper.Mint(ctx, nft.NFT{ClassId: classID, Id: tokenID}, sender))
	sequence, found := app.IBCKeeper.ChannelKeeper.GetNextSequenceSend(ctx, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
	suite.Require().True(found)

	_, err := app.NFTTransferKeeper.Transfer(ctx, &types.MsgTransfer{
		SourcePort:    path.EndpointA.ChannelConfig.PortID,
		SourceChannel: path.EndpointA.ChannelID,
		ClassId:       classID,
		TokenIds:      []string{tokenID},
		Sender:        sender.String(),
		Receiver:      suite.chainB.SenderAccount.GetAddress().String(),
		TimeoutHeight: suite.chainB.GetTimeoutHeight(),
	})
	suite.Require().ErrorIs(err, porttypes.ErrInvalidPort)
	suite.Require().Equal(sender, app.NFTKeeper.GetOwner(ctx, classID, tokenID))
	nextSequence, found := app.IBCKeeper.ChannelKeeper.GetNextSequenceSend(ctx, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
	suite.Require().True(found)
	suite.Require().Equal(sequence, nextSequence)
	suite.Require().Empty(app.IBCKeeper.ChannelKeeper.GetPacketCommitment(ctx, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sequence))
}

func (suite *KeeperTestSuite) TestCallbacksRejectForeignPort() {
	ctx := suite.chainA.GetContext()
	app := suite.GetSimApp(suite.chainA)
	relayer := suite.chainA.SenderAccount.GetAddress()
	const foreignPort = "other-nft-transfer"
	// A pre-existing foreign port can resolve to NFT transfer through IBC v10's
	// substring routing even though a new handshake would reject that port.
	module, found := app.IBCKeeper.PortKeeper.Route(foreignPort)
	suite.Require().True(found)
	escrowAddress := types.GetEscrowAddress(foreignPort, "channel-0")
	suite.Require().False(app.AccountKeeper.HasAccount(ctx, escrowAddress))
	_, err := module.OnChanOpenInit(ctx, channeltypes.UNORDERED, nil, foreignPort, "channel-0", channeltypes.Counterparty{}, types.Version)
	suite.Require().ErrorIs(err, porttypes.ErrInvalidPort)
	_, err = module.OnChanOpenTry(ctx, channeltypes.UNORDERED, nil, foreignPort, "channel-0", channeltypes.Counterparty{}, types.Version)
	suite.Require().ErrorIs(err, porttypes.ErrInvalidPort)
	suite.Require().ErrorIs(module.OnChanOpenAck(ctx, foreignPort, "channel-0", "channel-1", types.Version), porttypes.ErrInvalidPort)
	suite.Require().False(app.AccountKeeper.HasAccount(ctx, escrowAddress))
	suite.Require().ErrorIs(module.OnChanOpenConfirm(ctx, foreignPort, "channel-0"), porttypes.ErrInvalidPort)
	suite.Require().False(app.AccountKeeper.HasAccount(ctx, escrowAddress))
	suite.Require().ErrorIs(module.OnChanCloseConfirm(ctx, foreignPort, "channel-0"), porttypes.ErrInvalidPort)
	packet := channeltypes.Packet{SourcePort: foreignPort, DestinationPort: foreignPort}
	ack := module.OnRecvPacket(ctx, types.Version, packet, relayer)
	suite.Require().Equal(channeltypes.NewErrorAcknowledgement(porttypes.ErrInvalidPort).Acknowledgement(), ack.Acknowledgement())
	suite.Require().ErrorIs(module.OnAcknowledgementPacket(ctx, types.Version, packet, nil, relayer), porttypes.ErrInvalidPort)
	suite.Require().ErrorIs(module.OnTimeoutPacket(ctx, types.Version, packet, relayer), porttypes.ErrInvalidPort)
}
