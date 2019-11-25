package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
	"muse/x/muse/types"
	"github.com/cosmos/cosmos-sdk/codec"
)
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	museQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the muse module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	museQueryCmd.AddCommand(client.GetCommands(
		GetCmdLyric(storeKey, cdc),
	)...)
	return museQueryCmd
}
func GetCmdLyric(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "lyric [lyricCode]",
		Short: "query lyric",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			code := args[0]
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/lyric/%s", queryRoute, code), nil)
			if err != nil {
				fmt.Printf("could not query lyricCode - %s \n", code)
				return nil
			}
			var out types.QueryResLyric
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
