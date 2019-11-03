package ecocredit

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"time"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   ModuleName,
		Short: "reDAOmint transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		GetCmdCreateCreditClass(cdc),
		GetCmdIssueCredit(cdc),
	)...)

	return txCmd
}

func GetCmdCreateCreditClass(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-class [name] [issuers]",
		Args:  cobra.ExactArgs(2),
		Short: "create a new credit class",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			from := cliCtx.GetFromAddress()

			var issuers []sdk.AccAddress
			for _, bech := range strings.Split(args[1], ",") {
				addr, err := sdk.AccAddressFromBech32(bech)
				if err != nil {
					return err
				}
				issuers = append(issuers, addr)
			}

			msg := MsgCreateCreditClass{CreditClassMetadata{Designer: from, Name: args[0], Issuers: issuers}}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{})
		},
	}
	return cmd
}

func GetCmdIssueCredit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue [credit-class] [geo-polygon] [start-date] [end-date] [units] [holder]",
		Args:  cobra.ExactArgs(6),
		Short: "issue a new ecosystem service credit",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			from := cliCtx.GetFromAddress()

			creditClass, err := CreditClassFromBech32(args[0])
			if err != nil {
				return err
			}

			startDate, err := time.Parse("2006-01-02T15:04:05-0700", args[2])
			if err != nil {
				return err
			}

			endDate, err := time.Parse("2006-01-02T15:04:05-0700", args[3])
			if err != nil {
				return err
			}

			holder, err := sdk.AccAddressFromBech32(args[5])
			if err != nil {
				return err
			}

			units, err := sdk.NewDecFromStr(args[4])
			if err != nil {
				return err
			}

			msg := MsgIssueCredit{CreditMetadata{
				Issuer:      from,
				CreditClass: creditClass,
				GeoPolygon:  []byte(args[1]),
				StartDate:   startDate,
				EndDate:     endDate,
				Units:       units,
			},
				holder,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{})
		},
	}
	return cmd
}
