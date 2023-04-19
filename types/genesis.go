package types

import (
	host "github.com/cosmos/ibc-go/v5/modules/core/24-host"
)

var DefaultPorts = []string{PortID, "erc-721"}

// NewGenesisState creates a new ibc nft-transfer GenesisState instance.
func NewGenesisState(portIDs []string, traces Traces, params Params) *GenesisState {
	return &GenesisState{
		PortIds: portIDs,
		Traces:  traces,
		Params:  params,
	}
}

// DefaultGenesisState returns a GenesisState with "nft-transfer" as the default PortID.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		PortIds: DefaultPorts,
		Traces:  Traces{},
		Params:  DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for _, port := range gs.PortIds {
		if err := host.PortIdentifierValidator(port); err != nil {
			return err
		}
	}
	return gs.Traces.Validate()
}
