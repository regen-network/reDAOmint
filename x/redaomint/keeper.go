package redaomint

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
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

func (k Keeper) WithdrawCredits(ctx sdk.Context, redaomint sdk.AccAddress, shareholder sdk.AccAddress) error {
	// NOTE: this is a hacky and incorrect way of determining a share at the time of withdraw,
	// but in order to have something more correct, the bank module would need to track the holders
	// of a single coin and possible allow access to historical balances
	denom := Denom(redaomint)
	coins := k.bankKeeper.GetCoins(ctx, shareholder).AmountOf(denom)
	totalCoins := k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)
	var share sdk.Dec
	share.Div(coins.BigInt(), totalCoins.BigInt())
	return nil
}

func (k Keeper) DistributeCredit(ctx sdk.Context, redaomint sdk.AccAddress, credit ecocredit.CreditID) error {
	holding, found := k.ecocreditKeeper.GetCreditHolding(ctx, credit, redaomint)
	if !found {
		return fmt.Errorf("not found")
	}
	denom := Denom(redaomint)
	totalCoins := k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)
	k.accountKeeper.IterateAccounts(ctx, func(account exported.Account) (stop bool) {
		coins := account.GetCoins().AmountOf(denom)
		var share sdk.Dec
		share.Div(coins.BigInt(), totalCoins.BigInt())
		units := holding.LiquidUnits.Mul(share)
		if units.IsPositive() {
			err := k.ecocreditKeeper.SendCredit(ctx, credit, redaomint, account.GetAddress(), units)
			if err != nil {
				return err
			}
		}
		return false
	})
	return nil
}
