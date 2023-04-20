package simulation_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/stretchr/testify/require"

	"github.com/bianjieai/nft-transfer/simulation"
	"github.com/bianjieai/nft-transfer/testing/simapp"
	"github.com/bianjieai/nft-transfer/types"
)

func TestDecodeStore(t *testing.T) {
	app := simapp.Setup(false)
	dec := simulation.NewDecodeStore(app.NFTTransferKeeper)

	trace := types.ClassTrace{
		BaseClassId: "kitty",
		Path:        "nft-transfer/channelToA",
	}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{
				Key:   types.PortKey,
				Value: []byte(types.NativePortID),
			},
			{
				Key:   types.ClassTraceKey,
				Value: app.NFTTransferKeeper.MustMarshalClassTrace(trace),
			},
			{
				Key:   []byte{0x99},
				Value: []byte{0x99},
			},
		},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"PortID", fmt.Sprintf("Port A: %s\nPort B: %s", types.NativePortID, types.NativePortID)},
		{"ClassTrace", fmt.Sprintf("ClassTrace A: %s\nClassTrace B: %s", trace.IBCClassID(), trace.IBCClassID())},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			if i == len(tests)-1 {
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			} else {
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
