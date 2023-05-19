package keeper

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	host "github.com/cosmos/ibc-go/v5/modules/core/24-host"

	"github.com/bianjieai/nft-transfer/types"
)

// Keeper defines the IBC non fungible transfer keeper
type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.Codec
	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority        string
	router           *types.Router
	defaultNFTKeeper types.NFTKeeper

	ics4Wrapper   types.ICS4Wrapper
	channelKeeper types.ChannelKeeper
	portKeeper    types.PortKeeper
	authKeeper    types.AccountKeeper
	scopedKeeper  capabilitykeeper.ScopedKeeper
}

// NewKeeper creates a new IBC nft-transfer Keeper instance
func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	authority string,
	defaultNFTKeeper types.NFTKeeper,
	ics4Wrapper types.ICS4Wrapper,
	channelKeeper types.ChannelKeeper,
	portKeeper types.PortKeeper,
	authKeeper types.AccountKeeper,
	scopedKeeper capabilitykeeper.ScopedKeeper,
) Keeper {
	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		router:           types.NewRouter(),
		authority:        authority,
		defaultNFTKeeper: defaultNFTKeeper,
		ics4Wrapper:      ics4Wrapper,
		channelKeeper:    channelKeeper,
		portKeeper:       portKeeper,
		authKeeper:       authKeeper,
		scopedKeeper:     scopedKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"-"+types.ModuleName)
}

// WithRouter set the router and return the Keeper
func (k Keeper) WithRouter(router *types.Router) Keeper {
	k.router = router
	return k
}

// SetPort sets the portID for the nft-transfer module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPort(portID), []byte(portID))
}

// GetPort returns the portID for the nft-transfer module.
func (k Keeper) GetPort(ctx sdk.Context, portID string) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.KeyPort(portID)))
}

// GetPort returns the portID for the nft-transfer module.
func (k Keeper) GetPorts(ctx sdk.Context) (ports []string) {
	store := ctx.KVStore(k.storeKey)
	portStore := prefix.NewStore(store, types.PortKey)

	iterator := portStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		ports = append(ports, string(iterator.Value()))
	}
	return ports
}

// IsBound checks if the transfer module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(
	ctx sdk.Context,
	cap *capabilitytypes.Capability,
	name string,
) bool {
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the nft-transfer module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(
	ctx sdk.Context,
	cap *capabilitytypes.Capability,
	name string,
) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

// SetEscrowAddress attempts to save a account to auth module
func (k Keeper) SetEscrowAddress(ctx sdk.Context, portID, channelID string) {
	// create the escrow address for the tokens
	escrowAddress := types.GetEscrowAddress(portID, channelID)
	if !k.authKeeper.HasAccount(ctx, escrowAddress) {
		acc := k.authKeeper.NewAccountWithAddress(ctx, escrowAddress)
		k.authKeeper.SetAccount(ctx, acc)
	}
}

// GetNFTKeeper return the keeper corresponding to the port
func (k Keeper) GetNFTKeeper(port string) (types.NFTKeeper, error) {
	nftKeeper, ok := k.router.GetRoute(port)
	if ok {
		return nftKeeper, nil

	}
	if k.defaultNFTKeeper != nil {
		return k.defaultNFTKeeper, nil
	}
	return nil, sdkerrors.Wrapf(types.ErrNotRegisterRoute, "port: %s", port)
}

// GetNFTKeeper return the keeper corresponding to the port
func (k Keeper) HasRoute(port string) bool {
	_, ok := k.router.GetRoute(port)
	return ok
}
