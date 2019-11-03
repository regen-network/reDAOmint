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
	"time"
)

type Keeper struct {
	cdc             *codec.Codec
	storeKey        sdk.StoreKey
	accountKeeper   auth.AccountKeeper
	bankKeeper      bank.Keeper
	supplyKeeper    supply.Keeper
	ecocreditKeeper ecocredit.Keeper
	ibcKeeper       ibc.Keeper
	router          sdk.Router
	metadataBucket  orm.AutoIDBucket
	landAllocations orm.NaturalKeyBucket
	proposalBucket  orm.AutoIDBucket
	votesBucket     orm.NaturalKeyBucket
}

const (
	IndexByReDAOMint = "by-redaomint"
	IndexByProposal  = "by-proposal"
)

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, supplyKeeper supply.Keeper, ecocreditKeeper ecocredit.Keeper, ibcKeeper ibc.Keeper, router sdk.Router) Keeper {
	return Keeper{cdc: cdc,
		storeKey:        storeKey,
		accountKeeper:   accountKeeper,
		bankKeeper:      bankKeeper,
		supplyKeeper:    supplyKeeper,
		ecocreditKeeper: ecocreditKeeper,
		ibcKeeper:       ibcKeeper,
		router: router,
		metadataBucket:  orm.NewAutoIDBucket(storeKey, "metadata", cdc, nil),
		landAllocations: orm.NewNaturalKeyBucket(storeKey, "allocations", cdc, []orm.Index{
			{Name: IndexByReDAOMint, Indexer: func(key []byte, value interface{}) (indexValue []byte, err error) {
				allocation := value.(LandAllocation)
				return allocation.ReDAOMint, nil
			}},
		}),
		proposalBucket: orm.NewAutoIDBucket(storeKey, "proposal", cdc, nil),
		votesBucket: orm.NewNaturalKeyBucket(storeKey, "votes", cdc, []orm.Index{
			{IndexByProposal,
				func(key []byte, value interface{}) (indexValue []byte, err error) {
					vote := value.(Vote)
					return vote.Proposal, nil
				},
			},
		}),
	}
}

// Denom returns the token denomination for a reDAOmint
func Denom(redaomint sdk.AccAddress) string {
	return fmt.Sprintf("redao:%x", redaomint)
}

// CreateReDAOMint creates a new reDAOmint account and token denomination for the reDAOmint.
// This event also distributes founder shares to the founder of the reDAOmint
func (k Keeper) CreateReDAOMint(ctx sdk.Context, metadata ReDAOMintMetadata, founder sdk.AccAddress, founderShares sdk.Int) (addr sdk.AccAddress, denom string, err error) {
	metadata.TotalLandAllocations = sdk.NewInt(0)
	addr, err = k.metadataBucket.Create(ctx, metadata)
	if err != nil {
		return nil, "", err
	}
	k.accountKeeper.SetAccount(ctx, &auth.BaseAccount{Address: addr})
	return addr, Denom(addr), err
}

// MintShares mints new shares for the reDAOmint and assigns them to the reDAOmint pool to be sold on a DEX
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

// SetLandAllocation gives a land steward on a specific piece of land some fractional allocation of the rewards
// in the reDAOmint. The exact fractional value of an allocation is up to the reDAOmint
func (k Keeper) SetLandAllocation(ctx sdk.Context, allocation LandAllocation) error {
	var metadata ReDAOMintMetadata
	err := k.metadataBucket.GetOne(ctx, allocation.ReDAOMint, &metadata)
	if err != nil {
		return err
	}
	// look for an existing allocation
	var existing LandAllocation
	err = k.landAllocations.GetOne(ctx, &existing)
	if err != nil {
		metadata.TotalLandAllocations = metadata.TotalLandAllocations.Sub(existing.Allocation)
	}
	if allocation.Allocation.IsZero() {
		// delete if the new allocation is zero
		err = k.landAllocations.Delete(ctx, allocation)
		if err != nil {
			return err
		}
	} else {
		err = k.landAllocations.Save(ctx, allocation)
		if err != nil {
			return err
		}
		metadata.TotalLandAllocations = metadata.TotalLandAllocations.Add(allocation.Allocation)
	}
	// update the master reDAOmint metadata
	err = k.metadataBucket.Save(ctx, allocation.ReDAOMint, metadata)
	if err != nil {
		return err
	}
	return nil
}

// DistributeCredit distributes fractional shares of a credit held by the reDAOmint to all reDAOmint
// shareholders. This event should ideally happen on a pre-defined schedule within a begin blocker
// for instance
func (k Keeper) DistributeCredit(ctx sdk.Context, redaomint sdk.AccAddress, credit ecocredit.CreditID) error {
	holding, found := k.ecocreditKeeper.GetCreditHolding(ctx, credit, redaomint)
	if !found {
		return fmt.Errorf("not found")
	}
	denom := Denom(redaomint)
	totalCoins := k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)
	var err error
	// TODO: figure out something more efficient for this
	k.accountKeeper.IterateAccounts(ctx, func(account exported.Account) (stop bool) {
		coins := account.GetCoins().AmountOf(denom)
		var share sdk.Dec
		share.Div(coins.BigInt(), totalCoins.BigInt())
		units := holding.LiquidUnits.Mul(share)
		if units.IsPositive() {
			err := k.ecocreditKeeper.SendCredit(ctx, credit, redaomint, account.GetAddress(), units)
			if err != nil {
				return true
			}
		}
		return false
	})
	return err
}

