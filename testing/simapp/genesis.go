package simapp

import (
	"github.com/cosmos/cosmos-sdk/codec"

	ibcsimapp "github.com/cosmos/ibc-go/v7/testing/simapp"
)

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) ibcsimapp.GenesisState {
	return ModuleBasics.DefaultGenesis(cdc)
}
