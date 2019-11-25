package muse

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"muse/x/muse/keeper"
	"muse/x/muse/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSetLyric:
			return handleMsgSetLyric(ctx, k, msg)
		case types.MsgResetLyric:
			return handleMsgReSetLyric(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized muse Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSetLyric(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSetLyric) sdk.Result {
	if msg.Owner.Equals(keeper.GetLyric(ctx, msg.LyricCode).Owner) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}
	err := keeper.SetLyric(ctx, msg.LyricCode, types.NewLyric(msg.LyricCode, msg.Author, msg.Title, msg.Hash, msg.Owner, msg.Token))
	if err != nil {
		return sdk.ErrUnknownRequest(err.Error()).Result()
	}
	return sdk.Result{}
}

func handleMsgReSetLyric(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgResetLyric) sdk.Result {
	if msg.Owner.Equals(keeper.GetLyric(ctx, msg.LyricCode).Owner) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}
	err := keeper.ResetLyric(ctx, msg.LyricCode, types.NewLyric(msg.LyricCode, msg.Author, msg.Title, msg.Hash, msg.Owner, msg.Token))
	if err != nil {
		return sdk.ErrUnknownRequest(err.Error()).Result()
	}
	return sdk.Result{}
}
