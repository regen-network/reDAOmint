# The ReDAOmint impact endowment protocol

## Inspiration

Protecting Earth’s most precious natural environments - such as the Amazon rainforest - is often dangerously dependent on unstable political and economic forces. While the Amazon’s critical importance as a carbon sink, source of oxygen, and biodiversity hotspot is widely acknowledged worldwide, the local economics in inhabited portions of the forest are often not conducive to protecting this planetary treasure. Instead of relying on government policies and inefficient charities, what if we could create a decentralized economic vehicle that made long-term protection of the forest economically beneficial for local residents?

## Project Overview

This is the Regen network team's submission for the DeFi hackathon. Our objective is to use decentralized finance to maximize and secure the living capital potential of impact projects around the globe and craft them into rational investment vehicles for both those interested in investing for social good and self-interest. This is an opportunity for the risk-tolerant to secure gains across both monetary and immaterial valuables in a once-in-a-lifetime underappreciated inflection point - it will never be easier to draw carbon into soils, recover endangered species, etc. To this end we suggest to use a multi-currency blockchain project to secure and manage an endowment - first for cultivating positive outcomes in rainforests, with the goal of eventually offering a marketplace to cultivate any human value.

A **reDAOmint** (*regenerative + endowment + DAO + tendermint*) is a shareholder DAO that holds a pool of assets for perpetuity, distributing the dividends of those assets with land stewards in critical bioregions in exchange for protecting the land. The DAO funds itself by minting new shares in exchange for valuable market assets that produce dividends or rewards (such as staking tokens).

## System Objectives

Our system is designed to foster the discovery, continual funding and growth of impactful projects across the globe. We wish to connect land stewards with new contribution streams in an open marketplace that punishes unacceptable behaviors and rewards both investor and good stewards of the land. Transparency, consistency and a sound financial basis are scarce in the impact industry, making it a compelling use-case for those interested in blockchain applications. Blockchain is inheritly powered by electricity, so some externalities on the climate are to be expected. The synergies between the two will lead to powerful new investment mechanisms and impact outcomes.
profit support could provide.

For investors in the **reDAOmint**, the investment is not, however, simply a charity and could also be consider a decentralize public hedgefund. **reDAOmint** shares provide dividends that come in the form of ecosystem service credits. These credits represent unaccounted ecological benefits from protecting the land which is an emerging asset class  estimated to grow from 36B USD in 2017 to 1T USD over next 20 years. An example of these are carbon credits, which will be increasingly demanded by governments and the public. **reDAOmint** shareholders can either use the ecosystem service credits they have received as dividends for their own offseting purposes or sell them on the market.

