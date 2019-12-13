package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/lib/server"
	"github.com/tendermint/tendermint/rpc/lib/types"
	"net/http"
)

func lyricHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bciCode := vars[bciCode]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("/custom/%s/lyric/%s", storeName, bciCode), nil)
		if err != nil {
			rpcserver.WriteRPCResponseHTTP(w, rpctypes.RPCInvalidRequestError(rpcid, err))
			return
		}
		rpcserver.WriteRPCResponseHTTP(w, rpctypes.NewRPCSuccessResponse(cliCtx.Codec, rpcid, string(res)))

	}
}

func musicHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bciCode := vars[bciCode]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("/custom/%s/music/%s", storeName, bciCode), nil)
		if err != nil {
			rpcserver.WriteRPCResponseHTTP(w, rpctypes.RPCInvalidRequestError(rpcid, err))
			return
		}
		rpcserver.WriteRPCResponseHTTP(w, rpctypes.NewRPCSuccessResponse(cliCtx.Codec, rpcid, string(res)))

	}
}