// VerifyOrSlashLandStewards cycles through all of the land allocations and checks if there is a credit for that piece
// of land for that time window in the class of approved credits, if not, the land allocation is slashed from the
// pool of land allocations for this reDAOmint and can no longer receive rewards. Receipt of an approved credit is
// required to keep receiving rewards. In the future the start and end dates would be set more automatically and
// this process would be run on a schedule
func (k Keeper) VerifyOrSlashLandStewards(ctx sdk.Context, redaomint sdk.AccAddress, startDate time.Time, endDate time.Time) error {
	iterator, err := k.landAllocations.ByIndexPrefixScan(ctx, IndexByReDAOMint, nil, nil, false)
	if err != nil {
		return nil
	}
	for {
		var allocation LandAllocation
		_, err = iterator.LoadNext(&allocation)
		if err != nil {
			break
		}
		found := false
		k.ecocreditKeeper.IteratorCreditsByGeoPolygon(ctx, allocation.GeoPolygon, func(metadata ecocredit.CreditMetadata) (stop bool) {
			// TODO: make this more robust so that different credits could span these dates
			if (metadata.StartDate.Before(startDate) || metadata.StartDate.Equal(startDate)) &&
				(metadata.EndDate.After(endDate) || metadata.EndDate.Equal(endDate)) {
				found = true
				return true
			}
			return false
		})
		if !found {
			allocation.Allocation = sdk.NewInt(0)
			_ = k.SetLandAllocation(ctx, allocation)
		}
	}
	return nil
}

// DistributeFunds distributes funds from the reDAOmint to land stewards (who have not been slashed). Ideally this
// would be some sort of stable coin that has been acquired on a DEX by converting dividends/block rewards from
// assets within the reDAOmint pool
func (k Keeper) DistributeFunds(ctx sdk.Context, redaomint sdk.AccAddress, funds sdk.Coins) error {
	total := k.bankKeeper.GetCoins(ctx, redaomint)
	_, insufficientFunds := total.SafeSub(funds)
	if insufficientFunds {
		return fmt.Errorf("insufficient funds")
	}
	var metadata ReDAOMintMetadata
	err := k.metadataBucket.GetOne(ctx, redaomint, &metadata)
	if err != nil {
		return err
	}
	totalAllocations := metadata.TotalLandAllocations
	iterator, err := k.landAllocations.ByIndexPrefixScan(ctx, IndexByReDAOMint, nil, nil, false)
	if err != nil {
		return nil
	}
	for {
		var allocation LandAllocation
		_, err = iterator.LoadNext(&allocation)
		if err != nil {
			break
		}
		var share sdk.Dec
		share.Div(allocation.Allocation.BigInt(), totalAllocations.BigInt())
		for _, coin := range funds {
			amount := share.MulInt(coin.Amount).TruncateInt()
			err := k.bankKeeper.SendCoins(ctx, redaomint, allocation.LandSteward, sdk.Coins{{Denom: coin.Denom, Amount: amount}})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k Keeper) CreateProposal(ctx sdk.Context, proposal Proposal) (ProposalID, error) {
	id, err := k.proposalBucket.Create(ctx, proposal)
	if err != nil {
		return id, err
	}
	return id, nil
}

type Vote struct {
	Proposal ProposalID
	Voter    sdk.AccAddress
	Vote     bool
}

func (v Vote) ID() []byte {
	return []byte(fmt.Sprintf("%x/%x", v.Proposal, v.Voter))
}

func (k Keeper) Vote(ctx sdk.Context, proposal ProposalID, voter sdk.AccAddress, vote bool) error {
	return k.votesBucket.Save(ctx, Vote{
		Proposal: proposal,
		Voter:    voter,
		Vote:     false,
	})
}

func (k Keeper) ExecProposal(ctx sdk.Context, id ProposalID) sdk.Result {
	var proposal Proposal
	err := k.proposalBucket.GetOne(ctx, id, &proposal)
	if err != nil {
		return sdk.ResultFromError(err)
	}
	denom := Denom(proposal.ReDAOMint)

	var votes sdk.Int
	iterator, err := k.votesBucket.ByIndex(ctx, IndexByProposal, id)
	if err != nil {
		return sdk.ResultFromError(err)
	}
	for {
		var vote Vote
		_, err := iterator.LoadNext(&vote)
		if err != nil {
			break
		}
		coins := k.bankKeeper.GetCoins(ctx, vote.Voter).AmountOf(denom)
		votes.Add(coins)
	}

	totalSupply := k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)

	var percentYes sdk.Dec
	percentYes.Div(votes.BigInt(), totalSupply.BigInt())
	threshold, _ := sdk.NewDecFromStr("0.50")
	if percentYes.LTE(threshold) {
		return sdk.ResultFromError(fmt.Errorf("proposal didn't pass"))
	}
	var res sdk.Result
	for _, msg := range proposal.Msgs {
		res = k.router.Route(msg.Route())(ctx, msg)
		if !res.IsOK() {
			return res
		}
	}
	return res
}
