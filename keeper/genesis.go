package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/nft-transfer/types"
)

// InitGenesis initializes the ibc nft-transfer state.
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	k.SetPort(ctx, state.PortId)

	for _, trace := range state.Traces {
		k.SetClassTrace(ctx, trace)
	}

	if err := k.SetParams(ctx, state.Params); err != nil {
		panic(fmt.Sprintf("SetParams failed: %v", err))
	}
}

// ExportGenesis exports ibc nft-transfer  module's portID and class trace info into its genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		PortId: k.GetPort(ctx),
		Traces: k.GetAllClassTraces(ctx),
		Params: k.GetParams(ctx),
	}
}
