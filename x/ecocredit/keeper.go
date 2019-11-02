package ecocredit

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey
}

func (k Keeper) CreateCreditClass(ctx sdk.Context, designer sdk.AccAddress, name string, issuers []sdk.AccAddress) (CreditClassID, error) {
	store := ctx.KVStore(k.storeKey)
	panic("TODO")
}

func (k Keeper) IssueCredit(ctx sdk.Context, metadata CreditMetadata, holder sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	panic("TODO")
}

func (k Keeper) SendCredit(ctx sdk.Context, credit CreditID, from sdk.AccAddress, to sdk.AccAddress, units sdk.Dec) error {
	store := ctx.KVStore(k.storeKey)
	panic("TODO")
}

func (k Keeper) BurnCredit(ctx sdk.Context, credit CreditID, holder sdk.AccAddress, units sdk.Dec) error {
	store := ctx.KVStore(k.storeKey)
	panic("TODO")
}
