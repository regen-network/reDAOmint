package redaomint

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/gaia/x/ecocredit"
	"github.com/spf13/cobra"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   ModuleName,
		Short: "reDAOmint transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		GetCmdCreateReDAOMint(cdc),
		GetCmdContributeReDAOMint(cdc),
	)...)

	return txCmd
}

func GetCmdCreateReDAOMint(cdc *codec.Codec) *cobra.Command {
	var creditClassStrs []string
	cmd := &cobra.Command{
		Use:   "create [description] (--credit-class [credit-class])*",
		Args:  cobra.ExactArgs(1),
		Short: "create a new ReDAOMint",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			from := cliCtx.GetFromAddress()

			var creditClasses []ecocredit.CreditClassID
			for _, cls := range creditClassStrs {
				bz, err := ecocredit.CreditClassFromBech32(cls)
				if err != nil {
					return err
				}
				creditClasses = append(creditClasses, bz)
			}

			msg := MsgCreateReDAOMint{ReDAOMintMetadata{args[0], creditClasses}, from}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{})
		},
	}
	cmd.Flags().StringArrayVar(&creditClassStrs, "credit-class", nil, "set an approved credit class for the reDAOmint")
	return cmd
}

func GetCmdContributeReDAOMint(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contribute [redaomint] [funds]",
		Args:  cobra.ExactArgs(2),
		Short: "contribute to a ReDAOMint",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			from := cliCtx.GetFromAddress()

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := MsgContributeReDAOMint{Sender: from, ReDAOMint: addr, Funds: coins}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{})
		},
	}
	return cmd
}
