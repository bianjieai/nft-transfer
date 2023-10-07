package types

import (
	"testing"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		name     string
		genState *GenesisState
		wantErr  bool
	}{
		{
			name:     "default",
			genState: DefaultGenesisState(),
			wantErr:  false,
		},
		{
			"valid genesis",
			&GenesisState{
				PortIds: []string{"portidone"},
			},
			false,
		},
		{
			"invalid port",
			&GenesisState{
				PortIds: []string{"(INVALIDPORT)"},
			},
			true,
		},
		{
			"duplicate port",
			&GenesisState{
				PortIds: []string{NativePortID, ERC721PortID, ERC721PortID},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.genState.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("GenesisState.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
