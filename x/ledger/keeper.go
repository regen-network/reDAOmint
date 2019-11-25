package ledger

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/orm"
	"time"
)

const (
	IndexByAsset = "by-asset"
)

type keeper struct {
	key                      sdk.StoreKey
	cdc                      *codec.Codec
	addressCompressionBucket orm.ExternalKeyBucket
	assetCompressionBucket   orm.ExternalKeyBucket
	holdingsBucket           orm.ExternalKeyBucket
	supplyBucket             orm.ExternalKeyBucket
}

type Keeper interface {
	GetHolding(ctx sdk.Context, acc sdk.Address, asset AssetID) HoldingView
	Transfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, asset AssetID, amount Dec) error
	Mint(ctx sdk.Context, acc sdk.AccAddress, asset AssetID, amount Dec) error
	Burn(ctx sdk.Context, acc sdk.AccAddress, asset AssetID, amount Dec) error
}

type Authority string
type Module string

type AssetMetadata struct {
	Asset     AssetID
	Authority Authority
	Module    Module
	Name      string
}

type AssetManager interface {
	GetMetadata() AssetMetadata
	GetBalance(ctx sdk.Context, acc sdk.Address) Dec
	Transfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amount Dec) error
}

type NativeAssetManager interface {
	AssetManager
	GetSupply() Dec
	Burn(ctx sdk.Context, acc sdk.AccAddress, amount Dec) error
}

type AssetsManager interface {
	GetMetadata(asset AssetID) AssetMetadata
	GetBalance(ctx sdk.Context, acc sdk.Address, asset AssetID) Dec
	Transfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, asset AssetID, amount Dec) error
	Burn(ctx sdk.Context, acc sdk.AccAddress, asset AssetID, amount Dec) error
}

func NewKeeper(key sdk.StoreKey, cdc *codec.Codec) Keeper {
	return keeper{
		key: key,
		cdc: cdc,
		holdingsBucket: orm.NewExternalKeyBucket(
			key,
			"holdings",
			cdc,
			[]orm.Index{
				{Name: IndexByAsset,
					Indexer: func(key []byte, value interface{}) (indexValue []byte, err error) {
						return []byte(value.(Holding).GetAsset()), nil

					}},
			},
		),
	}
}

type AssetID []byte

type Dec = sdk.Dec

type HoldingView interface {
	GetAddress() sdk.AccAddress
	GetAsset() AssetID
	GetBalance() Dec
	SpendableBalance(time time.Time) Dec
}

type Holding interface {
	HoldingView
	Transfer(amount Dec) (Dec, error)
	Burn(amount Dec) (Dec, error)
}

type BaseHolding struct {
	Address sdk.AccAddress `json:"address"`
	Asset   AssetID        `json:"asset"`
	Balance Dec            `json:"balance"`
}

func (b BaseHolding) GetAddress() sdk.AccAddress {
	return b.Address
}

func (b BaseHolding) GetAsset() AssetID {
	return b.Asset
}

func (b BaseHolding) GetBalance() Dec {
	return b.Balance
}

func (b BaseHolding) SpendableBalance(time time.Time) Dec {
	return b.Balance
}

func (b BaseHolding) Transfer(amount Dec) (Dec, error) {
	newBalance := b.Balance.Add(amount)
	if newBalance.IsNegative() {
		return b.Balance, fmt.Errorf("insufficient balance")
	}
	b.Balance = newBalance
	return newBalance, nil
}

func (b BaseHolding) Burn(amount Dec) (Dec, error) {
	return b.Transfer(amount.Neg())
}

func (k keeper) GetHolding(ctx sdk.Context, from sdk.Address, asset AssetID) HoldingView {
	panic("TODO")
}

func (k keeper) HoldingKey(ctx sdk.Context, address sdk.AccAddress, asset AssetID) []byte {
	// TODO: get asset key from assetCompressionBucket
	panic("TODO")
}

func (k keeper) NewHolding(ctx sdk.Context, address sdk.AccAddress, asset AssetID) (Holding, error) {
	return BaseHolding{}, nil
}

func (k keeper) Transfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, asset AssetID, amount Dec) error {
	var holding Holding
	err := k.holdingsBucket.GetOne(ctx, k.HoldingKey(ctx, from, asset), &holding)
	if err != nil {
		return err
	}

	_, err = holding.Transfer(amount.Neg())
	if err != nil {
		return err
	}

	var holding2 Holding
	err = k.holdingsBucket.GetOne(ctx, k.HoldingKey(ctx, to, asset), &holding2)
	if err != nil {
		holding2, err = k.NewHolding(ctx, to, asset)
		if err != nil {
			return err
		}
	}

	_, err = holding.Transfer(amount)
	if err != nil {
		return err
	}

	return nil
}

func (k keeper) Mint(ctx sdk.Context, acc sdk.AccAddress, asset AssetID, amount Dec) error {
	var holding Holding
	err := k.holdingsBucket.GetOne(ctx, k.HoldingKey(ctx, acc, asset), &holding)
	if err != nil {
		holding, err = k.NewHolding(ctx, acc, asset)
		if err != nil {
			return err
		}
	}

	_, err = holding.Transfer(amount)
	if err != nil {
		return err
	}

	return nil
}

func (k keeper) Burn(ctx sdk.Context, acc sdk.AccAddress, asset AssetID, amount Dec) error {
	var holding Holding
	err := k.holdingsBucket.GetOne(ctx, k.HoldingKey(ctx, acc, asset), &holding)
	if err != nil {
		return err
	}

	_, err = holding.Burn(amount)
	if err != nil {
		return err
	}

	return nil
}
