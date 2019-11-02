package ecocredit

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateCreditClass:
			_, err := k.CreateCreditClass(ctx, msg.CreditClassMetadata)
			return sdk.ResultFromError(err)
		case MsgIssueCredit:
			_, err := k.IssueCredit(ctx, msg.CreditMetadata, msg.Holder)
			return sdk.ResultFromError(err)
		case MsgSendCredit:
			return sdk.Result{}
		case MsgBurnCredit:
			return sdk.Result{}
		default:
			errMsg := fmt.Sprintf("Unrecognized data Msg type: %s", ModuleName)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
