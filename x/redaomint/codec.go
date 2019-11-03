package redaomint

import (
"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateReDAOMint{}, "redaomint/MsgCreateReDAOMint", nil)
	cdc.RegisterConcrete(MsgMintShares{}, "redaomint/MsgMintShares", nil)
	cdc.RegisterConcrete(MsgAllocateLandShares{}, "redaomint/MsgAllocateLandShares", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}
