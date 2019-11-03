package redaomint

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/gaia/orm"
	"github.com/cosmos/gaia/x/ecocredit"
)

type Keeper struct {
	cdc             *codec.Codec
	storeKey        sdk.StoreKey
	accountKeeper   auth.AccountKeeper
	bankKeeper      bank.Keeper
	supplyKeeper    supply.Keeper
	ecocreditKeeper ecocredit.Keeper
	ibcKeeper       ibc.Keeper
	metadataBucket  orm.AutoIDBucket
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, supplyKeeper supply.Keeper, ecocreditKeeper ecocredit.Keeper, ibcKeeper ibc.Keeper) Keeper {
	return Keeper{cdc: cdc, storeKey: storeKey, accountKeeper: accountKeeper, bankKeeper: bankKeeper, supplyKeeper: supplyKeeper, ecocreditKeeper: ecocreditKeeper, ibcKeeper: ibcKeeper}
}

func Denom(redaomint sdk.AccAddress) string {
	return fmt.Sprintf("redao:%x", redaomint)
}

func (k Keeper) CreateReDAOMint(ctx sdk.Context, metadata ReDAOMintMetadata) (addr sdk.AccAddress, denom string, err error) {
	addr, err = k.metadataBucket.Create(ctx, metadata)
	if err != nil {
		return nil, "", err
	}
	k.accountKeeper.SetAccount(ctx, &auth.BaseAccount{Address: addr})
	return addr, Denom(addr), err
}

func (k Keeper) MintShares(ctx sdk.Context, redaomint sdk.AccAddress, shares sdk.Int) error {
	coins := sdk.Coins{sdk.Coin{Denom: Denom(redaomint), Amount: shares}}
	err := k.supplyKeeper.MintCoins(ctx, ModuleName, coins)
	if err != nil {
		return err
	}
	_, err = k.bankKeeper.AddCoins(ctx, redaomint, coins)
	if err != nil {
		return err
	}
	return nil
}
