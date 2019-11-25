package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"muse/x/muse/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey  sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey: storeKey,
		cdc: cdc,
	}
}

func (k Keeper) setLyric(ctx sdk.Context, lyricCode string, lyric types.Lyric){
	if lyric.Owner.Empty() || len(lyric.Hash) <= 0 {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(lyricCode), k.cdc.MustMarshalBinaryBare(lyric))
}

func (k Keeper)GetLyric(ctx sdk.Context, lyricCode string) types.Lyric {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(lyricCode)) {
		return types.Lyric{}
	}
	bz := store.Get([]byte(lyricCode))
	var ly types.Lyric
	k.cdc.MustUnmarshalBinaryBare(bz, &ly)
	return ly
}

func (k Keeper) GetLyricsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

func (k Keeper) SetLyric(ctx sdk.Context, lyricCode string, lyric types.Lyric) error {
	store := ctx.KVStore(k.storeKey)
	if store.Has([]byte(lyricCode)) {
		return sdk.ErrUnknownRequest("do not repeat setting lyrics")
	}
	//lyric.Token = sdk.Coins{sdk.NewInt64Coin(tokenName, types.TokenAmount)}
	k.setLyric(ctx, lyricCode,lyric)
	return nil
}

func (k Keeper) ResetLyric(ctx sdk.Context, lyricCode string, lyric types.Lyric) error{
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(lyricCode)) {
		return sdk.ErrUnknownRequest("the lyricCode you want to reset do not exist")
	}
	bz := store.Get([]byte(lyricCode))
	var newLyric types.Lyric
	k.cdc.MustUnmarshalBinaryBare(bz, &newLyric)
	newLyric = lyric

	k.setLyric(ctx, lyricCode, newLyric)
	return nil
}

func (k Keeper) QueryLyric(ctx sdk.Context, lyricCode string) string {
	return k.GetLyric(ctx, lyricCode).String()
}