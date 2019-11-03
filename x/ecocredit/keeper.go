package ecocredit

import (
	"fmt"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/orm"
)

type Keeper struct {
	cdc                  *codec.Codec
	storeKey             sdk.StoreKey
	creditClassBucket    orm.AutoIDBucket
	creditBucket         orm.AutoIDBucket
	creditHoldingsBucket orm.NaturalKeyBucket
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{cdc: cdc, storeKey: storeKey,
		creditClassBucket:    orm.NewAutoIDBucket(storeKey, "credit-class", cdc, nil),
		creditBucket:         orm.NewAutoIDBucket(storeKey, "credit", cdc, nil),
		creditHoldingsBucket: orm.NewNaturalKeyBucket(storeKey, "credit-holdings", cdc, nil),
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

type CreditHolding struct {
	Credit      CreditID       `json:"id"`
	Holder      sdk.AccAddress `json:"holder"`
	LiquidUnits sdk.Dec        `json:"liquid_units"`
	BurnedUnits sdk.Dec        `json:"burned_units"`
}

func (c CreditHolding) ID() []byte {
	return []byte(fmt.Sprintf("%x/%x", c.Credit, c.Holder))
}

func (k Keeper) IssueCredit(ctx sdk.Context, metadata CreditMetadata, holder sdk.AccAddress) (CreditID, error) {
	id, err := k.creditBucket.Create(ctx, metadata)
	if err != nil {
		return nil, err
	}
	err = k.creditHoldingsBucket.Save(ctx, CreditHolding{Credit: id, Holder: holder, LiquidUnits: metadata.LiquidUnits, BurnedUnits: metadata.BurnedUnits})
	if err != nil {
		return nil, err
	}
	return id, err
}

func (k Keeper) SendCredit(ctx sdk.Context, credit CreditID, from sdk.AccAddress, to sdk.AccAddress, units sdk.Dec) error {
	var holding CreditHolding
	err := k.creditHoldingsBucket.One(ctx, CreditHolding{Credit: credit, Holder: from}.ID(), &holding)
	if err != nil {
		return err
	}
	holding.LiquidUnits = holding.LiquidUnits.Sub(units)
	if holding.LiquidUnits.IsNegative() {
		return fmt.Errorf("not enough units")
	}
	err = k.creditHoldingsBucket.Save(ctx, holding)
	if err != nil {
		return err
	}
	var holding2 CreditHolding
	err = k.creditHoldingsBucket.One(ctx, CreditHolding{Credit: credit, Holder: to}.ID(), &holding2)
	if err != nil {
		err = k.creditHoldingsBucket.Save(ctx, CreditHolding{Credit: credit, Holder: to, LiquidUnits: units})
		if err != nil {
			return err
		}
	} else {
		holding2.LiquidUnits = holding2.LiquidUnits.Add(units)
		err = k.creditHoldingsBucket.Save(ctx, holding2)
	}
	return nil
}

func (k Keeper) BurnCredit(ctx sdk.Context, credit CreditID, holder sdk.AccAddress, units sdk.Dec) error {
	var holding CreditHolding
	err := k.creditHoldingsBucket.One(ctx, CreditHolding{Credit: credit, Holder: holder}.ID(), &holding)
	if err != nil {
		return err
	}
	holding.LiquidUnits = holding.LiquidUnits.Sub(units)
	if holding.LiquidUnits.IsNegative() {
		return fmt.Errorf("not enough units")
	}
	holding.BurnedUnits = holding.BurnedUnits.Add(units)
	err = k.creditHoldingsBucket.Save(ctx, holding)
	if err != nil {
		return err
	}
	return nil
}
