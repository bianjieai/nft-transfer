package keeper_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	ibctesting "github.com/bianjieai/nft-transfer/testing"
	ics721testing "github.com/bianjieai/nft-transfer/testing"
	"github.com/bianjieai/nft-transfer/testing/mock"
	"github.com/bianjieai/nft-transfer/testing/simapp"
	"github.com/bianjieai/nft-transfer/types"
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ics721testing.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	chainC *ibctesting.TestChain

	queryClient types.QueryClient

	classMetadata *codectypes.Any
	tokenMetadata *codectypes.Any
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ics721testing.NewCoordinator(suite.T(), 3)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))
	suite.chainC = suite.coordinator.GetChain(ibctesting.GetChainID(3))

	queryHelper := baseapp.NewQueryServerTestHelper(suite.chainA.GetContext(), suite.GetSimApp(suite.chainA).InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.GetSimApp(suite.chainA).NFTTransferKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	classData := &mock.ClassMetadata{
		Creator:          "test creator",
		Schema:           "test schema",
		MintRestricted:   true,
		UpdateRestricted: true,
		Data:             "test data",
	}
	classMetadata, err := codectypes.NewAnyWithValue(classData)
	suite.Require().NoError(err, "NewAnyWithValue error")

	suite.classMetadata = classMetadata

	tokenData := &mock.TokenMetadata{
		Name: "kitty",
		Data: "test data",
	}
	tokenMetadata, err := codectypes.NewAnyWithValue(tokenData)
	suite.Require().NoError(err, "NewAnyWithValue error")

	suite.tokenMetadata = tokenMetadata
}

func (suite *KeeperTestSuite) MarshalClassMetadata() string {
	codec := suite.chainA.App.AppCodec()
	bz, err := codec.MarshalJSON(suite.classMetadata)
	suite.Require().NoError(err, "MarshalClassMetadata error")
	return base64.RawStdEncoding.EncodeToString(bz)
}

func (suite *KeeperTestSuite) MarshalTokenMetadata() string {
	codec := suite.chainA.App.AppCodec()
	bz, err := codec.MarshalJSON(suite.tokenMetadata)
	suite.Require().NoError(err, "MarshalTokenMetadata error")
	return base64.RawStdEncoding.EncodeToString(bz)
}

func (suite *KeeperTestSuite) GetSimApp(chain *ibctesting.TestChain) *simapp.SimApp {
	app := chain.App.(*simapp.SimApp)
	return app
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
