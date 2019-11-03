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
			_, _, err := k.CreateReDAOMint(ctx, msg.ReDAOMintMetadata, msg.Founder, msg.FounderShares)
			return sdk.ResultFromError(err)
		case MsgMintShares:
			err := k.MintShares(ctx, msg.ReDAOMint, msg.Shares)
			return sdk.ResultFromError(err)
		case MsgAllocateLandShares:
			err := k.SetLandAllocation(ctx, msg.LandAllocation)
			return sdk.ResultFromError(err)
		case MsgPropose:
			_, err := k.CreateProposal(ctx, msg.Proposal)
			return sdk.ResultFromError(err)
		case MsgVote:
			err := k.Vote(ctx, msg.ProposalID, msg.Voter, msg.Vote)
			return sdk.ResultFromError(err)
		case MsgExecProposal:
			return k.ExecProposal(ctx, msg.ProposalID)
		default:
			errMsg := fmt.Sprintf("Unrecognized data Msg type: %s", ModuleName)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
