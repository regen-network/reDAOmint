package redaomint

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateReDAOMint:
			_, _, err :=k.CreateReDAOMint(ctx, msg.ReDAOMintMetadata, msg.Founder, msg.FounderShares)
			return sdk.ResultFromError(err)
		case MsgMintShares:
			err := k.MintShares(ctx, msg.ReDAOMint, msg.Shares)
			return sdk.ResultFromError(err)
		case MsgAllocateLandShares:
			return sdk.Result{}
		case MsgPropose:
			return sdk.Result{}
		case MsgVote:
			return sdk.Result{}
		case MsgExecProposal:
			return sdk.Result{}
		default:
			errMsg := fmt.Sprintf("Unrecognized data Msg type: %s", ModuleName)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
