package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// IBC transfer sentinel errors
var (
	ErrInvalidPacketTimeout = sdkerrors.Register(ModuleName, 2, "invalid packet timeout")
	ErrInvalidVersion       = sdkerrors.Register(ModuleName, 3, "invalid ICS721 version")
	ErrMaxTransferChannels  = sdkerrors.Register(ModuleName, 4, "max nft-transfer channels")
	ErrInvalidClassID       = sdkerrors.Register(ModuleName, 5, "invalid class id")
	ErrInvalidTokenID       = sdkerrors.Register(ModuleName, 6, "invalid token id")
	ErrInvalidPacket        = sdkerrors.Register(ModuleName, 7, "invalid non-fungible token packet")
	ErrTraceNotFound        = sdkerrors.Register(ModuleName, 8, "classTrace trace not found")
	ErrMarshal              = sdkerrors.Register(ModuleName, 9, "failed to marshal token data")
	ErrSendDisabled         = sdkerrors.Register(ModuleName, 10, "non-fungible token transfers from this chain are disabled")
	ErrReceiveDisabled      = sdkerrors.Register(ModuleName, 11, "non-fungible token transfers to this chain are disabled")
)
