package keeper
import (
	"github.com/cosmos/cosmos-sdk/codec"
	"muse/x/muse/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)
const (
	QueryLyric = "lyric"
	QueryLyrics = "lyrics"
	QueryWork = "work"
)
// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryLyric:
			return queryLyric(ctx, path[1:], req, keeper)
		//case QueryLyrics:
			//return queryLyrics(ctx, path[1:], req, keeper)
		//case QueryWork:
			//return queryWork(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown muse query endpoint")
		}
	}
}

func queryLyric(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error)  {
	val := keeper.QueryLyric(ctx, path[0])
	if val == "" {
		return []byte{}, sdk.ErrUnknownRequest("could not get lyric")
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, types.QueryResLyric{Value: val})
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return res, nil
}