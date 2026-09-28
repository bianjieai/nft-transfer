package cli

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"

	"github.com/bianjieai/nft-transfer/types"
)

const (
	flagPacketTimeoutHeight    = "packet-timeout-height"
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	flagPacketMemo             = "packet-memo"
	flagAbsoluteTimeouts       = "absolute-timeouts"
)

// NewTransferTxCmd returns the command to create a NewMsgTransfer transaction
func NewTransferTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [src-port] [src-channel] [receiver] [classID] [tokenIDs]",
		Short: "Transfer a non-fungible token through IBC",
		Long: strings.TrimSpace(`Transfer a non-fungible token through IBC. Timeouts can be specified
as absolute or relative using the "absolute-timeouts" flag. Timeout height can be set by passing in the height string
in the form {revision}-{height} using the "packet-timeout-height" flag and "absolute-timeouts".
Relative timeouts support timestamps only and are added to the local clock time.
Any absolute timeout set to 0 is disabled.`),
		Example: fmt.Sprintf("%s tx nft-transfer transfer [src-port] [src-channel] [receiver] [classID] [tokenIDs]", version.AppName),
		Args:    cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sender := clientCtx.GetFromAddress().String()
			srcPort := args[0]
			srcChannel := args[1]
			receiver := args[2]
			classID := args[3]
			tokenIDs := strings.Split(args[4], ",")

			if len(tokenIDs) == 0 {
				return errors.New("tokenIDs cannot be empty")
			}

			timeoutHeightStr, err := cmd.Flags().GetString(flagPacketTimeoutHeight)
			if err != nil {
				return err
			}
			timeoutHeight, err := clienttypes.ParseHeight(timeoutHeightStr)
			if err != nil {
				return err
			}

			timeoutTimestamp, err := cmd.Flags().GetUint64(flagPacketTimeoutTimestamp)
			if err != nil {
				return err
			}

			absoluteTimeouts, err := cmd.Flags().GetBool(flagAbsoluteTimeouts)
			if err != nil {
				return err
			}

			memo, err := cmd.Flags().GetString(flagPacketMemo)
			if err != nil {
				return err
			}

			if !absoluteTimeouts {
				timeoutTimestamp, err = relativeTimeoutTimestamp(timeoutHeight, timeoutTimestamp, time.Now())
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgTransfer(
				srcPort, srcChannel, classID, tokenIDs, sender, receiver, timeoutHeight, timeoutTimestamp, memo,
			)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagPacketTimeoutHeight, types.DefaultRelativePacketTimeoutHeight, "Packet timeout block height. The timeout is disabled when set to 0-0.")
	cmd.Flags().Uint64(flagPacketTimeoutTimestamp, types.DefaultRelativePacketTimeoutTimestamp, "Packet timeout timestamp in nanoseconds from now. Default is 10 minutes. The timeout is disabled when set to 0.")
	cmd.Flags().String(flagPacketMemo, "", "Packet memo. Default is empty")
	cmd.Flags().Bool(flagAbsoluteTimeouts, false, "Timeout flags are used as absolute timeouts.")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// relativeTimeoutTimestamp follows IBC v10's timestamp-only relative timeouts.
func relativeTimeoutTimestamp(height clienttypes.Height, timestamp uint64, now time.Time) (uint64, error) {
	if !height.IsZero() {
		return 0, errors.New("relative timeouts using block height are not supported; use --absolute-timeouts")
	}
	if timestamp == 0 {
		return 0, errors.New("relative timeouts must provide a non-zero timestamp")
	}
	if now.UnixNano() <= 0 {
		return 0, errors.New("local clock time is not greater than Jan 1st, 1970 12:00 AM")
	}
	return uint64(now.UnixNano()) + timestamp, nil
}
