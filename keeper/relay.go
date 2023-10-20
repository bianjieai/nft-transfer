package keeper

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	coretypes "github.com/cosmos/ibc-go/v7/modules/core/types"

	"github.com/bianjieai/nft-transfer/types"
)

// SendTransfer handles nft-transfer sending logic.
// A sending chain may be acting as a source or sink zone.
//
// when a chain is sending tokens across a port and channel which are
// not equal to the last prefixed port and channel pair, it is acting as a source zone.
// when tokens are sent from a source zone, the destination port and
// channel will be prefixed onto the classId (once the tokens are received)
// adding another hop to the tokens record.
//
// when a chain is sending tokens across a port and channel which are
// equal to the last prefixed port and channel pair, it is acting as a sink zone.
// when tokens are sent from a sink zone, the last prefixed port and channel
// pair on the classId is removed (once the tokens are received), undoing the last hop in the tokens record.
//
// For example, assume these steps of transfer occur:
// A -> B -> C -> A -> C -> B -> A
//
// |                    sender  chain                      |                       receiver     chain              |
// | :-----: | -------------------------: | :------------: | :------------: | -------------------------: | :-----: |
// |  chain  |                    classID | (port,channel) | (port,channel) |                    classID |  chain  |
// |    A    |                   nftClass |    (p1,c1)     |    (p2,c2)     |             p2/c2/nftClass |    B    |
// |    B    |             p2/c2/nftClass |    (p3,c3)     |    (p4,c4)     |       p4/c4/p2/c2/nftClass |    C    |
// |    C    |       p4/c4/p2/c2/nftClass |    (p5,c5)     |    (p6,c6)     | p6/c6/p4/c4/p2/c2/nftClass |    A    |
// |    A    | p6/c6/p4/c4/p2/c2/nftClass |    (p6,c6)     |    (p5,c5)     |       p4/c4/p2/c2/nftClass |    C    |
// |    C    |       p4/c4/p2/c2/nftClass |    (p4,c4)     |    (p3,c3)     |             p2/c2/nftClass |    B    |
// |    B    |             p2/c2/nftClass |    (p2,c2)     |    (p1,c1)     |                   nftClass |    A    |
func (k Keeper) SendTransfer(
	ctx sdk.Context,
	sourcePort,
	sourceChannel,
	classID string,
	tokenIDs []string,
	sender sdk.AccAddress,
	receiver string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	memo string,
) (uint64, error) {
	if !k.GetSendEnabled(ctx) {
		return 0, types.ErrSendDisabled
	}

	channel, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return 0, errorsmod.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}

	destinationPort := channel.GetCounterparty().GetPortID()
	destinationChannel := channel.GetCounterparty().GetChannelID()

	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, errorsmod.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	// See spec for this logic: https://github.com/cosmos/ibc/blob/master/spec/app/ics-721-nft-transfer/README.md#packet-relay
	packet, err := k.createOutgoingPacket(ctx,
		sourcePort,
		sourceChannel,
		classID,
		tokenIDs,
		sender,
		receiver,
		memo,
	)
	if err != nil {
		return 0, err
	}

	sequence, err := k.ics4Wrapper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packet.GetBytes())
	if err != nil {
		return 0, err
	}

	defer func() {
		labels := []metrics.Label{
			telemetry.NewLabel(coretypes.LabelDestinationPort, destinationPort),
			telemetry.NewLabel(coretypes.LabelDestinationChannel, destinationChannel),
		}

		telemetry.SetGaugeWithLabels(
			[]string{"tx", "msg", "ibc", "nft-transfer"},
			float32(len(tokenIDs)),
			[]metrics.Label{telemetry.NewLabel("class_id", classID)},
		)

		telemetry.IncrCounterWithLabels(
			[]string{"ibc", types.ModuleName, "send"},
			1,
			labels,
		)
	}()
	return sequence, nil
}

// OnRecvPacket processes a cross chain fungible token transfer. If the
// sender chain is the source of minted tokens then vouchers will be minted
// and sent to the receiving address. Otherwise if the sender chain is sending
// back tokens this chain originally transferred to it, the tokens are
// unescrowed and sent to the receiving address.
func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet,
	data types.NonFungibleTokenPacketData) error {
	if !k.GetReceiveEnabled(ctx) {
		return types.ErrReceiveDisabled
	}

	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return err
	}

	// See spec for this logic: https://github.com/cosmos/ibc/blob/master/spec/app/ics-721-nft-transfer/README.md#packet-relay
	return k.processReceivedPacket(ctx, packet, data)
}

// OnAcknowledgementPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain. If the acknowledgement
// was a success then nothing occurs. If the acknowledgement failed, then
// the sender is refunded their tokens using the refundPacketToken function.
func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, data types.NonFungibleTokenPacketData, ack channeltypes.Acknowledgement) error {
	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		return k.refundPacketToken(ctx, packet, data)
	default:
		// the acknowledgement succeeded on the receiving chain so nothing
		// needs to be executed and no error needs to be returned
		return nil
	}
}

// OnTimeoutPacket refunds the sender since the original packet sent was
// never received and has been timed out.
func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data types.NonFungibleTokenPacketData) error {
	return k.refundPacketToken(ctx, packet, data)
}

// refundPacketToken will unescrow and send back the tokens back to sender
// if the sending chain was the source chain. Otherwise, the sent tokens
// were burnt in the original send so new tokens are minted and sent to
// the sending address.
func (k Keeper) refundPacketToken(ctx sdk.Context, packet channeltypes.Packet, data types.NonFungibleTokenPacketData) error {
	sender, err := sdk.AccAddressFromBech32(data.Sender)
	if err != nil {
		return err
	}

	voucherClassID, err := k.GetVoucherClassID(ctx, data.ClassId)
	if err != nil {
		return err
	}
	if types.IsAwayFromOrigin(packet.GetSourcePort(), packet.GetSourceChannel(), data.ClassId) {
		for i, tokenID := range data.TokenIds {
			if err := k.nftKeeper.Transfer(ctx, voucherClassID, tokenID, types.GetIfExist(i, data.TokenData), sender); err != nil {
				return err
			}
		}
		return nil
	}

	for i, tokenID := range data.TokenIds {
		if err := k.nftKeeper.Mint(ctx,
			voucherClassID,
			tokenID,
			types.GetIfExist(i, data.TokenUris),
			types.GetIfExist(i, data.TokenData),
			sender,
		); err != nil {
			return err
		}
	}
	return nil
}

// createOutgoingPacket will escrow the tokens to escrow account
// if the token was away from origin chain . Otherwise, the sent tokens
// were burnt in the sending chain and will unescrow the token to receiver
// in the destination chain
func (k Keeper) createOutgoingPacket(ctx sdk.Context,
	sourcePort,
	sourceChannel,
	classID string,
	tokenIDs []string,
	sender sdk.AccAddress,
	receiver string,
	memo string,
) (types.NonFungibleTokenPacketData, error) {
	class, exist := k.nftKeeper.GetClass(ctx, classID)
	if !exist {
		return types.NonFungibleTokenPacketData{}, errorsmod.Wrap(types.ErrInvalidClassID, "classId not exist")
	}

	var (
		// NOTE: class and hex hash correctness checked during msg.ValidateBasic
		fullClassPath = classID
		err           error
		tokenURIs     = make([]string, len(tokenIDs))
		tokenData     = make([]string, len(tokenIDs))
	)

	// deconstruct the token denomination into the denomination trace info
	// to determine if the sender is the source chain
	if strings.HasPrefix(classID, "ibc/") {
		fullClassPath, err = k.ClassPathFromHash(ctx, classID)
		if err != nil {
			return types.NonFungibleTokenPacketData{}, err
		}
	}

	isAwayFromOrigin := types.IsAwayFromOrigin(sourcePort,
		sourceChannel, fullClassPath)
	for i, tokenID := range tokenIDs {
		nft, exist := k.nftKeeper.GetNFT(ctx, classID, tokenID)
		if !exist {
			return types.NonFungibleTokenPacketData{}, errorsmod.Wrap(types.ErrInvalidTokenID, "tokenId not exist")
		}

		owner := k.nftKeeper.GetOwner(ctx, classID, tokenID)
		if !sender.Equals(owner) {
			return types.NonFungibleTokenPacketData{}, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "not token owner")
		}

		tokenURIs[i] = nft.GetURI()
		tokenData[i] = nft.GetData()

		if isAwayFromOrigin {
			// create the escrow address for the tokens
			escrowAddress := types.GetEscrowAddress(sourcePort, sourceChannel)
			if err := k.nftKeeper.Transfer(ctx, classID, tokenID, nft.GetData(), escrowAddress); err != nil {
				return types.NonFungibleTokenPacketData{}, err
			}
		} else {
			if err := k.nftKeeper.Burn(ctx, classID, tokenID); err != nil {
				return types.NonFungibleTokenPacketData{}, err
			}
		}
	}

	packetData := types.NewNonFungibleTokenPacketData(
		fullClassPath,
		class.GetURI(),
		class.GetData(),
		tokenIDs,
		tokenURIs,
		sender.String(),
		receiver,
		tokenData,
		memo,
	)
	return packetData, packetData.ValidateBasic()
}

