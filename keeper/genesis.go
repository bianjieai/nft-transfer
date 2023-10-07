package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/nft-transfer/types"
)

// InitGenesis initializes the ibc nft-transfer state and binds to PortID.
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	for _, trace := range state.Traces {
		k.SetClassTrace(ctx, trace)
	}

	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	for _, portID := range state.PortIds {
		k.SetPort(ctx, portID)

		if !k.IsBound(ctx, portID) {
			// nft-transfer module binds to the nft-transfer port on InitChain
			// and claims the returned capability
			err := k.BindPort(ctx, portID)
			if err != nil {
				panic(fmt.Sprintf("could not claim port capability: %v", err))
			}
		}
	}

	if err := k.SetParams(ctx, state.Params);err != nil {
		panic(fmt.Sprintf("SetParams failed: %v", err))
	}
}

// ExportGenesis exports ibc nft-transfer  module's portID and class trace info into its genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		PortIds: k.GetPorts(ctx),
		Traces:  k.GetAllClassTraces(ctx),
		Params:  k.GetParams(ctx),
	}
}
