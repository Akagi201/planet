package keeper

import (
	"context"

	"github.com/Akagi201/planet/x/blog/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
)

func (k msgServer) SendIbcUpdatePost(goCtx context.Context, msg *types.MsgSendIbcUpdatePost) (*types.MsgSendIbcUpdatePostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: logic before transmitting the packet

	// Construct the packet
	var packet types.IbcUpdatePostPacketData

	packet.PostID = msg.PostID
	packet.Title = msg.Title
	packet.Content = msg.Content

	// Transmit the packet
	_, err := k.TransmitIbcUpdatePostPacket(
		ctx,
		packet,
		msg.Port,
		msg.ChannelID,
		clienttypes.ZeroHeight(),
		msg.TimeoutTimestamp,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendIbcUpdatePostResponse{}, nil
}
