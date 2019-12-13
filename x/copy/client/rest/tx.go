package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/rpc/lib/server"
	"github.com/tendermint/tendermint/rpc/lib/types"
	"math/big"
	"muse/x/copy/types"
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
	Author    string `json:"author"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Hash      string `json:"hash"`
	TokenName string `json:"token_name"`
}

type Stakeholder struct {
	Name     string `json:"name"`
	IDI      string `json:"idi,omitempty"`  //现实身份标识
	Type     string `json:"type,omitempty"` //词曲作者、出版商、发行人、演出团体
	Address  string `json:"address"`
	Describe string `json:"describe,omitempty"`
	Weights  string `json:"weights"` //权重
}

type setMusicReq struct {
	Req          `json:"req"`
	Title        string        `json:"title"`
	Hash         string        `json:"hash"`
	LyricCode    string        `json:"lyric_code"`
	Stakeholders []Stakeholder `json:"stakeholders"`
	WorksID      string        `json:"works_id"` //ISWC、DCI等
}

func setLyricHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req setLyricReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rpcserver.WriteRPCResponseHTTP(w, rpctypes.RPCInvalidRequestError(rpcid, errors.New("failed to parse request")))
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rpcserver.WriteRPCResponseHTTP(w, rpctypes.RPCInvalidRequestError(rpcid, err))
			return
		}

		bciCode, err := types.GenBCI(cliCtx, req.BaseReq.From, 'L')
		if err != nil {
			rpcserver.WriteRPCResponseHTTP(w, rpctypes.RPCInvalidRequestError(rpcid, err))
			return
		}
		// create the message
		msg := types.NewMsgSetLyric(bciCode, req.Author, req.Title, req.Content, []byte(req.Hash), addr, sdk.Coins{sdk.NewInt64Coin(req.TokenName, types.TokenAmount)})
		err = msg.ValidateBasic()
		if err != nil {
			rpcserver.WriteRPCResponseHTTP(w, rpctypes.RPCInvalidRequestError(rpcid, err))
			return
		}
		txRes, err := sendTx(w, cliCtx, req.BaseReq, req.Sign, []sdk.Msg{msg})
		if err != nil {
			rpcserver.WriteRPCResponseHTTP(w, rpctypes.RPCInternalError(rpcid, err))
			return
		}

		res := struct {
			sdk.TxResponse
			BCICode string `json:"bci_code"`
		}{
			txRes,
			bciCode,
		}
		rpcserver.WriteRPCResponseHTTP(w, rpctypes.NewRPCSuccessResponse(cliCtx.Codec, rpcid, res))
	}
}

func setMusicHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req setMusicReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rpcserver.WriteRPCResponseHTTP(w, rpctypes.RPCInvalidRequestError(rpcid, errors.New("failed to parse request")))
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		bciCode, err := types.GenBCI(cliCtx, req.BaseReq.From, 'M')
		if err != nil {
			rpcserver.WriteRPCResponseHTTP(w, rpctypes.RPCInvalidRequestError(rpcid, err))
			return
		}

		var holders []types.Stakeholder
		for _, s := range req.Stakeholders {
			f, _, err := big.ParseFloat(s.Weights, 0, 1000, big.ToNearestEven)
			if err != nil {
				rpcserver.WriteRPCResponseHTTP(w, rpctypes.RPCInvalidRequestError(rpcid, err))
				return
			}

			if s.Address == "" {
				s.Address = req.BaseReq.From
			}

			addr, err := sdk.AccAddressFromBech32(s.Address)
			if err != nil {
				rpcserver.WriteRPCResponseHTTP(w, rpctypes.RPCInvalidRequestError(rpcid, err))
				return
			}

			holders = append(holders, types.Stakeholder{
				Name:     s.Name,
				IDI:      s.IDI,
				Type:     s.Type,
				Address:  addr,
				Describe: s.Describe,
				Weights:  *f,
			})
		}

		// create the message
		msg := types.NewMsgSetMusic(bciCode, req.Title, []byte(req.Hash), req.LyricCode, holders, req.WorksID)
		err = msg.ValidateBasic()
		if err != nil {
			rpcserver.WriteRPCResponseHTTP(w, rpctypes.RPCInvalidRequestError(rpcid, err))
			return
		}
		txRes, err := sendTx(w, cliCtx, req.BaseReq, req.Sign, []sdk.Msg{msg})
		if err != nil {
			rpcserver.WriteRPCResponseHTTP(w, rpctypes.RPCInternalError(rpcid, err))
			return
		}

		res := struct {
			sdk.TxResponse
			BCICode string `json:"bci_code"`
		}{
			txRes,
			bciCode,
		}

		rpcserver.WriteRPCResponseHTTP(w, rpctypes.NewRPCSuccessResponse(cliCtx.Codec, rpcid, res))
	}
}

/**
curl -XPUT -s http://localhost:1317/muse/lyric --data-binary '{"req":{"base_req":{"from":"cosmos193jutxkx74xx8yaufcx9pcp3cwd90nwsegklpa","chain_id":"musenetwork","gas":"auto","memo":"how hard"},"sign":{"from_name":"name","password":"password"}},"lyric_code":"d0911","author":"jon","title":"helloworld","hash":"333333333333333333333333","token_name":"hwtoken"}'
*/

//return [sync]
//success
//{"height":"0","txhash":"4224962F1A55B3D27B03A89B1FB4022C7DAD0F985FE15E014573B0BFC8A17F42","raw_log":"[{\"msg_index\":0,\"success\":true,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"set_lyric\"}]}]}]","logs":[{"msg_index":0,"success":true,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"set_lyric"}]}]}]}
//error
//{"height":"0","txhash":"8B41A79B1CED9CF238C4BDC2328369D44A57F3A60D226B166F8F8D5C8D49CD23","code":4,"raw_log":"{\"codespace\":\"sdk\",\"code\":4,\"message\":\"signature verification failed; verify correct account sequence and chain-id\"}"}

//return [async]
//success
//{"height":"0","txhash":"F648CAD239F3BC30614495F78FAFD8494A1787322FD295CCC8EFD17DA96F5879"}
//error
//{"error":"invalid account password"}

func sendTx(w http.ResponseWriter, cliCtx context.CLIContext, br rest.BaseReq, sign Sign, msgs []sdk.Msg) (sdk.TxResponse, error) {

	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, br.GasAdjustment, flags.DefaultGasAdjustment)
	if !ok {
		return sdk.TxResponse{}, errors.New("gasAdj error")
	}

	simAndExec, gas, err := flags.ParseGas(br.Gas)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	txBldr := authtypes.NewTxBuilder(
		utils.GetTxEncoder(cliCtx.Codec), br.AccountNumber, br.Sequence, gas, gasAdj,
		br.Simulate, br.ChainID, br.Memo, br.Fees, br.GasPrices,
	)

	if br.Simulate || simAndExec {
		if gasAdj < 0 {
			return sdk.TxResponse{}, errors.New("invalid gas adjustment")
		}

		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return sdk.TxResponse{}, err
		}

		if br.Simulate {
			rest.WriteSimulationResponse(w, cliCtx.Codec, txBldr.Gas())
			return sdk.TxResponse{}, nil
		}
	}

	addr, err := sdk.AccAddressFromBech32(br.From)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	//准备好正确的accountNumber和sequence
	txBldr, err = utils.PrepareTxBuilder(txBldr, cliCtx.WithFromAddress(addr))
	if err != nil {
		return sdk.TxResponse{}, err
	}

	stdMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	//新建交易
	tx := auth.NewStdTx(stdMsg.Msgs, stdMsg.Fee, nil, stdMsg.Memo)

	//对交易签名
	signedTx, err := txBldr.SignStdTx(sign.FromName, sign.Password, tx, false)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	//交易签名后的字节
	txBytes, err := utils.GetTxEncoder(cliCtx.Codec)(signedTx)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	// broadcast to a Tendermint node
	return cliCtx.BroadcastTxAsync(txBytes)
}
