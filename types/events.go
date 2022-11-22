package types

// IBC transfer events
const (
	EventTypeTimeout      = "timeout"
	EventTypePacket       = "non_fungible_token_packet"
	EventTypeTransfer     = "ibc_nft_transfer"
	EventTypeChannelClose = "channel_closed"
	EventTypeClassTrace   = "class_trace"

	AttributeKeySender     = "sender"
	AttributeKeyReceiver   = "receiver"
	AttributeKeyClassID    = "classID"
	AttributeKeyTokenIDs   = "tokenIDs"
	AttributeKeyAckSuccess = "success"
	AttributeKeyAckError   = "error"
	AttributeKeyTraceHash  = "trace_hash"
)
