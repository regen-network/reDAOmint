package redaomint

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/gaia/x/ecocredit"
)

type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey
	accountKeeper auth.AccountKeeper
	ecocreditKeeper ecocredit.Keeper
}

