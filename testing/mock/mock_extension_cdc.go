package mock

import (
	"encoding/base64"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	proto "github.com/gogo/protobuf/proto"

	"github.com/bianjieai/nft-transfer/types"
)

const ExtensionTypeURL = "/mock.Extension"

type (
	ClassMetadataResolver interface {
		MarshalClassMetadata(any *codectypes.Any) (string, error)
		UnmarshalClassMetadata(bz string) (*codectypes.Any, error)
	}

	TokenMetadataResolver interface {
		MarshalTokenMetadata(any *codectypes.Any) (string, error)
		UnmarshalTokenMetadata(bz string) (*codectypes.Any, error)
	}
)

// MarshalClassMetadata is responsible for serializing ClassMetadata
func (w MockNFTKeeper) MarshalClassMetadata(any *codectypes.Any) (string, error) {
	if any == nil {
		return "", nil
	}

	if any.GetTypeUrl() != ExtensionTypeURL {
		bz, err := w.cdc.MarshalJSON(any)
		if err != nil {
			return "", err
		}
		return base64.RawStdEncoding.EncodeToString(bz), nil
	}

	var message proto.Message
	if err := w.cdc.UnpackAny(any, &message); err != nil {
		return "", err
	}
	tokenData, ok := message.(*Extension)
	if !ok {
		return "", types.ErrMarshal
	}
	return tokenData.Data, nil
}

// UnmarshalClassMetadata is responsible for deserializing Metadata.
// Notice: If it is an unregistered type in this system, this method will save the original data in the `Extension` type,
// which is compatible with the definition of other chain Metadata types.
func (w MockNFTKeeper) UnmarshalClassMetadata(data string) (*codectypes.Any, error) {
	if len(data) == 0 {
		return nil, nil
	}

	bz, err := base64.RawStdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	var any codectypes.Any
	if err := w.cdc.UnmarshalJSON(bz, &any); err == nil {
		return &any, nil
	}
	return codectypes.NewAnyWithValue(&Extension{Data: data})
}

// MarshalTokenMetadata is responsible for serializing tokendata
func (w MockNFTKeeper) MarshalTokenMetadata(any *codectypes.Any) (string, error) {
	if any == nil {
		return "", nil
	}

	if any.GetTypeUrl() != ExtensionTypeURL {
		bz, err := w.cdc.MarshalJSON(any)
		if err != nil {
			return "", err
		}
		return base64.RawStdEncoding.EncodeToString(bz), nil
	}

	var message proto.Message
	if err := w.cdc.UnpackAny(any, &message); err != nil {
		return "", err
	}
	tokenData, ok := message.(*Extension)
	if !ok {
		return "", types.ErrMarshal
	}
	return tokenData.Data, nil

}

// Unmarshal is responsible for deserializing tokendata.
// Notice: If it is an unregistered type in this system, this method will save the original data in the `UnknownTokenData` type,
// which is compatible with the definition of other chain Tokendata types.
func (w MockNFTKeeper) UnmarshalTokenMetadata(data string) (*codectypes.Any, error) {
	if len(data) == 0 {
		return nil, nil
	}

	bz, err := base64.RawStdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	var any codectypes.Any
	if err := w.cdc.UnmarshalJSON(bz, &any); err == nil {
		return &any, nil
	}

	return codectypes.NewAnyWithValue(&Extension{Data: data})
}
