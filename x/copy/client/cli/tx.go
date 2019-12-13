package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"muse/x/copy/types"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	museTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "copyright transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	museTxCmd.AddCommand(client.PostCommands(
		GetCmdSetLyric(cdc),
	)...)

	return museTxCmd
}

//交易from为owner，owner不作为命令入参
func GetCmdSetLyric(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "set-lyric [bciCode] [author] [content] [title] [hash] [tokenName]",
		Short: "set the lyric that you own",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgSetLyric(args[0], args[1], args[2], args[3],[]byte(args[4]),  cliCtx.GetFromAddress(), sdk.Coins{sdk.NewInt64Coin(args[4], types.TokenAmount)})

			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}