// processReceivedPacket will mint the tokens to receiver account
// if the token was away from origin chain . Otherwise, the sent tokens
// were burnt in the sending chain and will unescrow the token to receiver
// in the destination chain
func (k Keeper) processReceivedPacket(ctx sdk.Context, packet channeltypes.Packet,
	data types.NonFungibleTokenPacketData) error {
	receiver, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return err
	}

	if types.IsAwayFromOrigin(packet.GetSourcePort(), packet.GetSourceChannel(), data.ClassId) {
		// since SendPacket did not prefix the classID, we must prefix classID here
		classPrefix := types.GetClassPrefix(packet.GetDestPort(), packet.GetDestChannel())
		// NOTE: sourcePrefix contains the trailing "/"
		prefixedClassID := classPrefix + data.ClassId

		// construct the class trace from the full raw classID
		classTrace := types.ParseClassTrace(prefixedClassID)
		if !k.HasClassTrace(ctx, classTrace.Hash()) {
			k.SetClassTrace(ctx, classTrace)
		}

		voucherClassID := classTrace.IBCClassID()
		if err := k.nftKeeper.CreateOrUpdateClass(ctx,
			voucherClassID, data.ClassUri, data.ClassData); err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeClassTrace,
				sdk.NewAttribute(types.AttributeKeyTraceHash, classTrace.Hash().String()),
				sdk.NewAttribute(types.AttributeKeyClassID, voucherClassID),
			),
		)
		for i, tokenID := range data.TokenIds {
			if err := k.nftKeeper.Mint(ctx,
				voucherClassID,
				tokenID,
				types.GetIfExist(i, data.TokenUris),
				types.GetIfExist(i, data.TokenData),
				receiver,
			); err != nil {
				return err
			}
		}
		return nil
	}

	// If the token moves in the direction of back to origin,
	// we need to unescrow the token and transfer it to the receiver

	// we should remove the prefix. For example:
	// p6/c6/p4/c4/p2/c2/nftClass -> p4/c4/p2/c2/nftClass
	unprefixedClassID, err := types.RemoveClassPrefix(packet.GetSourcePort(),
		packet.GetSourceChannel(), data.ClassId)
	if err != nil {
		return err
	}

	voucherClassID, err := k.GetVoucherClassID(ctx, unprefixedClassID)
	if err != nil {
		return err
	}

	escrowAddress := types.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
	for i, tokenID := range data.TokenIds {
		//NOTE: It must be verified here whether the nft is escrowed by the <destPort, destChannel> account
		//FIX https://github.com/game-of-nfts/gon-evidence/issues/346
		owner := k.nftKeeper.GetOwner(ctx, voucherClassID, tokenID)
		if !escrowAddress.Equals(owner) {
			return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "not token owner")
		}

		if err := k.nftKeeper.Transfer(ctx,
			voucherClassID, tokenID, types.GetIfExist(i, data.TokenData), receiver); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) GetVoucherClassID(ctx sdk.Context, classID string) (string, error) {

	// If "/" is not included after removing the prefix,
	// it means that nft has returned to the initial chain, and the classID after removing the prefix is the real classID
	if !strings.Contains(classID, "/") {
		return classID, nil
	}

	// If "/" is included after removing the prefix, there are two situations:
	//	1. The original classID itself contains "/",
	//	2. The current nft returns to the relay chain, not the original chain

	// First deal with case 1, if the classID can be found, return the result
	if k.nftKeeper.HasClass(ctx, classID) {
		return classID, nil
	}

	// If not found, generate classID according to classTrace
	return types.ParseClassTrace(classID).IBCClassID(), nil
}
