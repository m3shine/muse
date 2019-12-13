package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"math/big"
	"regexp"
	"time"
)

const (
	RouterKey  = ModuleName
	BCIExp     = `^(BCI)\d{4}[A-Z]\d+\S{39}$`
	CtypeLyric = 'L'
	CtypeMusic = 'M'
)

type MsgSetLyric struct {
	BCICode string         `json:"bci_code"`
	Author  string         `json:"author"`
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Hash    []byte         `json:"hash"`
	Owner   sdk.AccAddress `json:"owner"`
	Token   sdk.Coins      `json:"token"`
}

func NewMsgSetLyric(bciCode, author, title, content string, hash []byte, owner sdk.AccAddress, token sdk.Coins) MsgSetLyric {
	return MsgSetLyric{
		BCICode: bciCode,
		Author:  author,
		Title:   title,
		Content: content,
		Hash:    hash,
		Owner:   owner,
		Token:   token,
	}
}

func (msg MsgSetLyric) Route() string { return RouterKey }

func (msg MsgSetLyric) Type() string { return "set_lyric" }

//执行交易时会检查
func (msg MsgSetLyric) ValidateBasic() sdk.Error {
	ok, err := ValidateBCI(msg.BCICode)
	if !ok || err != nil {
		return sdk.ErrUnknownRequest("invalid BCI code " + msg.BCICode)
	}
	b := []byte(msg.BCICode)
	if rune(b[7]) != msg.BCIType() {
		return sdk.ErrUnknownRequest("invalid BCI code "  + msg.BCICode)
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Content) == 0 {
		return sdk.ErrUnknownRequest("lyric content cannot be empty")
	}

	if len(msg.Hash) == 0 {
		return sdk.ErrUnknownRequest("BCICode and/or Hash cannot be empty")
	}
	return nil
}

func (msg MsgSetLyric) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSetLyric) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func (msg MsgSetLyric) BCIType() rune {
	return CtypeLyric
}

/*********************************************/

type MsgSetMusic struct {
	BCICode      string        `json:"bci_code"`
	Title        string        `json:"title"`
	Hash         []byte        `json:"hash"`
	LyricCode    string        `json:"lyric_code"`
	Stakeholders []Stakeholder `json:"stakeholders"`
	WorksID      string        `json:"works_id"` //ISWC、DCI等
}

func NewMsgSetMusic(bciCode, title string, hash []byte, lyricCode string,stakeholders []Stakeholder , worksId string) MsgSetMusic {
	return MsgSetMusic{
		BCICode: bciCode,
		Title:   title,
		Hash:    hash,
		LyricCode: lyricCode,
		Stakeholders: stakeholders,
		WorksID: worksId,
	}
}

func (msg MsgSetMusic) Route() string { return RouterKey }

func (msg MsgSetMusic) Type() string { return "set_music" }

func (msg MsgSetMusic) ValidateBasic() sdk.Error {
	ok, err := ValidateBCI(msg.BCICode)
	if !ok || err != nil {
		return sdk.ErrUnknownRequest("invalid BCI code " + msg.BCICode)
	}

	b := []byte(msg.BCICode)
	if rune(b[7]) != msg.BCIType() {
		return sdk.ErrUnknownRequest("invalid BCI code "  + msg.BCICode)
	}

	if len(msg.Stakeholders) <= 0 {
		return sdk.ErrInternal("no owners")
	} else {
		var sum, hundred big.Float
		sum.SetFloat64(0)
		hundred.SetFloat64(1)
		for _, s := range msg.Stakeholders {
			if s.Address.Empty() {
				return sdk.ErrInvalidAddress(s.Address.String())
			}
			sum.Add(&sum, &s.Weights)
		}
		if sum.Cmp(&hundred) != 0 {
			return sdk.ErrUnknownRequest("invalid weights")
		}
	}

	if len(msg.Hash) == 0 {
		return sdk.ErrUnknownRequest("BCICode and/or Hash cannot be empty")
	}

	return nil
}

func (msg MsgSetMusic) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSetMusic) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.Stakeholders))
	for _, s := range msg.Stakeholders {
		addrs = append(addrs, s.Address)
	}
	return addrs
}

func (msg MsgSetMusic) BCIType() rune {
	return CtypeMusic
}

/*****/

func ValidateBCI(code string) (bool, error) {
	return regexp.MatchString(BCIExp, code)
}

func GenBCI(cliCtx context.CLIContext, from string, ctype rune) (string, error) {
	addr, err := sdk.AccAddressFromBech32(from)
	_, seq, err := authtypes.NewAccountRetriever(cliCtx).GetAccountNumberSequence(addr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("BCI%d%c%d%s", time.Now().Year(), ctype, seq, []byte(from)[6:]), nil
}
