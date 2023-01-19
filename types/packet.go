package types

import (
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// DefaultRelativePacketTimeoutHeight is the default packet timeout height (in blocks) relative
	// to the current block height of the counterparty chain provided by the client state. The
	// timeout is disabled when set to 0.
	DefaultRelativePacketTimeoutHeight = "0-1000"

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
		return sdkerrors.Wrap(ErrInvalidClassID, "classId cannot be blank")
	}

	if len(nftpd.TokenIds) == 0 {
		return sdkerrors.Wrap(ErrInvalidTokenID, "tokenId cannot be empty")
	}

	for _, id := range nftpd.TokenIds {
		if strings.TrimSpace(id) == "" {
			return sdkerrors.Wrap(ErrInvalidTokenID, "tokenId cannot be blank")
		}
	}

	if (len(nftpd.TokenUris) != 0) && len(nftpd.TokenIds) != len(nftpd.TokenUris) {
		return sdkerrors.Wrap(ErrInvalidPacket, "the length of tokenUri must be 0 or the same as the length of TokenIds")
	}

	if _, err := ValidateContent(nftpd.TokenUris); err != nil {
		return err
	}

	if (len(nftpd.TokenData) != 0) && (len(nftpd.TokenIds) != len(nftpd.TokenData)) {
		return sdkerrors.Wrap(ErrInvalidPacket, "the length of tokenData must be 0 or the same as the length of TokenIds")
	}

	if _, err := ValidateContent(nftpd.TokenData); err != nil {
		return err
	}

	if strings.TrimSpace(nftpd.Sender) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender address cannot be blank")
	}

	if strings.TrimSpace(nftpd.Receiver) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "receiver address cannot be blank")
	}
	return nil
}

// ShapeContent will validate and reshape tokenUris and tokenData in NonFungibleTokenPacketData:
// 1. if tokenUris/tokenData is ["","",""], then set it to nil.
// 2. if tokenUris/tokenData is ["a","b","c"], then keep it.
// 3. if tokenUris/tokenData is ["a","","c"], then it's invalid.
// NOTE: Only use this before sending pkg.
func (nftpd *NonFungibleTokenPacketData) ShapeContent() error {

	shape, err := ValidateContent(nftpd.TokenUris)
	if err != nil {
		return sdkerrors.Wrap(err, "entries of TokenUris must be either all empty string or all non-empty string")
	}
	if shape {
		nftpd.TokenUris = nil
	}

	shape, err = ValidateContent(nftpd.TokenData)
	if err != nil {
		return sdkerrors.Wrap(err, "entries of TokenData must be either all empty string or all non-empty string")
	}
	if shape {
		nftpd.TokenData = nil
	}

	return nil
}

// GetBytes is a helper for serializing
func (nftpd NonFungibleTokenPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(MustProtoMarshalJSON(&nftpd))
}

// ValidateContent is used to validate contents of TokenUris/TokenData if TokenUris/TokenData
// has the invalid form of ["a","","c"]. It returns true if contents need to be shaped, false if
// needless, and error if invalid.
func ValidateContent(contents []string) (shape bool, err error) {
	if contents == nil {
		return false, nil
	}

	if len(contents) == 0 {
		return true, nil
	}

	emptyStringCount := 0
	for _, v := range contents {
		if v == "" {
			emptyStringCount++
		}
	}

	if emptyStringCount == len(contents) {
		return true, nil
	}

	if emptyStringCount == 0 {
		return false, nil
	}

	return false, ErrInvalidPacket
}

func GetIfExist(i int, data []string) string {
	if i < 0 || i >= len(data) {
		return ""
	}
	return data[i]
}
