package redaomint

import (
	"fmt"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/gaia/x/ecocredit"
	"github.com/spf13/cobra"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   ModuleName,
		Short: "reDAOmint transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		GetCmdCreateReDAOMint(cdc),
	)...)

	return txCmd
}

func GetCmdCreateReDAOMint(cdc *codec.Codec) *cobra.Command {
	var creditClassStrs []string
	cmd := &cobra.Command{
		Use:   "create-redaomint [description] (--credit-class [credit-class])*",
		Args:  cobra.ExactArgs(1),
		Short: "create a new ReDAOMint",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			from := cliCtx.GetFromAddress()

			var creditClasses []ecocredit.CreditClassID
			for _, cls := range creditClassStrs {
				hrp, bz, err := bech32.Decode(cls)
				if err != nil {
					return err
				}
				if hrp != "ecocls" {
					return fmt.Errorf("not a credit class %s", cls)
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
