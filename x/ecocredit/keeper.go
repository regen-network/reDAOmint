package ecocredit

import (
	"fmt"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/orm"
)

type Keeper struct {
	cdc               *codec.Codec
	storeKey          sdk.StoreKey
	creditClassBucket orm.AutoIDBucket
	creditBucket      orm.AutoIDBucket
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{cdc: cdc, storeKey: storeKey,
		creditClassBucket: orm.NewAutoIDBucket(storeKey, "credit-class", cdc, nil),
		creditBucket:orm.NewAutoIDBucket(storeKey, "credit", cdc, nil),
	}
}

func CreditClassFromBech32(bech string) (CreditClassID, error) {
	hrp, bz, err := bech32.Decode(bech)
	if err != nil {
		return nil, err
	}
	if hrp != "ecocls" {
		return nil, fmt.Errorf("not a credit class %s", bech)
	}
	return bz, err
}

func (k Keeper) CreateCreditClass(ctx sdk.Context, metadata CreditClassMetadata) (CreditClassID, error) {
	return k.creditClassBucket.Create(ctx, metadata)
}

func (k Keeper) IssueCredit(ctx sdk.Context, metadata CreditMetadata, holder sdk.AccAddress) (CreditID, error) {
	id, err := k.creditBucket.Create(ctx, metadata)
	if err != nil {
		return nil, err
	}
	return id, err
}

func (k Keeper) SendCredit(ctx sdk.Context, credit CreditID, from sdk.AccAddress, to sdk.AccAddress, units sdk.Dec) error {
	//store := ctx.KVStore(k.storeKey)
	panic("TODO")
}

func (k Keeper) BurnCredit(ctx sdk.Context, credit CreditID, holder sdk.AccAddress, units sdk.Dec) error {
	//store := ctx.KVStore(k.storeKey)
	panic("TODO")
}
