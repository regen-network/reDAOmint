package ecocredit

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateCreditClass{}, "ecocredit/MsgCreateCreditClass", nil)
	cdc.RegisterConcrete(MsgIssueCredit{}, "ecocredit/MsgIssueCredit", nil)
	cdc.RegisterConcrete(MsgSendCredit{}, "ecocredit/MsgSendCredit", nil)
	cdc.RegisterConcrete(MsgBurnCredit{}, "ecocredit/MsgBurnCredit", nil)
	cdc.RegisterConcrete(CreditClassMetadata{}, "ecocredit/CreditClassMetadata", nil)
	cdc.RegisterConcrete(CreditMetadata{}, "ecocredit/CreditMetadata", nil)
	cdc.RegisterConcrete(CreditHolding{}, "ecocredit/CreditHolding", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}
