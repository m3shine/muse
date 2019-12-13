package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/lib/types"
	"net/http"
	"time"
)

const (
	bciCode                       = "bci_code"
	rpcid   rpctypes.JSONRPCIntID = 112
)

var httpClient = http.Client{Timeout: 10 * time.Second}

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/lyric/{%s}", storeName, bciCode), lyricHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/lyric", storeName), setLyricHandler(cliCtx)).Methods("PUT")


	r.HandleFunc(fmt.Sprintf("/music/{%s}", storeName, bciCode), musicHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/music", storeName), setMusicHandler(cliCtx)).Methods("PUT")
}