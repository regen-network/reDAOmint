package redaomint

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/gaia/orm"
	"github.com/cosmos/gaia/x/ecocredit"
)

type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey
	accountKeeper auth.AccountKeeper
	bankKeeper bank.Keeper
	ecocreditKeeper ecocredit.Keeper
	ibcKeeper ibc.Keeper
	metadataBucket orm.AutoIDBucket
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, ecocreditKeeper ecocredit.Keeper, ibcKeeper ibc.Keeper) Keeper {
	return Keeper{cdc: cdc, storeKey: storeKey, accountKeeper: accountKeeper, bankKeeper: bankKeeper, ecocreditKeeper: ecocreditKeeper, ibcKeeper: ibcKeeper}
}

func (k Keeper) CreateReDAOMint(ctx sdk.Context, metadata ReDAOMintMetadata) (addr sdk.AccAddress, denom string, err error) {
	addr, err = k.metadataBucket.Create(ctx, metadata)
	if err != nil {
		return nil, "", err
	}
	k.accountKeeper.SetAccount(ctx, &auth.BaseAccount{Address:addr})
	return addr, fmt.Sprintf("redao:%x", addr), err
}

func (k Keeper) ContributeReDAOMint(ctx sdk.Context, contributor sdk.AccAddress, redaomint sdk.AccAddress, funds sdk.Coins, priceInfo []byte) (sdk.Coins, sdk.Error) {
	err := k.bankKeeper.SendCoins(ctx, contributor, redaomint, funds)
	if err != nil {
		return sdk.Coins{}, err
	}
	// Use this to verify price info
	// k.ibcConnKeeper.VerifyMembership(ctx,)
	panic("TODO: mint redaomint shares")
}

