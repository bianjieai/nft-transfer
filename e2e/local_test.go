package e2e

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	interchaintest "github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/strangelove-ventures/interchaintest/v6/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	port    = "nft-transfer"
	version = "ics721-1"
)

func TestLocal(t *testing.T) {
	LocalTest(t, "irishub")
}

func LocalTest(t *testing.T, chainName string) {
	t.Parallel()

	chainConfig1 := ibc.ChainConfig{
		Name:    chainName,
		Type:    "cosmos",
		ChainID: "irishub-1",
		Images: []ibc.DockerImage{
			{
				Repository: "zhiqiangz/irishub",  // FOR LOCAL IMAGE USE: Docker Image Name
				Version:    "v1.4.1-gon-testnet", // FOR LOCAL IMAGE USE: Docker Image Tag
			},
		},
		Bin:            "iris",
		Bech32Prefix:   "iaa",
		Denom:          "stake",
		GasPrices:      "0.001stake",
		GasAdjustment:  1.3,
		TrustingPeriod: "508h",
		NoHostMount:    false,
	}

	chainConfig2 := ibc.ChainConfig{
		Name:    chainName,
		Type:    "cosmos",
		ChainID: "irishub-2",
		Images: []ibc.DockerImage{
			{
				Repository: "zhiqiangz/irishub",  // FOR LOCAL IMAGE USE: Docker Image Name
				Version:    "v1.4.1-gon-testnet", // FOR LOCAL IMAGE USE: Docker Image Tag
			},
		},
		Bin:            "iris",
		Bech32Prefix:   "iaa",
		Denom:          "stake",
		GasPrices:      "0.001stake",
		GasAdjustment:  1.3,
		TrustingPeriod: "508h",
		NoHostMount:    false,
	}

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{ChainName: chainName + "1", Version: "v1.4.1-gon-testnet", ChainConfig: chainConfig1},
		{ChainName: chainName + "2", Version: "v1.4.1-gon-testnet", ChainConfig: chainConfig2},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)
	chain1, chain2 := chains[0], chains[1]

	// Relayer Factory
	client, network := interchaintest.DockerSetup(t)
	rly := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
	).Build(t, client, network)

	// Prep Interchain
	const ibcPath = "gaia-osmo-demo"
	ic := interchaintest.NewInterchain().
		AddChain(chain1).
		AddChain(chain2).
		AddRelayer(rly, "relayer").
		AddLink(interchaintest.InterchainLink{
			Chain1:  chain1,
			Chain2:  chain2,
			Relayer: rly,
			Path:    ibcPath,
			CreateChannelOpts: ibc.CreateChannelOptions{
				SourcePortName: port,
				DestPortName:   port,
				Order:          ibc.Unordered,
				Version:        version,
			},
		})

	ctx := context.Background()

	// Log location
	f, err := interchaintest.CreateLogFile(fmt.Sprintf("%d.json", time.Now().Unix()))
	require.NoError(t, err)
	// Reporter/logs
	rep := testreporter.NewReporter(f)
	eRep := rep.RelayerExecReporter(t)

	// start chains, rly
	err = ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation:  false,
	})

	require.NoError(t, err)
	t.Cleanup(func() {
		_ = ic.Close()
	})

	const userFunds = int64(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, chain1, chain2)
	chain1User := users[0]
	chain2User := users[1]

	chain1BalInitial, err := chain1.GetBalance(ctx, chain1User.FormattedAddress(), chain1.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, userFunds, chain1BalInitial)

	// Get Channel ID
	chain1ChannelInfo, err := rly.GetChannels(ctx, eRep, chain1.Config().ChainID)
	require.NoError(t, err)
	chain1ChannelID := chain1ChannelInfo[0].ChannelID

	chain2ChannelInfo, err := rly.GetChannels(ctx, eRep, chain2.Config().ChainID)
	require.NoError(t, err)
	chain2ChannelID := chain2ChannelInfo[0].ChannelID

	classID := "badbody"
	cmd1 := classIssueCmd(classID, "", "", "", "", "", true, true)
	chain1Node := chain1.(*cosmos.CosmosChain).FullNodes[0]
	hash, err := chain1Node.ExecTx(ctx, chain1User.KeyName(), cmd1...)
	require.NoError(t, err)
	t.Logf("issue class txHash=%s", hash)

	tokenID := "tom"
	cmd2 := tokenMintCmd(classID, tokenID, "", "", "", "", "")
	hash, err = chain1Node.ExecTx(ctx, chain1User.KeyName(), cmd2...)
	require.NoError(t, err)
	t.Logf("mint token txHash=%s", hash)

	output, _, err := chain1Node.ExecQuery(ctx, tokenQueryCmd(classID, tokenID)...)
	require.NoError(t, err)
	t.Logf("token: %s", string(output))

	cmd3 := tokenInterTransferCmd(port, chain1ChannelID, chain2User.FormattedAddress(), classID, tokenID)
	hash, err = chain1Node.ExecTx(ctx, chain1User.KeyName(), cmd3...)
	require.NoError(t, err)
	t.Logf("inter transfer token txHash=%s", hash)

	// relay packets and acknoledgments
	require.NoError(t, rly.FlushPackets(ctx, eRep, ibcPath, chain2ChannelID))
	require.NoError(t, rly.FlushAcknowledgements(ctx, eRep, ibcPath, chain1ChannelID))

	chain2Node := chain2.(*cosmos.CosmosChain).FullNodes[0]
	traceHash := sha256.Sum256([]byte(fmt.Sprintf("%s/%s/%s", port, chain2ChannelID, classID)))
	ibcClassID := fmt.Sprintf("%s/%s", "ibc", strings.ToUpper(hex.EncodeToString(traceHash[:])))
	output, _, err = chain2Node.ExecQuery(ctx, tokenQueryCmd(ibcClassID, tokenID)...)
	require.NoError(t, err)
	t.Logf("ibc token: %s", string(output))
}
