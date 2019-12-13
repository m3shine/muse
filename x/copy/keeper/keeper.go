package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"muse/x/copy/types"
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

func (k Keeper)GetLyric(ctx sdk.Context, bciCode string) types.Lyric {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(bciCode)) {
		return types.Lyric{}
	}
	bz := store.Get([]byte(bciCode))
	var ly types.Lyric
	k.cdc.MustUnmarshalBinaryBare(bz, &ly)
	return ly
}

func (k Keeper) GetLyricsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

func (k Keeper) SetLyric(ctx sdk.Context, lyric types.Lyric) error {
	store := ctx.KVStore(k.storeKey)
	if store.Has([]byte(lyric.BCICode)) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("the same bci code already exists: %s", lyric.BCICode))
	}
	store.Set([]byte(lyric.BCICode), k.cdc.MustMarshalBinaryBare(lyric))
	return nil
}

func (k Keeper) QueryLyric(ctx sdk.Context, bciCode string) string {
	return k.GetLyric(ctx, bciCode).String()
}


func (k Keeper)GetMusic(ctx sdk.Context, bciCode string) types.MusicWorks {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(bciCode)) {
		return types.MusicWorks{}
	}
	bz := store.Get([]byte(bciCode))
	var mw types.MusicWorks
	k.cdc.MustUnmarshalBinaryBare(bz, &mw)
	return mw
}

func (k Keeper) GetMusicsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

func (k Keeper) SetMusic(ctx sdk.Context, mw types.MusicWorks) error {
	store := ctx.KVStore(k.storeKey)
	if store.Has([]byte(mw.BCICode)) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("the same bci code already exists: %s", mw.BCICode))
	}
	store.Set([]byte(mw.BCICode), k.cdc.MustMarshalBinaryBare(mw))
	return nil
}

func (k Keeper) QueryMusic(ctx sdk.Context, bciCode string) string {
	return k.GetMusic(ctx, bciCode).String()
}
