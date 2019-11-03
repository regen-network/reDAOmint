# Project Overview

This is the Regen network team's submission for the DeFi hackathon. Our objective is to use decentralized finance to realize the living capital potential of the various projects around the globe and turn them into more understandable investment vehicles for both those interested in investing for social good and for self-interested parties as well. This is backed by the observation that there are many low-hanging fruit - projects that sequester carbon for less than the future carbon tax, etc. This opportunity for the risk-tolerant to secure gains across both monetary and immaterial valuables is a once-in-a-lifetime opportunity that has yet to be made avalable as a product to retail investors.

## System Objectives

Our system is designed to foster the discovery, continual funding and growth of impactful projects across the globe. We wish to connect land stewards with new contribution streams in an open marketplace that punishes unacceptable behaviors and rewards both investor and good stewards of the land.

## Design

We rely on a mix of currencies offered for trade on our managed DEX. This DEX cultivates flows of value by managing buybacks and limiting trade pairs. In this way, our purpose-built staking token represents inclusion in the DAO's endowment fund, and a stake in the positive outcomes of DAO activities. Careful definition of the various actors in the system allows us to model it as a [three-sided market](https://github.com/BlockScience/cadCAD-Tutorials/blob/master/02%20Reference%20Models/ThreeSidedBasic/3%20sided%20model.ipynb) producing Argument of Impact backed living capital tokens. It matches DAO funds to land stewards through a marketplace for argument providers.

## Tokens

### Contributed tokens

The DEX accepts any token the DAO has been configured to manage: first Cosmos tokens like TRUEstory and Atoms, perhaps bitcoins and litecoin later.

### EnDAOmint tokens

Roughly represent a share in the enDAOmint and a stake in the living capital outputs from its activities. This token has a set limit to encourage value accrual by early-adopters. Once all tokens are minted, the only entrypoint is to buy tokens from current holders. The enDAOmint gets a portion of all minted tokens in exchange for development and managment of the platform.

### undifferentiated contributor credits

Expiring quarterly rewards distributed to every holder of EnDAOmint tokens. They can be redeemed by trading them on the DEX for living capital tokens issued by various Authorities. The DAO then buys back the undifferentiated contributor credits from the Authorities at a rate set at the beginning of the quarter. They must be redeemed within the next quarter after they were issued and are always issued by the end of the quarter.

### Living capital tokens

Represent positive outcomes that the DAO community has voted to value. These tokens can be issued by multiple Authorities (many authorities, one AOI token) or a single Authority, but no Authority can issue multiple types of token against one stake.

## Actors

### Investors

Contribute paired funds to the DAO by buying the Ecoshare token from our DEX. Ecoshares are minted to match demand, and contributed funds are transferred from the DEX to be held in the DAO funding basket.

### Authorities

Authorities are companies, government agencies or perhaps even self-representing stewards that produce Arguments Of Impact (AOI) and issue living capital tokens under their own authority. Each Authority can have one or more Land Stewards producing against managed land towards the Authority's primary objective. A few examples: An IPCC-run carbon crediting authority, a remote sensing company that creates biomass indexes and manages forests and wetlands through our system. The Authorities must maintain a stake for slashing by the DAO community that is proportional to their expected living capital tokens. 

### Stewards

Stewards are a user class tied directly to plots of land. Mostly, they will represent existing projects tied by context to an authority that manages the Arguments Of Impact on their behalf. An example could be an indigenous tribe that manages ancestral lands, a forestry commons, a logging company trying to improve its impact or a government-run remediation project. 

### Governance

Because we have funds under management, the governance is a bit more complicated. If we are not managing a basket fund then governance could be a simple direct democracy. Instead we elect for a hybrid model based on prediction markets where users vote to approve or continue the approval of authorities, and overall system objectives which are mostly autonomous. The Regen team will take on the task of balancing the basket across assets and ensuring that the fund distribution process stays in sync with community-voted system objectives. The open-source, forkable nature of the project ensures that no user or project is doomed by a misbehaving managment team.

## User stories {{NEEDS WORK}}

### ACME Reforestation Group (Authority)

* Issues ARG tokens under an argument of positive plant biomass impact they process quarterly using a government-approved model and open landsat data.
* Manages AOI submission for a government-funded reforestation project in Papa New Guinea
* Manages AOI for a community-supported-agriculture consortium in Ecuador
* Issues and buys back ARG tokens on the DEX with trading pairs to: the stablecoin and the undifferentiated investor dividend token.

### Larry (investor)

* Discovers the DAO through a relative that gifted him enDAOmint tokens
* Sets a buy order to purchase more enDAOmint tokens at the beginning of each month
* Deeply concerned about climate change and the wellbeing of elephants.
* Buys biomass and carbon sequestration tokens
* Advocates for new living capital tokens to represent elephant interests (habitat reclamation, elephant rescue operations)

## Interaction

1. Larry recieves his quarterly dividend in the form of 120 undifferentiated contributor credits.
2. Quarterly harvest for ARG is positive - they issue corresponding percentages of stablecoin outputs to their managed project and puts the AOI-tokens for biomass (ARG) on the exchange to trade against contributor credits.
3. Larry goes to the DEX, and redeems most of his credits for biomass tokens like ARG and carbon credits issued by IPCC.
4. Larry sends a portion of his credits as a bounty for inclusion of an elephant habitat authority.
5. Stewards of ARG projects submit claims for expected biomass increase next quarter (or allow ARG to do it on their behalf).
