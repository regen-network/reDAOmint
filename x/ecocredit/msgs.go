package ecocredit

import (
	"time"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CreditClassMetadata struct {
	// Designer is the entity which designs a credit class at the top-level and
	// certifies issuers
	Designer sdk.AccAddress
	// Name is the name the Designer gives to the credit, internally credits
	// are identified by their CreditClassID
	Name string
	// Issuers are those entities authorized to issue credits via MsgIssueCredit
	Issuers []sdk.AccAddress
}

// MsgCreateCreditClass creates a class of credits and returns a new CreditClassID
type MsgCreateCreditClass struct {
	CreditClassMetadata `json:"metadata"`
}

type CreditClassID []byte

type CreditMetadata struct {
	Issuer      sdk.AccAddress `json:"issuer"`
	CreditClass CreditClassID  `json:"credit_class"`
	GeoPolygon  []byte         `json:"geo_polygon"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	// LiquidUnits specifies how many tradeable units of this credit are issued for this polygon
	LiquidUnits sdk.Dec `json:"liquid_units"`
	BurnedUnits sdk.Dec `json:"burned_units"`
}

// MsgIssueCredit issues a credit to the Holder with the number of LiquidUnits provided
// for the provided credit class, polygon, and start and end dates. A new CreditID
// is returned. It is illegal to issue a credit where the provided polygon and dates
// overlaps with those of an existing credit of the same class
type MsgIssueCredit struct {
	CreditMetadata `json:"metadata"`
	// Holder receives the credit from the issuer and can send it to other holders
	// or consume it
	Holder sdk.AccAddress `json:"holder"`
}

type CreditID []byte

// MsgSendCredit sends the provided number of units of the credit from the from
// address to the to address
type MsgSendCredit struct {
	Credit CreditID
	From   sdk.AccAddress
	To     sdk.AccAddress
	Units  sdk.Dec
}

// MsgBurnCredit consumes the provided number of units of the credit, essentially
// burning or retiring those units. This operation is used to actually use
// the credit as an offset. Otherwise, the holder of the credit is simply
// holding the credit as an asset but has not claimed the offset. Once a
// credit has been consumed, it can no longer be transferred
type MsgBurnCredit struct {
	Credit CreditID
	Holder sdk.AccAddress
	Units  sdk.Dec
}

func (m MsgCreateCreditClass) Route() string {
	return "ecocredit"
}

func (m MsgCreateCreditClass) Type() string {
	return "create-credit-class"
}

func (m MsgCreateCreditClass) ValidateBasic() sdk.Error {
	return nil
}

func (m MsgCreateCreditClass) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgCreateCreditClass) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Designer}
}

func (m MsgIssueCredit) Route() string {
	return "ecocredit"
}

func (m MsgIssueCredit) Type() string {
	return "issue-credit"
}

func (m MsgIssueCredit) ValidateBasic() sdk.Error {
	return nil
}

func (m MsgIssueCredit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgIssueCredit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Issuer}
}

func (m MsgSendCredit) Route() string {
	return "ecocredit"
}

func (m MsgSendCredit) Type() string {
	return "send-credit"
}

func (m MsgSendCredit) ValidateBasic() sdk.Error {
	return nil
}

func (m MsgSendCredit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgSendCredit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.From}
}

func (m MsgBurnCredit) Route() string {
	return "ecocredit"
}

func (m MsgBurnCredit) Type() string {
	return "burn-credit"
}

func (m MsgBurnCredit) ValidateBasic() sdk.Error {
	return nil
}

func (m MsgBurnCredit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgBurnCredit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Holder}
}
