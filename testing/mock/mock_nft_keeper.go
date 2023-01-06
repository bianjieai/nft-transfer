package mock

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"

	nfttransfer "github.com/bianjieai/nft-transfer/types"
)

type (
	MockNFTKeeper struct {
		nk  nftkeeper.Keeper
		cdc codec.Codec
	}

	WrappedClass struct {
		nft.Class
		cdc ClassMetadataResolver
	}
	WrappedNFT struct {
		nft.NFT
		cdc TokenMetadataResolver
	}
)

func Wrap(cdc codec.Codec, nk nftkeeper.Keeper) MockNFTKeeper {
	return MockNFTKeeper{nk, cdc}
}

func (w MockNFTKeeper) CreateOrUpdateClass(ctx sdk.Context,
	classID,
	classURI,
	classData string,
) error {
	if !w.nk.HasClass(ctx, classID) {
		any, err := w.UnmarshalClassMetadata(classData)
		if err != nil {
			return err
		}
		return w.nk.SaveClass(ctx, nft.Class{
			Id:   classID,
			Uri:  classURI,
			Data: any,
		})
	}
	if len(classData) == 0 {
		return nil
	}

	class, _ := w.nk.GetClass(ctx, classID)
	class.Uri = classURI
	any, err := w.UnmarshalClassMetadata(classData)
	if err != nil {
		return err
	}
	class.Data = any
	return w.nk.UpdateClass(ctx, class)
}
func (w MockNFTKeeper) Mint(ctx sdk.Context,
	classID,
	tokenID,
	tokenURI,
	tokenData string,
	receiver sdk.AccAddress,
) error {
	any, err := w.UnmarshalTokenMetadata(tokenData)
	if err != nil {
		return err
	}
	nft := nft.NFT{
		ClassId: classID,
		Id:      tokenID,
		Uri:     tokenURI,
		Data:    any,
	}
	return w.nk.Mint(ctx, nft, receiver)
}
func (w MockNFTKeeper) Transfer(ctx sdk.Context,
	classID,
	tokenID string,
	tokenData string,
	receiver sdk.AccAddress,
) error {
	if err := w.nk.Transfer(ctx, classID, tokenID, receiver); err != nil {
		return err
	}
	if len(tokenData) == 0 {
		return nil
	}

	any, err := w.UnmarshalTokenMetadata(tokenData)
	if err != nil {
		return err
	}
	nft, _ := w.nk.GetNFT(ctx, classID, tokenID)
	nft.Data = any
	return w.nk.Update(ctx, nft)
}

func (w MockNFTKeeper) Burn(ctx sdk.Context, classID string, tokenID string) error {
	return w.nk.Burn(ctx, classID, tokenID)
}
func (w MockNFTKeeper) GetOwner(ctx sdk.Context, classID string, tokenID string) sdk.AccAddress {
	return w.nk.GetOwner(ctx, classID, tokenID)
}
func (w MockNFTKeeper) HasClass(ctx sdk.Context, classID string) bool {
	return w.nk.HasClass(ctx, classID)
}
func (w MockNFTKeeper) GetClass(ctx sdk.Context, classID string) (nfttransfer.Class, bool) {
	class, exist := w.nk.GetClass(ctx, classID)
	if !exist {
		return nil, exist
	}
	return WrappedClass{class, w}, true
}
func (w MockNFTKeeper) GetNFT(ctx sdk.Context, classID, tokenID string) (nfttransfer.NFT, bool) {
	nft, exist := w.nk.GetNFT(ctx, classID, tokenID)
	if !exist {
		return nil, exist
	}
	return WrappedNFT{nft, w}, true
}

func (wc WrappedClass) GetData() string {
	data, _ := wc.cdc.MarshalClassMetadata(wc.Data)
	return data
}

func (wc WrappedClass) GetID() string {
	return wc.Id
}

func (wc WrappedClass) GetURI() string {
	return wc.Uri
}

func (wc WrappedNFT) GetData() string {
	data, _ := wc.cdc.MarshalTokenMetadata(wc.Data)
	return data
}

func (wc WrappedNFT) GetID() string {
	return wc.Id
}

func (wc WrappedNFT) GetClassID() string {
	return wc.ClassId
}

func (wc WrappedNFT) GetURI() string {
	return wc.Uri
}
