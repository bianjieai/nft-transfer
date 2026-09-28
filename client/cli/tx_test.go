package cli

import (
	"testing"
	"time"

	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	"github.com/stretchr/testify/require"
)

func TestRelativeTimeoutTimestamp(t *testing.T) {
	now := time.Unix(100, 0)
	timestamp, err := relativeTimeoutTimestamp(clienttypes.ZeroHeight(), uint64(time.Minute), now)
	require.NoError(t, err)
	require.Equal(t, uint64(now.Add(time.Minute).UnixNano()), timestamp)

	_, err = relativeTimeoutTimestamp(clienttypes.NewHeight(0, 100), uint64(time.Minute), now)
	require.ErrorContains(t, err, "--absolute-timeouts")
	_, err = relativeTimeoutTimestamp(clienttypes.ZeroHeight(), 0, now)
	require.ErrorContains(t, err, "non-zero timestamp")
	_, err = relativeTimeoutTimestamp(clienttypes.ZeroHeight(), uint64(time.Minute), time.Unix(0, 0))
	require.ErrorContains(t, err, "local clock time")
}
