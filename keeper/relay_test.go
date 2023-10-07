package keeper_test

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"

	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/bianjieai/nft-transfer/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
)

func (suite *KeeperTestSuite) TestSendTransfer() {
	var (
		path          *ibctesting.Path
		err           error
		classID       string
		timeoutHeight clienttypes.Height
	)

	baseClassID := "cryptoCat"
	classURI := "cat_uri"
	nftID := "kitty"
	nftURI := "kittt_uri"

	testCases := []struct {
		msg              string
		malleate         func()
		isAwayFromOrigin bool
	}{
		{
			"successful transfer from chainA to chainB",
			func() {
				suite.coordinator.CreateChannels(path)
				classID = baseClassID

				nftKeeper := suite.GetSimApp(path.EndpointA.Chain).NFTKeeper
				err = nftKeeper.SaveClass(path.EndpointA.Chain.GetContext(), nft.Class{
					Id:  classID,
					Uri: classURI,
				})
				suite.Require().NoError(err, "SaveClass error")

				err = nftKeeper.Mint(path.EndpointA.Chain.GetContext(), nft.NFT{
					ClassId: classID,
					Id:      nftID,
					Uri:     nftURI,
					Data:    suite.tokenMetadata,
				}, path.EndpointA.Chain.SenderAccount.GetAddress())
				suite.Require().NoError(err, "Mint error")
				timeoutHeight = path.EndpointB.Chain.GetTimeoutHeight()
			},
			true,
		},
		{
			"successful transfer from chainB to chainA",
			func() {
				suite.coordinator.CreateChannels(path)
				trace := types.ParseClassTrace(
					types.GetClassPrefix(
						path.EndpointB.ChannelConfig.PortID,
						path.EndpointB.ChannelID,
					) + baseClassID)
				suite.GetSimApp(path.EndpointB.Chain).NFTTransferKeeper.SetClassTrace(path.EndpointB.Chain.GetContext(), trace)

				classID = trace.IBCClassID()
				nftKeeper := suite.GetSimApp(path.EndpointB.Chain).NFTKeeper
				err = nftKeeper.SaveClass(path.EndpointB.Chain.GetContext(), nft.Class{
					Id:   classID,
					Uri:  classURI,
					Data: suite.classMetadata,
				})
				suite.Require().NoError(err, "SaveClass error")

				err = nftKeeper.Mint(path.EndpointB.Chain.GetContext(), nft.NFT{
					ClassId: classID,
					Id:      nftID,
					Uri:     nftURI,
				}, path.EndpointB.Chain.SenderAccount.GetAddress())
				suite.Require().NoError(err, "Mint error")
				timeoutHeight = path.EndpointA.Chain.GetTimeoutHeight()
			},
			false,
		},
	}

	kts := suite

	for _, tc := range testCases {
		tc := tc

		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			path = NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)
			tc.malleate()

			if !tc.isAwayFromOrigin {
				ctx := path.EndpointB.Chain.GetContext()

				_, err = kts.GetSimApp(path.EndpointB.Chain).NFTTransferKeeper.SendTransfer(
					ctx,
					path.EndpointB.ChannelConfig.PortID,
					path.EndpointB.ChannelID,
					classID,
					[]string{nftID},
					path.EndpointB.Chain.SenderAccount.GetAddress(),
					path.EndpointA.Chain.SenderAccount.GetAddress().String(),
					timeoutHeight,
					0,
					"memo",
				)
				suite.Require().NoError(err)

				suite.Require().False(
					kts.GetSimApp(path.EndpointB.Chain).NFTKeeper.HasNFT(ctx, classID, nftID),
					"burn nft failed",
				)
				return
			}

			ctx := path.EndpointA.Chain.GetContext()
			_, err = kts.GetSimApp(path.EndpointA.Chain).NFTTransferKeeper.SendTransfer(
				ctx,
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
				classID,
				[]string{nftID},
				path.EndpointA.Chain.SenderAccount.GetAddress(),
				path.EndpointB.Chain.SenderAccount.GetAddress().String(),
				timeoutHeight,
				0,
				"memo",
			)

			suite.Require().NoError(err)
			suite.Require().Equal(
				types.GetEscrowAddress(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID),
				kts.GetSimApp(path.EndpointA.Chain).NFTKeeper.GetOwner(ctx, classID, nftID),
				"escrow nft failed",
			)
		})
	}
}

