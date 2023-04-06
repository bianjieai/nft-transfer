package types

import (
	errorsmod "cosmossdk.io/errors"
)

// IBC transfer sentinel errors
var (
	ErrInvalidPacketTimeout = errorsmod.Register(ModuleName, 2, "invalid packet timeout")
	ErrInvalidVersion       = errorsmod.Register(ModuleName, 3, "invalid ICS721 version")
	ErrMaxTransferChannels  = errorsmod.Register(ModuleName, 4, "max nft-transfer channels")
	ErrInvalidClassID       = errorsmod.Register(ModuleName, 5, "invalid class id")
	ErrInvalidTokenID       = errorsmod.Register(ModuleName, 6, "invalid token id")
	ErrInvalidPacket        = errorsmod.Register(ModuleName, 7, "invalid non-fungible token packet")
	ErrTraceNotFound        = errorsmod.Register(ModuleName, 8, "classTrace trace not found")
	ErrMarshal              = errorsmod.Register(ModuleName, 9, "failed to marshal token data")
	ErrSendDisabled         = errorsmod.Register(ModuleName, 10, "non-fungible token transfers from this chain are disabled")
	ErrReceiveDisabled      = errorsmod.Register(ModuleName, 11, "non-fungible token transfers to this chain are disabled")
)
