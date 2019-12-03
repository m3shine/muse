package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"muse/x/muse/types"
	"net/http"
)

type setLyricReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	LyricCode string       `json:"lyric_code"`
	Author    string       `json:"author"`
	Title     string       `json:"title"`
	Hash      []byte       `json:"hash"`
	Owner     string       `json:"owner"`
	TokenName string       `json:"token_name"`
}

func setLyricHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req setLyricReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgSetLyric(req.LyricCode, req.Author, req.Title, req.Hash, addr, sdk.Coins{sdk.NewInt64Coin(req.TokenName, types.TokenAmount)})
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
