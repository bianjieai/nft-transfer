package types

import (
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/bianjieai/nft-transfer/testing/mock"
	"github.com/stretchr/testify/require"
)

func TestMarshalAnyAndUnmarshalAny(t *testing.T) {
	any, tokenDataBz := MockTokenMetadata()
	anyAct, err := UnmarshalAny(tokenDataBz)
	require.NoError(t, err)
	require.True(t, any.Equal(anyAct))
}

func MockTokenMetadata() (*codectypes.Any, []byte) {
	tokenData := &mock.TokenMetadata{
		Name:                 "kitty",
		Description:          "fertile digital cats",
		Image:                "external-link-url/image.png",
		ExternalLink:         "external-link-url/image.png",
		SellerFeeBasisPoints: "100",
	}
	any, err := codectypes.NewAnyWithValue(tokenData)
	if err != nil {
		panic(err)
	}

	bz, err := MarshalAny(any)
	if err != nil {
		panic(err)
	}
	return any, bz
}
