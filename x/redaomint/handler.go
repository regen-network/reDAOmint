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
			_, _, err :=k.CreateReDAOMint(ctx, msg.ReDAOMintMetadata)
			return sdk.ResultFromError(err)
		case MsgContributeReDAOMint:
			_, err := k.ContributeReDAOMint(ctx, msg.Sender, msg.ReDAOMint, msg.Funds, msg.PriceInfo)
			if err != nil {
				return err.Result()
			}
			return sdk.Result{}
		case MsgAllocateLandShares:
			return sdk.Result{}
		case MsgWithdrawCredits:
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
