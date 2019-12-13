package muse

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"muse/x/copy/keeper"
	"muse/x/copy/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSetLyric:
			return handleMsgSetLyric(ctx, k, msg)
		case types.MsgSetMusic:
			return handleMsgSetMusic(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized copyrights msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
/*
func NewAnteHandler(k keeper.Keeper) sdk.AnteHandler{
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, result sdk.Result, abort bool) {
		switch tx := tx.(type) {


		}
	}
}
*/
func handleMsgSetLyric(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSetLyric) sdk.Result {
	err := keeper.SetLyric(ctx, types.NewLyric(msg.BCICode, msg.Author, msg.Title, msg.Content, msg.Hash, msg.Owner, msg.Token))
	if err != nil {
		return sdk.ErrUnknownRequest(err.Error()).Result()
	}
	return sdk.Result{}
}

func handleMsgSetMusic(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSetMusic) sdk.Result {
	err := keeper.SetMusic(ctx, types.NewMusicWorks(msg.BCICode, msg.Title, msg.Hash, msg.Stakeholders, msg.WorksID))
	if err != nil {
		return sdk.ErrUnknownRequest(err.Error()).Result()
	}
	return sdk.Result{}
}