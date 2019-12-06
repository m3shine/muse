package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
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