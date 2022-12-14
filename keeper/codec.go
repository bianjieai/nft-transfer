package keeper

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gogo/protobuf/proto"

	"github.com/bianjieai/nft-transfer/types"
)

var _ TokenDataResolver = Keeper{}

type TokenDataResolver interface {
	Marshal(any *codectypes.Any) ([]byte, error)
	Unmarshal(bz []byte) (*codectypes.Any, error)
}

// Marshal is responsible for serializing tokendata
func (k Keeper) Marshal(any *codectypes.Any) ([]byte, error) {
	if any == nil {
		return nil, nil
	}

	if any.GetTypeUrl() != types.UnknownTokenDataTypeURL {
		return k.cdc.MarshalJSON(any)
	}

	var message proto.Message
	if err := k.cdc.UnpackAny(any, &message); err != nil {
		return nil, err
	}
	tokenData, ok := message.(*types.UnknownTokenData)
	if !ok {
		return nil, types.ErrMarshal
	}
	return tokenData.Data, nil

}

// Unmarshal is responsible for deserializing tokendata.
// Notice: If it is an unregistered type in this system, this method will save the original data in the `UnknownTokenData` type,
// which is compatible with the definition of other chain Tokendata types.
func (k Keeper) Unmarshal(bz []byte) (*codectypes.Any, error) {
	if bz == nil || len(bz) == 0 {
		return nil, nil
	}

	var any codectypes.Any
	if err := k.cdc.UnmarshalJSON(bz, &any); err == nil {
		return &any, nil
	}

	return codectypes.NewAnyWithValue(&types.UnknownTokenData{Data: bz})
}

// TokenDataResolver returns the parser of Tokendata
func (k Keeper) TokenDataResolver() TokenDataResolver {
	return k.resolver
}