![](https://i.imgur.com/v7lKV1q.png)

In comparison to traditional mechanisms for protecting land, a **reDAOmint** provides the following benefits:

* insulates land stewards from seasonal variation in the ecological health of land
* insulates land stewards from volatility in the price of ecosystem service credits
* efficiently generates ecosystem services by good capital management
* auditability of all funding and verification activity

## Construction

is built using the Cosmos SDK using a fork of cosmos/gaia that includes IBC support. A diff of our work can be seen [here](https://github.com/regen-network/reDAOmint/pull/3).

We produced two new Cosmos modules and an ORM package to create the **reDAOmint** for this hackathon.

### `redaomint` module

The [`redaomint` module](https://github.com/regen-network/reDAOmint/tree/reDAOmint/x/redaomint) implements a shareholder DAO that produces dividends of ecosystem service credits for shareholders and dividends of the underlying asset pool for allocated land stewards as described above. It also provides governance proposal and voting support for shareholders. It uses the existing Cosmos `bank` and `supply` modules to track holdings of shares, and interacts with the `ecocredit` module for verification of good land stewardship and distributing credits. The primary documentation for the module can be found [here](https://github.com/regen-network/reDAOmint/blob/reDAOmint/x/redaomint/keeper.go).

### `ecocredit` module

The [`ecocredit` module](https://github.com/regen-network/reDAOmint/tree/reDAOmint/x/ecocredit) provides a fractional NFT with metadata specific to ecosystem service credits. It allows for:

* the creation of new credit classes with a list of approved issuers
* the issuance of credits for a specific piece of land and time frame
* exchange of fractional portions of individual credits
* burning credits in order to remove them from circulation (in the language of carbon credits this is called retiring and means that you are using the credit as an offset).
Documentation can be found [here](https://github.com/regen-network/reDAOmint/blob/reDAOmint/x/ecocredit/keeper.go)

### `orm` package

In order to make the implementation of the above modules easier, we implemented an [`orm` package](https://github.com/regen-network/reDAOmint/tree/reDAOmint/orm) to handle secondary indexes and the automatic generation of ID's. This is inspired by the [Weave SDK's](https://github.com/iov-one/weave) and is something we've been intending to build for a while. While it was its own "mini-project", it greatly simplified implementation of the other code.

## Design

![3-Sided Market](imgs/threesidedmarket.jpeg?raw=true "Market Diagram")
We rely on a core flow of investor contributions to grow the DAO. This DAO cultivates flows of value by managing authority to issue . In this way, our purpose-built staking token represents inclusion in the DAO's endowment fund, and a stake in the positive outcomes of DAO activities. Careful definition of the actors and allowable actions within the system allows us to model our objectives as a [three-sided market](https://github.com/BlockScience/cadCAD-Tutorials/blob/master/02%20Reference%20Models/ThreeSidedBasic/3%20sided%20model.ipynb) producing living capital tokens backed by Arguments of Impact (AOI). It matches DAO funds to land stewards through a marketplace for argument providers.

### Quarterly Cycle

![DAO cycle](imgs/quarterlyCycle.jpg?raw=true "Cycle Diagram")

The DAO's system evolution is driven by this overturning quarterly cycle. We incorporate game theory concepts like Skin in the Game via authority reputation and land staking, governance is updated via voting on values for the DAO and reviewing proposals for improvement of the process or inclusion of new measures, projects and authorities.

### Tokens

#### Contributed tokens

The DEX accepts any token the DAO has been configured to manage: first Cosmos tokens like TRUEstory and Atoms, perhaps bitcoins and litecoin later.

#### EnDAOmint tokenized shares

Roughly represent a share in the enDAOmint and a stake in the living capital outputs from its activities. This token has a set limit to encourage value accrual by early-adopters. Once all tokens are minted, the only entrypoint is to buy tokens from current holders. The enDAOmint gets a portion of all minted tokens in exchange for development and managment of the platform.

#### Living capital tokens

Represent positive outcomes that the DAO community has voted to value. These tokens are implemented as a fractional NFT containing information about the catalogue, issuer and time period of outcomes they capture. It can be issued by multiple Authorities (many authorities -> one AOI token) or a single Authority, but no Authority can issue multiple types of token against one stake.

#### Stablecoin

A stablecoin is employed for paying land stewards in a denomination that would be useful to them.

### Trading Pairs

#### Rationale

Trading pairs are the primary mechanism for designing DAO interactions. Because our DAO is the only authority for transferring some of the assets (the undifferentiated credits) we can control how the constituent value pools interact to some degree, while still maintaining the open and decentralized nature of the overall protocol.

#### Examples

* EnDAOmint Shares (EDS) can only be sold by the DAO, which accepts offers for supported currencies via the DEX auction method. {ANY -> EDS} pairing. Exclusive to DAO.
* LCT can be traded for any other LCT, but its initial offering on the DEX by an Authority is tied to EDS holder accounts. Once an Authority has produced LCT they are allowed to trade it to the DAO for stablecoins to be delivered to its stewards by the DAO directly, or participate as EDS investors as well. LCT can also be traded on other exchanges. Examples of LCT include biodiversity, biomass, carbon sequestration tokens.

## Actors

### Investors

Contribute paired funds to the DAO by buying the Ecoshare token from our DEX. Ecoshares are minted to match demand, and contributed funds are transferred from the DEX to be held in the DAO funding basket.

### Authorities

Authorities are companies, government agencies or perhaps even self-representing stewards that produce Arguments Of Impact (AOI) and issue living capital tokens under their own authority. Each Authority can have one or more Land Stewards producing against managed land towards the Authority's primary objective. A few examples: An IPCC-run carbon crediting authority, a remote sensing company that creates biomass indexes and manages forests and wetlands through our system. The Authorities must maintain a stake for slashing by the DAO community that is proportional to their expected living capital tokens.

### Stewards

Stewards are a user class tied directly to plots of land. Mostly, they will represent existing projects tied by context to an authority that manages the Arguments Of Impact on their behalf. An example could be an indigenous tribe that manages ancestral lands, a forestry commons, a logging company trying to improve its impact or a government-run remediation project.

### Governance

Because we have funds under management, the governance is a bit more complicated. If we are not managing a basket fund then governance could be a simple direct democracy. Instead we elect for a hybrid model based on prediction markets where users vote to approve or continue the approval of authorities, and overall system objectives in an autonomous process. The Regen team will take on the task of balancing the basket across assets and ensuring that the fund distribution process stays in sync with community-voted system objectives. The open-source, forkable nature of the project ensures that no user or project is doomed by a misbehaving managment team.

## User stories

### ACME Reforestation Group (Authority)

* Issues ARG tokens under an argument of positive plant biomass impact they process quarterly using a government-approved model and open landsat data.
* Manages AOI submission for a government-funded reforestation project in Papa New Guinea
* Manages AOI for a community-supported-agriculture consortium in Ecuador
* Issues and buys back ARG tokens on the DEX with trading pairs to: the stablecoin and the EDS investor token.

### Larry (impact investor)

* Discovers the DAO through a relative that gifted him enDAOmint tokens
* Sets a buy order to purchase more enDAOmint tokens at the beginning of each month
* Deeply concerned about climate change and the wellbeing of elephants.
* Recieves and redeems biomass and carbon sequestration tokens
* Advocates for new living capital tokens to represent elephant interests (habitat reclamation, elephant rescue operations)

### Joe (homo economicus)

* Profit-motivated
* Observes the quarterly issuance of carbon credits are sold below market price.
* Buys in to EDS, gets Carbon credits as quarterly dividend and buys more when they are all dumped on the DEX.
* Sells the credits over the course of the quarter as price recovers for healty profit.

## Interaction

1. Larry recieves his quarterly dividend in the form of 120 undifferentiated contributor credits.
2. Quarterly harvest for ARG is positive - they issue corresponding percentages of stablecoin outputs to their managed project and puts the AOI-tokens for biomass (ARG) on the exchange to trade against contributor credits.
3. Larry redeems most of his credits as biomass tokens like ARG and carbon credits issued by IPCC.
4. Larry is never required to be aware of the land stewards producing ARG, or other investors motivated differently like Joe. But, he could trace his issued living capital tokens via the issuing authorities, or explore other account holder information within the DAO if he wished.
5. Larry sends a portion of his credits as a bounty for inclusion of an elephant habitat authority.
6. Stewards of ARG projects submit claims for expected biomass increase next quarter (or allow ARG to do it on their behalf).
