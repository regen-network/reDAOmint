/* Package ecocredit defines a fractional NFT (non-fungible token) module with metadata specifically targeted
for usage as ecosystem service credits. The model could be generalized for other fractional asset types in the future.
 */
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

const (
	IndexByGeoPolygon = "polygon"
)

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{cdc: cdc, storeKey: storeKey,
		creditClassBucket: orm.NewAutoIDBucket(storeKey, "credit-class", cdc, nil),
		creditBucket: orm.NewAutoIDBucket(storeKey, "credit", cdc, []orm.Index{
			{IndexByGeoPolygon, func(key []byte, value interface{}) (indexValue []byte, err error) {
				meta := value.(CreditMetadata)
				return meta.GeoPolygon, nil
			}},
		}),
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

// CreateCreditClass creates a new credit class with a set of authorized issuers
func (k Keeper) CreateCreditClass(ctx sdk.Context, metadata CreditClassMetadata) (CreditClassID, error) {
	return k.creditClassBucket.Create(ctx, metadata)
}

// CreditHolding describes the fractional holdings of a specific credit including units burned or in the language
// of carbon credits "retired", and liquid units that can still be transferred
type CreditHolding struct {
	Credit      CreditID       `json:"id"`
	Holder      sdk.AccAddress `json:"holder"`
	LiquidUnits sdk.Dec        `json:"liquid_units"`
	BurnedUnits sdk.Dec        `json:"burned_units"`
}

func (c CreditHolding) ID() []byte {
	return []byte(fmt.Sprintf("%x/%x", c.Credit, c.Holder))
}

// Issue credits issues some units of a credit class for a specific land area over a specific date range
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

// SendCredit sends fractional units of a credit from one account to another account
func (k Keeper) SendCredit(ctx sdk.Context, credit CreditID, from sdk.AccAddress, to sdk.AccAddress, units sdk.Dec) error {
	holding := CreditHolding{Credit: credit, Holder: from}
	err := k.creditHoldingsBucket.GetOne(ctx, &holding)
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
	holding2 := CreditHolding{Credit: credit, Holder: to}
	err = k.creditHoldingsBucket.GetOne(ctx, &holding2)
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

// BurnCredit burns some units of a credit that the holder holds. Burned units are still attached to the account that
// burned them for record of where they were ultimately "retired". In the language of carbon credits, retirement
// is used to take credits out of circulation which means that the holder retiring them is using them as an offset.
// So basically "burning" credits corresponds to the actual usage of ecosystem services.
func (k Keeper) BurnCredit(ctx sdk.Context, credit CreditID, holder sdk.AccAddress, units sdk.Dec) error {
	holding := CreditHolding{Credit: credit, Holder: holder}
	err := k.creditHoldingsBucket.GetOne(ctx, &holding)
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
	// TODO update credit metadata
	return nil
}

// GetCreditHolding gets the holdings of a specific credit by a specific holder
func (k Keeper) GetCreditHolding(ctx sdk.Context, credit CreditID, holder sdk.AccAddress) (holding CreditHolding, found bool) {
	holding = CreditHolding{Credit: credit, Holder: holder}
	err := k.creditHoldingsBucket.GetOne(ctx, &holding)
	if err != nil {
		return holding, false
	}
	return holding, true
}

// IterateCreditsByGeoPolygon iterators overall all credits for a specific geo-polygon. NOTE: this approach is not
// for use in production as it does not handle polygon overlaps. This method is used for demonstration purposes only
// until we have on-chain geo-index support or this iteration gets moved off-chain.
func (k Keeper) IterateCreditsByGeoPolygon(ctx sdk.Context, geoPolygon []byte, callback func(metadata CreditMetadata) (stop bool)) {
	iterator, err := k.creditBucket.ByIndex(ctx, IndexByGeoPolygon, geoPolygon)
	if err != nil {
		return
	}
	for {
		var metadata CreditMetadata
		_, err = iterator.LoadNext(&metadata)
		if err != nil {
			break
		}
		if callback(metadata) {
			return
		}
	}
}
