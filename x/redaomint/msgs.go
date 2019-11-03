package redaomint

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/ecocredit"
)

type ReDAOMintMetadata struct {
	Description           string                    `json:"description"`
	ApprovedCreditClasses []ecocredit.CreditClassID `json:"credit_classes"`
	TotalLandAllocations  sdk.Int                   `json:"total_land_allocations"`
}

type MsgCreateReDAOMint struct {
	ReDAOMintMetadata
	Founder       sdk.AccAddress `json:"founder"`
	FounderShares sdk.Int        `json:"founder_shares"`
}

type MsgMintShares struct {
	ReDAOMint sdk.AccAddress `json:"re_dao_mint"`
	Shares    sdk.Int        `json:"shares"`
}

type LandAllocation struct {
	ReDAOMint   sdk.AccAddress `json:"re_dao_mint"`
	LandSteward sdk.AccAddress `json:"land_steward"`
	GeoPolygon  []byte         `json:"geo_polygon"`
	Allocation  sdk.Int        `json:"allocation"`
}

type MsgAllocateLandShares struct {
	LandAllocation
}

type MsgDistributeCredit struct {
	ReDAOMint sdk.AccAddress
	Credit    ecocredit.CreditID
}

type MsgDistributeFunds struct {
	ReDAOMint sdk.AccAddress
	Funds     sdk.Coins
}

type ProposalID []byte

type Proposal struct {
	ReDAOMint sdk.AccAddress `json:"re_dao_mint"`
	Msgs      []sdk.Msg      `json:"msgs"`
	Proposer  sdk.AccAddress `json:"proposer"`
}

type MsgPropose struct {
	Proposal
}

type MsgVote struct {
	ProposalID ProposalID     `json:"proposal_id"`
	Voter      sdk.AccAddress `json:"voter"`
	Vote       bool           `json:"vote"`
}

type MsgExecProposal struct {
	ProposalID ProposalID     `json:"proposal_id"`
	Signer     sdk.AccAddress `json:"signer"`
}

func (m MsgCreateReDAOMint) Route() string {
	return RouterKey
}

func (m MsgCreateReDAOMint) Type() string {
	return "create-redaomint"
}

func (m MsgCreateReDAOMint) ValidateBasic() sdk.Error {
	if m.Founder.Empty() {
		return sdk.ErrInvalidAddress(DefaultCodespace)
	}
	if !m.FounderShares.IsPositive() {
		return sdk.ErrUnknownRequest("invalid founder shares")
	}

	return nil
}

func (m MsgCreateReDAOMint) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgCreateReDAOMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Founder}
}

func (m MsgMintShares) Route() string {
	return RouterKey
}

func (m MsgMintShares) Type() string {
	return "mint-shares"
}

func (m MsgMintShares) ValidateBasic() sdk.Error {
	panic("implement me")
}

func (m MsgMintShares) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgMintShares) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.ReDAOMint}
}

func (m MsgAllocateLandShares) Route() string {
	return RouterKey
}

func (m MsgAllocateLandShares) Type() string {
	return "allocate-landshares"
}

func (m MsgAllocateLandShares) ValidateBasic() sdk.Error {
	if m.ReDAOMint.Empty() {
		return sdk.ErrInvalidAddress(DefaultCodespace)
	}
	if m.LandSteward.Empty() {
		return sdk.ErrInvalidAddress(DefaultCodespace)
	}
	if !(len(m.GeoPolygon) > 0) {
		return sdk.ErrUnknownRequest("invalid geo polygon")
	}
	if !m.Allocation.IsPositive() {
		return sdk.ErrUnknownRequest("invalid allocation")
	}

	return nil
}

func (m MsgAllocateLandShares) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgAllocateLandShares) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.ReDAOMint, m.LandSteward}
}

func (m MsgPropose) Route() string {
	return RouterKey
}

func (m MsgPropose) Type() string {
	return "propose"
}

func (m MsgPropose) ValidateBasic() sdk.Error {
	if m.ReDAOMint.Empty() {
		return sdk.ErrInvalidAddress(DefaultCodespace)
	}
	if !(len(m.Msgs) > 0) {
		return sdk.ErrUnknownRequest("invalid number of messages")
	}
	if m.Proposer.Empty() {
		return sdk.ErrInvalidAddress(DefaultCodespace)
	}

	return nil
}

func (m MsgPropose) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgPropose) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Proposer, m.ReDAOMint}
}

func (m MsgVote) Route() string {
	return RouterKey
}

func (m MsgVote) Type() string {
	return "vote"
}

func (m MsgVote) ValidateBasic() sdk.Error {
	if !(len(m.ProposalID) > 0) {
		return sdk.ErrUnknownRequest("invalid proposal id")
	}
	if m.Voter.Empty() {
		return sdk.ErrInvalidAddress(DefaultCodespace)
	}
	return nil
}

func (m MsgVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Voter}
}

func (m MsgExecProposal) Route() string {
	return RouterKey
}

func (m MsgExecProposal) Type() string {
	return "exec-proposal"
}

func (m MsgExecProposal) ValidateBasic() sdk.Error {
	if !(len(m.ProposalID) > 0) {
		return sdk.ErrUnknownRequest("invalid proposal id")
	}
	if m.Signer.Empty() {
		return sdk.ErrInvalidAddress(DefaultCodespace)
	}

	return nil
}

func (m MsgExecProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgExecProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

func (a LandAllocation) ID() []byte {
	return []byte(fmt.Sprintf("%x/%x", a.ReDAOMint, a.GeoPolygon))
}
