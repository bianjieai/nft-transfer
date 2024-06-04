package types

import (
	"bytes"

	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/cosmos/gogoproto/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary nft-transfer interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgTransfer{}, "cosmos-sdk/MsgTransferNFT", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "cosmos-sdk/MsgUpdateParams", nil)
}

// RegisterInterfaces register the ibc nft-transfer module interfaces to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTransfer{},
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	// ModuleCdc references the global nft-transfer module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding.
	//
	// The actual codec used for serialization should be provided to nft-transfer and
	// defined at the application level.
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

	// AminoCdc is a amino codec created to support amino json compatible msgs.
	AminoCdc = codec.NewLegacyAmino()
)

func init() {
	RegisterLegacyAminoCodec(AminoCdc)
	AminoCdc.Seal()
}

// MustProtoMarshalJSON marshals a protobuf message to JSON and panics if there is an error.
//
// It takes a protobuf message as input and returns the JSON-encoded byte array.
// The function uses the provided anyResolver to resolve any protobuf Any types in the message.
// If there is an error during the marshaling process, the function panics.
//
// Parameters:
// - msg: The protobuf message to be marshaled.
//
// Returns:
// - []byte: The JSON-encoded byte array.
func MustProtoMarshalJSON(msg proto.Message) []byte {
	anyResolver := codectypes.NewInterfaceRegistry()
	bz, err := ProtoMarshalJSON(msg, anyResolver)
	if err != nil {
		panic(err)
	}
	return bz
}

// ProtoMarshalJSON provides an auxiliary function to return Proto3 JSON encoded
// bytes of a message.
func ProtoMarshalJSON(msg proto.Message, resolver jsonpb.AnyResolver) ([]byte, error) {
	jm := &jsonpb.Marshaler{OrigName: false, EmitDefaults: false, AnyResolver: resolver}
	err := codectypes.UnpackInterfaces(msg, codectypes.ProtoJSONPacker{JSONPBMarshaler: jm})
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err := jm.Marshal(buf, msg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
