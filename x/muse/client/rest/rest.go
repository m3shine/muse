package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	restLyricCode = "lyricCode"
)

var httpClient = http.Client{Timeout: 10 * time.Second}

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/lyric/{%s}", storeName, restLyricCode), lyricHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/lyric", storeName), setLyricHandler(cliCtx)).Methods("PUT")
}

func sendTendRequest(method string, params interface{}) ([]byte, error){
	//连接muse区块链节点的服务器app key，由节点管理后台生成
	//appKey := config.Conf.GetValue("eth", "app_key")
	url := "http://localhost:26657/"
	jsonRpcId := "whocare"

	args := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      jsonRpcId,
		"method":  method,
		"params":  params,
	}

	body, err := json.Marshal(args)
	if err != nil {
		return []byte{}, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))

	return do(req)
}

func do(r *http.Request) ([]byte, error) {
	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}