func (suite *KeeperTestSuite) TestOnRecvPacket() {
	var (
		path              *ibctesting.Path
		trace             types.ClassTrace
		classID, receiver string
		nftIDs, nftURIs   []string
		nftMetaDatas      []string
	)

	baseClassID := "cryptoCat"
	classURI := "cat_uri"
	nftID := "kitty"
	nftURI := "kittt_uri"

	kts := suite

	testCases := []struct {
		msg              string
		malleate         func()
		isAwayFromOrigin bool // the receiving chain is the source of the coin originally
		expPass          bool
	}{
		{"success receive chain is away from origin chain", func() {}, true, true},
		{"success receive chain is not away from origin chain", func() {
			classID = types.GetClassPrefix(
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
			) + baseClassID

			err := kts.GetSimApp(kts.chainB).NFTKeeper.SaveClass(suite.chainB.GetContext(), nft.Class{
				Id:   baseClassID,
				Uri:  classURI,
				Data: suite.classMetadata,
			})
			kts.Require().NoError(err, "SaveClass failed")

			escrowAddress := types.GetEscrowAddress(
				path.EndpointB.ChannelConfig.PortID,
				path.EndpointB.ChannelID,
			)
			err = kts.GetSimApp(kts.chainB).NFTKeeper.Mint(suite.chainB.GetContext(), nft.NFT{
				ClassId: baseClassID,
				Id:      nftID,
				Uri:     nftURI,
				Data:    suite.tokenMetadata,
			}, escrowAddress)
			kts.Require().NoError(err, "Mint failed")

		}, false, true},
		{"empty classID", func() {
			classID = ""
		}, true, false},
		{"empty nftIDs", func() {
			nftIDs = nil
		}, true, false},
		{"empty nftURIs", func() {
			nftURIs = nil
		}, true, true},
		{"invalid receiver address", func() {
			receiver = "gaia1scqhwpgsmr6vmztaa7suurfl52my6nd2kmrudl"
		}, true, false},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			path = NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			classID = baseClassID
			receiver = suite.chainB.SenderAccount.GetAddress().String()
			nftIDs = []string{nftID}
			nftURIs = []string{nftURI}
			nftMetaDatas = []string{suite.MarshalTokenMetadata()}

			tc.malleate()

			trace = types.ParseClassTrace(classID)
			data := types.NewNonFungibleTokenPacketData(
				trace.GetFullClassPath(),
				classURI,
				suite.MarshalClassMetadata(),
				nftIDs,
				nftURIs,
				suite.chainA.SenderAccount.GetAddress().String(),
				receiver,
				nftMetaDatas,
				"memo",
			)

			packet := channeltypes.NewPacket(
				data.GetBytes(),
				1, //not check sequence
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
				path.EndpointB.ChannelConfig.PortID,
				path.EndpointB.ChannelID,
				clienttypes.NewHeight(0, 100),
				0,
			)

			err := kts.GetSimApp(kts.chainB).
				NFTTransferKeeper.OnRecvPacket(suite.chainB.GetContext(), packet, data)

			if !tc.expPass {
				suite.Require().Error(err)
				return
			}

			if tc.isAwayFromOrigin {
				prefixedClassID := types.GetClassPrefix(
					path.EndpointB.ChannelConfig.PortID,
					path.EndpointB.ChannelID,
				) + baseClassID
				trace = types.ParseClassTrace(prefixedClassID)

				suite.Require().Equal(
					receiver,
					kts.GetSimApp(kts.chainB).NFTKeeper.GetOwner(suite.chainB.GetContext(), trace.IBCClassID(), nftID).String(),
					"receive packet failed",
				)

				suite.Require().True(
					kts.GetSimApp(kts.chainB).NFTTransferKeeper.HasClassTrace(suite.chainB.GetContext(), trace.Hash()),
					"not found class trace",
				)

			} else {
				suite.Require().False(
					kts.GetSimApp(kts.chainB).NFTKeeper.HasNFT(suite.chainB.GetContext(), classID, nftID),
					"burn nft failed")
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOnAcknowledgementPacket() {
	var (
		successAck      = channeltypes.NewResultAcknowledgement([]byte{byte(1)})
		failedAck       = channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "failed packet transfer"))
		path            *ibctesting.Path
		trace           types.ClassTrace
		classID         string
		nftIDs, nftURIs []string
	)

	baseClassID := "cryptoCat"
	classURI := "cat_uri"
	nftID := "kitty"
	nftURI := "kittt_uri"

	kts := suite

	testCases := []struct {
		msg      string
		ack      channeltypes.Acknowledgement
		malleate func()
		success  bool // success of ack
		expPass  bool
	}{
		{"success ack causes no-op", successAck, func() {}, true, true},
		{"successful refund when isAwayFromOrigin is false", failedAck, func() {
			// if isAwayFromOrigin is false, OnAcknowledgementPacket will mint nft to sender again

			// mock SendTransfer
			classID = types.GetClassPrefix(
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
			) + baseClassID

			ibcClassID := types.ParseClassTrace(classID).IBCClassID()
			err := kts.GetSimApp(kts.chainA).NFTKeeper.SaveClass(suite.chainA.GetContext(), nft.Class{
				Id:  ibcClassID,
				Uri: classURI,
			})
			kts.Require().NoError(err, "SaveClass failed")

		}, false, true},
		{"successful refund when isAwayFromOrigin is true", failedAck, func() {
			// if isAwayFromOrigin is true, OnAcknowledgementPacket will unescrow nft to sender

			// mock SendTransfer
			classID = types.GetClassPrefix(
				path.EndpointB.ChannelConfig.PortID,
				"channel-1",
			) + baseClassID

			ibcClassID := types.ParseClassTrace(classID).IBCClassID()
			err := kts.GetSimApp(kts.chainA).NFTKeeper.SaveClass(suite.chainA.GetContext(), nft.Class{
				Id:  ibcClassID,
				Uri: classURI,
			})
			kts.Require().NoError(err, "SaveClass failed")

			escrowAddress := types.GetEscrowAddress(
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
			)

			err = kts.GetSimApp(kts.chainA).NFTKeeper.Mint(suite.chainA.GetContext(), nft.NFT{
				ClassId: ibcClassID,
				Id:      nftID,
				Uri:     nftURI,
			}, escrowAddress)
			kts.Require().NoError(err, "Mint failed")

		}, false, true},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			path = NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			classID = baseClassID
			nftIDs = []string{nftID}
			nftURIs = []string{nftURI}

			tc.malleate()

			trace = types.ParseClassTrace(classID)
			data := types.NewNonFungibleTokenPacketData(
				trace.GetFullClassPath(),
				classURI,
				suite.MarshalClassMetadata(),
				nftIDs,
				nftURIs,
				suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainB.SenderAccount.GetAddress().String(),
				nil,
				"memo",
			)

			packet := channeltypes.NewPacket(
				data.GetBytes(),
				1, //not check sequence
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
				path.EndpointB.ChannelConfig.PortID,
				path.EndpointB.ChannelID,
				clienttypes.NewHeight(0, 100),
				0,
			)

			err := kts.GetSimApp(kts.chainA).NFTTransferKeeper.OnAcknowledgementPacket(
				suite.chainA.GetContext(),
				packet,
				data, tc.ack,
			)

			if !tc.expPass {
				suite.Require().Error(err)
			}

			suite.Require().NoError(err, "OnAcknowledgementPacket failed")
			if tc.success {
				// if successful, nft is hosted in account a or destroyed(executed when SendTransfer)
				return
			}

			suite.Require().Equal(
				suite.chainA.SenderAccount.GetAddress().String(),
				kts.GetSimApp(kts.chainA).NFTKeeper.GetOwner(suite.chainA.GetContext(), trace.IBCClassID(), nftID).String(),
				"refund failed",
			)
		})
	}
}
