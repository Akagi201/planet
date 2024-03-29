package keeper

import (
	"errors"
	"strconv"

	"github.com/Akagi201/planet/x/blog/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
)

// TransmitIbcUpdatePostPacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitIbcUpdatePostPacket(
	ctx sdk.Context,
	packetData types.IbcUpdatePostPacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) (uint64, error) {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packetBytes, err := packetData.GetBytes()
	if err != nil {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: %w", err)
	}

	return k.channelKeeper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packetBytes)
}

// OnRecvIbcUpdatePostPacket processes packet reception
func (k Keeper) OnRecvIbcUpdatePostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcUpdatePostPacketData) (packetAck types.IbcUpdatePostPacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// TODO: packet reception logic
	pId, _ := strconv.Atoi(data.PostID)
	k.SetPost(ctx, types.Post{
		Id:      uint64(pId),
		Title:   data.Title,
		Content: data.Content,
	})
	packetAck.Ok = true
	return packetAck, nil
}

// OnAcknowledgementIbcUpdatePostPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementIbcUpdatePostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcUpdatePostPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:

		// TODO: failed acknowledgement logic
		_ = dispatchedAck.Error

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.IbcUpdatePostPacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// TODO: successful acknowledgement logic
		if !packetAck.Ok {
			return errors.New("fail to update")
		}
		k.SetSentPost(ctx,
			types.SentPost{
				PostID: data.PostID,
				Title:  data.Title,
			})

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutIbcUpdatePostPacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutIbcUpdatePostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcUpdatePostPacketData) error {

	// TODO: packet timeout logic
	k.AppendTimedoutPost(
		ctx,
		types.TimedoutPost{
			Creator: "", // TODO: add postID field
			Title:   data.Title,
			Chain:   packet.DestinationPort + "-" + packet.DestinationChannel,
		},
	)

	return nil
}
