package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
)

// Class defines the interface specifications of collection that can be transferred across chains
type Class interface {
	GetID() string
	GetURI() string
	GetData() string
}

// NFT defines the interface specification of nft that can be transferred across chains
type NFT interface {
	GetClassID() string
	GetID() string
	GetURI() string
	GetData() string
}

// NFTKeeper defines the expected nft keeper
type NFTKeeper interface {
	CreateOrUpdateClass(ctx sdk.Context, classID, classURI string, classData string) error
	Mint(ctx sdk.Context, classID, tokenID, tokenURI string, tokenData string, receiver sdk.AccAddress) error
	Transfer(ctx sdk.Context, classID string, tokenID string, tokenData string, receiver sdk.AccAddress) error
	Burn(ctx sdk.Context, classID string, tokenID string) error

	GetOwner(ctx sdk.Context, classID string, tokenID string) sdk.AccAddress
	HasClass(ctx sdk.Context, classID string) bool
	GetClass(ctx sdk.Context, classID string) (Class, bool)
	GetNFT(ctx sdk.Context, classID, tokenID string) (NFT, bool)
}

// ICS4Wrapper defines the expected ICS4Wrapper for middleware
type ICS4Wrapper interface {
	SendPacket(
		ctx sdk.Context,
		chanCap *capabilitytypes.Capability,
		sourcePort string,
		sourceChannel string,
		timeoutHeight clienttypes.Height,
		timeoutTimestamp uint64,
		data []byte,
	) (sequence uint64, err error)
}

// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channeltypes.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
}

// PortKeeper defines the expected IBC port keeper
type PortKeeper interface {
	BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability
}

// AccountKeeper defines the contract required for account APIs.
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Set an account in the store.
	SetAccount(sdk.Context, types.AccountI)
	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
	GetModuleAddress(name string) sdk.AccAddress
}
