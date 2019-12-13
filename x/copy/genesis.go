package muse

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"muse/x/copy/keeper"
	"muse/x/copy/types"
)

type GenesisState struct {
	LyricRecords []types.Lyric `json:"lyric_records"`
}

func NewGenesisState(lyricRecords []types.Lyric) GenesisState {
	return GenesisState{LyricRecords: nil}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.LyricRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid LyricRecords: Owner: %s. Error: Missing Owner", record.Owner)
		}
		if record.String() == "" {
			return fmt.Errorf("invalid LyricRecords: Value: %s. Error: Missing Value", record.String())
		}
		if len(record.Hash) <= 0 {
			return fmt.Errorf("invalid LyricRecords: Hash: %s. Error: Missing Hash", record.Hash)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		LyricRecords: []types.Lyric{},
	}
}

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.LyricRecords {
		keeper.SetLyric(ctx, record)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) GenesisState {
	var records []types.Lyric
	iterator := k.GetLyricsIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {

		lyricCode := string(iterator.Key())
		lyric := k.GetLyric(ctx, lyricCode)
		records = append(records, lyric)

	}
	return GenesisState{LyricRecords: records}
}
