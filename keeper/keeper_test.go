package keeper_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	ibctesting "github.com/bianjieai/nft-transfer/testing"
	"github.com/bianjieai/nft-transfer/testing/mock"
	"github.com/bianjieai/nft-transfer/types"
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	chainC *ibctesting.TestChain

	queryClient types.QueryClient

	any         *codectypes.Any
	nftMetadata []byte
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 3)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))
	suite.chainC = suite.coordinator.GetChain(ibctesting.GetChainID(3))

	queryHelper := baseapp.NewQueryServerTestHelper(suite.chainA.GetContext(), suite.chainA.GetSimApp().InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.chainA.GetSimApp().NFTTransferKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	tokenData := &mock.TokenMetadata{
		Name:                 "kitty",
		Description:          "fertile digital cats",
		Image:                "external-link-url/image.png",
		ExternalLink:         "external-link-url/image.png",
		SellerFeeBasisPoints: "100",
	}
	any, err := codectypes.NewAnyWithValue(tokenData)
	suite.Require().NoError(err, "NewAnyWithValue error")
	suite.any = any

	suite.nftMetadata, err = suite.chainA.GetSimApp().NFTTransferKeeper.TokenDataResolver().Marshal(any)
	suite.Require().NoError(err, "MarshalAny error")
}

func (suite *KeeperTestSuite) TestMarshalAnyAndUnmarshalAny() {
	var data = []byte(`{"image":null,"image_data":null,"external_url":null,"description":"acme","name":null,"attributes":null,"background_color":null,"animation_url":null,"youtube_url":null}`)

	exp, err := codectypes.NewAnyWithValue(&types.UnknownTokenData{Data: data})
	suite.Require().NoError(err, "NewAnyWithValue error")

	any, err := suite.chainA.GetSimApp().NFTTransferKeeper.TokenDataResolver().Unmarshal(data)
	suite.Require().NoError(err, "UnmarshalAny error")
	suite.Require().True(reflect.DeepEqual(exp, any), "not equal")

	bz, err := suite.chainA.GetSimApp().NFTTransferKeeper.TokenDataResolver().Marshal(any)
	suite.Require().NoError(err, "MarshalAny error")
	suite.Require().Equal(data, bz, "MarshalAny failed")
}

func NewTransferPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = types.PortID
	path.EndpointB.ChannelConfig.PortID = types.PortID
	path.EndpointA.ChannelConfig.Version = types.Version
	path.EndpointB.ChannelConfig.Version = types.Version

	return path
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
