package types

import (
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// DefaultRelativePacketTimeoutHeight is the default packet timeout height (in blocks) relative
	// to the current block height of the counterparty chain provided by the client state. The
	// timeout is disabled when set to 0.
	DefaultRelativePacketTimeoutHeight = "0-0"

	// DefaultRelativePacketTimeoutTimestamp is the default packet timeout timestamp (in nanoseconds)
	// relative to the current block timestamp of the counterparty chain provided by the client
	// state. The timeout is disabled when set to 0. The default is currently set to a 10 minute
	// timeout.
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

// NewNonFungibleTokenPacketData constructs a new NonFungibleTokenPacketData instance
func NewNonFungibleTokenPacketData(
	classID, classURI, classData string,
	tokenIDs, tokenURI []string,
	sender, receiver string,
	tokenData []string,
	memo string,
) NonFungibleTokenPacketData {
	return NonFungibleTokenPacketData{
		ClassId:   classID,
		ClassUri:  classURI,
		ClassData: classData,
		TokenIds:  tokenIDs,
		TokenUris: tokenURI,
		TokenData: tokenData,
		Sender:    sender,
		Receiver:  receiver,
		Memo:      memo,
	}
}

// ValidateBasic is used for validating the nft transfer.
// NOTE: The addresses formats are not validated as the sender and recipient can have different
// formats defined by their corresponding chains that are not known to IBC.
func (nftpd NonFungibleTokenPacketData) ValidateBasic() error {
	if strings.TrimSpace(nftpd.ClassId) == "" {
		return errorsmod.Wrap(ErrInvalidClassID, "classId cannot be blank")
	}

	if len(nftpd.TokenIds) == 0 {
		return errorsmod.Wrap(ErrInvalidTokenID, "tokenId cannot be empty")
	}

	seen := make(map[string]int64)
	for i, id := range nftpd.TokenIds {
		if strings.TrimSpace(id) == "" {
			return errorsmod.Wrap(ErrInvalidTokenID, "tokenId cannot be blank")
		}
		if j, exist := seen[id]; exist {
			return errorsmod.Wrapf(ErrInvalidTokenID, "the tokenId at positions %d and %d in the array are repeated", i, j)
		}
		seen[id] = int64(i)
	}

	if (len(nftpd.TokenUris) != 0) && len(nftpd.TokenIds) != len(nftpd.TokenUris) {
		return errorsmod.Wrap(ErrInvalidPacket, "the length of tokenUri must be 0 or the same as the length of TokenIds")
	}

	if (len(nftpd.TokenData) != 0) && (len(nftpd.TokenIds) != len(nftpd.TokenData)) {
		return errorsmod.Wrap(ErrInvalidPacket, "the length of tokenData must be 0 or the same as the length of TokenIds")
	}

	if strings.TrimSpace(nftpd.Sender) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "sender address cannot be blank")
	}

	if strings.TrimSpace(nftpd.Receiver) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "receiver address cannot be blank")
	}
	return nil
}

// GetBytes is a helper for serializing
func (nftpd NonFungibleTokenPacketData) GetBytes() []byte {
	// Format will reshape tokenUris and tokenData in NonFungibleTokenPacketData:
	// 1. if tokenUris/tokenData is ["","",""] or [], then set it to nil.
	// 2. if tokenUris/tokenData is ["a","b","c"] or ["a", "", "c"], then keep it.
	// NOTE: Only use this before sending pkg.
	if requireShape(nftpd.TokenUris) {
		nftpd.TokenUris = nil
	}

	if requireShape(nftpd.TokenData) {
		nftpd.TokenData = nil
	}
	return sdk.MustSortJSON(MustProtoMarshalJSON(&nftpd))
}

func GetIfExist(i int, data []string) string {
	if i < 0 || i >= len(data) {
		return ""
	}
	return data[i]
}

// requireShape checks if TokenUris/TokenData needs to be set as nil
func requireShape(contents []string) bool {
	if contents == nil {
		return false
	}

	// empty slice of string
	if len(contents) == 0 {
		return true
	}

	emptyStringCount := 0
	for _, v := range contents {
		if len(v) == 0 {
			emptyStringCount++
		}
	}
	return emptyStringCount == len(contents)
}
