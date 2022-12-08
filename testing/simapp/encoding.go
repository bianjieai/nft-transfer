package simapp

import (
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/gogo/protobuf/proto"

	"github.com/bianjieai/nft-transfer/testing/mock"
	simappparams "github.com/bianjieai/nft-transfer/testing/simapp/params"
)

// MakeTestEncodingConfig creates an EncodingConfig for testing. This function
// should be used only in tests or when creating a new app instance (NewApp*()).
// App user shouldn't create new codecs - use the app.AppCodec instead.
// [DEPRECATED]
func MakeTestEncodingConfig() simappparams.EncodingConfig {
	encodingConfig := simappparams.MakeTestEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	encodingConfig.InterfaceRegistry.RegisterImplementations((*proto.Message)(nil), &mock.TokenMetadata{})
	return encodingConfig
}
