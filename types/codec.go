package types

import (
	"bytes"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary nft-transfer interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgTransfer{}, "cosmos-sdk/MsgTransferNFT", nil)
}

// RegisterInterfaces register the ibc nft-transfer module interfaces to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgTransfer{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global nft-transfer module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding.
	//
	// The actual codec used for serialization should be provided to nft-transfer and
	// defined at the application level.
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

	// AminoCdc is a amino codec created to support amino json compatible msgs.
	AminoCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}

func mustProtoMarshalJSON(msg proto.Message) []byte {
	anyResolver := codectypes.NewInterfaceRegistry()

	// EmitDefaults is set to false to prevent marshalling of unpopulated fields (memo)
	// OrigName and the anyResovler match the fields the original SDK function would expect
	// in order to minimize changes.
	jm := &jsonpb.Marshaler{OrigName: false, EmitDefaults: false, AnyResolver: anyResolver}

	err := codectypes.UnpackInterfaces(msg, codectypes.ProtoJSONPacker{JSONPBMarshaler: jm})
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	if err := jm.Marshal(buf, msg); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func MarshalAny(any *codectypes.Any) ([]byte, error) {
	if any == nil {
		return nil, nil
	}
	bz, err := any.Marshal()
	if err != nil {
		return bz, err
	}
	return bz, nil
}

func UnmarshalAny(bz []byte) (*codectypes.Any, error) {
	if bz == nil || len(bz) == 0 {
		return nil, nil
	}
	any := &codectypes.Any{}
	err := any.Unmarshal(bz)
	return any, err
}
