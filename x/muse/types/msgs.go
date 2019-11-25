package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName

type MsgSetLyric struct {
	LyricCode string         `json:"lyric_code"`
	Author    string         `json:"author"`
	Title     string         `json:"title"`
	Hash      []byte         `json:"hash"`
	Owner     sdk.AccAddress `json:"owner"`
	Token     sdk.Coins      `json:"token"`
}

type MsgResetLyric struct {
	LyricCode string         `json:"lyric_code"`
	Author    string         `json:"author"`
	Title     string         `json:"title"`
	Hash      []byte         `json:"hash"`
	Owner     sdk.AccAddress `json:"owner"`
	Token     sdk.Coins      `json:"token"`
}

func NewMsgSetLyric(lyricCode string, author string, title string, hash []byte, owner sdk.AccAddress, token sdk.Coins) MsgSetLyric {
	return MsgSetLyric{
		LyricCode: lyricCode,
		Author:    author,
		Title:     title,
		Hash:      hash,
		Owner:     owner,
		Token:     token,
	}
}

func (msg MsgSetLyric) Route() string { return RouterKey }

func (msg MsgSetLyric) Type() string { return "set_lyric" }

func (msg MsgSetLyric) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	if len(msg.LyricCode) == 0 || len(msg.Hash) == 0 {
		return sdk.ErrUnknownRequest("LyricCode and/or Hash cannot be empty")
	}
	return nil
}

func (msg MsgSetLyric) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSetLyric) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}


func (msg MsgResetLyric) Route() string { return RouterKey }

func (msg MsgResetLyric) Type() string { return "reset_lyric" }

func (msg MsgResetLyric) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	if len(msg.LyricCode) == 0 || len(msg.Hash) == 0 {
		return sdk.ErrUnknownRequest("LyricCode and/or Hash cannot be empty")
	}
	return nil
}

func (msg MsgResetLyric) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgResetLyric) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
