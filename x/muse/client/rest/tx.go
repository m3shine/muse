package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"log"
	"muse/x/muse/types"
	"net/http"
)


type Sign struct {
	FromName string `json:"from_name"`
	Password string `json:"password"`
}

type Req struct {
	rest.BaseReq `json:"base_req"`
	Sign         `json:"sign"`
}

type setLyricReq struct {
	Req       `json:"req"`
	LyricCode string `json:"lyric_code"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Hash      []byte `json:"hash"`
	TokenName string `json:"token_name"`
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

		addr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
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
		sendTx(w, cliCtx, req.BaseReq, req.Sign, []sdk.Msg{msg})
	}
}

/**
curl -XPUT -s http://localhost:1317/muse/lyric --data-binary '{"req":{"base_req":{"from":"cosmos193jutxkx74xx8yaufcx9pcp3cwd90nwsegklpa","chain_id":"musenetwork"},"sign":{"from_name":"jon","password":"9ijn8uhb"}},"lyric_code":"d0911","author":"jon","title":"helloworld","hash":"333333333333333333333333","token_name":"hwtoken"}'
 */
func sendTx( w http.ResponseWriter,cliCtx context.CLIContext, br rest.BaseReq, sign Sign, msgs []sdk.Msg) {

	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, br.GasAdjustment, flags.DefaultGasAdjustment)
	if !ok {
		return
	}

	simAndExec, gas, err := flags.ParseGas(br.Gas)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	txBldr := authtypes.NewTxBuilder(
		utils.GetTxEncoder(cliCtx.Codec), br.AccountNumber, br.Sequence, gas, gasAdj,
		br.Simulate, br.ChainID, br.Memo, br.Fees, br.GasPrices,
	)

	if br.Simulate || simAndExec {
		if gasAdj < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid gas adjustment")
			return
		}

		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if br.Simulate {
			rest.WriteSimulationResponse(w, cliCtx.Codec, txBldr.Gas())
			return
		}
	}
	addr, err := sdk.AccAddressFromBech32(br.From)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	cliCtx.FromAddress = addr
	//准备好正确的accountNumber和sequence
	txBldr, err = utils.PrepareTxBuilder(txBldr, cliCtx)

	stdMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	//新建交易
	tx := auth.NewStdTx(stdMsg.Msgs, stdMsg.Fee , nil, stdMsg.Memo)

	//对交易签名
	signedTx, err := txBldr.SignStdTx(sign.FromName, sign.Password, tx, false)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	//交易签名后的字节
	txBytes, err := utils.GetTxEncoder(cliCtx.Codec)(signedTx)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// broadcast to a Tendermint node
	txRes, err := cliCtx.BroadcastTxCommit(txBytes)

	output, err := cliCtx.Codec.MarshalJSON(txRes)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(output); err != nil {
		log.Printf("could not write response: %v", err)
	}
	//success
	//{"height":"0","txhash":"4224962F1A55B3D27B03A89B1FB4022C7DAD0F985FE15E014573B0BFC8A17F42","raw_log":"[{\"msg_index\":0,\"success\":true,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"set_lyric\"}]}]}]","logs":[{"msg_index":0,"success":true,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"set_lyric"}]}]}]}
	//error
	//{"height":"0","txhash":"8B41A79B1CED9CF238C4BDC2328369D44A57F3A60D226B166F8F8D5C8D49CD23","code":4,"raw_log":"{\"codespace\":\"sdk\",\"code\":4,\"message\":\"signature verification failed; verify correct account sequence and chain-id\"}"}
}