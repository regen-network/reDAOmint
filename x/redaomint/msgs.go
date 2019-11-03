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
	Founder sdk.AccAddress `json:"founder"`
	FounderShares sdk.Int `json:"founder_shares"`
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
	panic("implement me")
}

func (m MsgCreateReDAOMint) Type() string {
	panic("implement me")
}

func (m MsgCreateReDAOMint) ValidateBasic() sdk.Error {
	panic("implement me")
}

func (m MsgCreateReDAOMint) GetSignBytes() []byte {
	panic("implement me")
}

func (m MsgCreateReDAOMint) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

func (m MsgMintShares) Route() string {
	panic("implement me")
}

func (m MsgMintShares) Type() string {
	panic("implement me")
}

func (m MsgMintShares) ValidateBasic() sdk.Error {
	panic("implement me")
}

func (m MsgMintShares) GetSignBytes() []byte {
	panic("implement me")
}

func (m MsgMintShares) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

func (m MsgAllocateLandShares) Route() string {
	panic("implement me")
}

func (m MsgAllocateLandShares) Type() string {
	panic("implement me")
}

func (m MsgAllocateLandShares) ValidateBasic() sdk.Error {
	panic("implement me")
}

func (m MsgAllocateLandShares) GetSignBytes() []byte {
	panic("implement me")
}

func (m MsgAllocateLandShares) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

func (m MsgPropose) Route() string {
	panic("implement me")
}

func (m MsgPropose) Type() string {
	panic("implement me")
}

func (m MsgPropose) ValidateBasic() sdk.Error {
	panic("implement me")
}

func (m MsgPropose) GetSignBytes() []byte {
	panic("implement me")
}

func (m MsgPropose) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

func (m MsgVote) Route() string {
	panic("implement me")
}

func (m MsgVote) Type() string {
	panic("implement me")
}

func (m MsgVote) ValidateBasic() sdk.Error {
	panic("implement me")
}

func (m MsgVote) GetSignBytes() []byte {
	panic("implement me")
}

func (m MsgVote) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

func (m MsgExecProposal) Route() string {
	panic("implement me")
}

func (m MsgExecProposal) Type() string {
	panic("implement me")
}

func (m MsgExecProposal) ValidateBasic() sdk.Error {
	panic("implement me")
}

func (m MsgExecProposal) GetSignBytes() []byte {
	panic("implement me")
}

func (m MsgExecProposal) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

func (a LandAllocation) ID() []byte {
	return []byte(fmt.Sprintf("%x/%x", a.ReDAOMint, a.GeoPolygon))
}
