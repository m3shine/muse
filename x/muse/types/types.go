package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

const TokenAmount int64 = 10e+7

//歌词
type Lyric struct {
	LyricCode string         `json:"lyric_code"`
	Author    string         `json:"author"`
	Title     string         `json:"title"`
	Hash      []byte         `json:"hash"`
	Owner     sdk.AccAddress `json:"owner"`
	Token     sdk.Coins      `json:"token"`
}

//歌曲
type Tuen struct {
}

//音乐作品
type ISWCWorks struct {
	Lyric
	Tuen
}

//音像作品
type ISRCWorks struct {
	ISWCWorks
}

func NewLyric(lyricCode string, author string, title string, hash []byte, owner sdk.AccAddress, token sdk.Coins) Lyric {
	return Lyric{
		LyricCode: lyricCode,
		Author: author,
		Title: title,
		Hash: hash,
		Owner: owner,
		Token: token,
	}
}

func (l Lyric) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Author: %s
Title %s
Hash %s
Owner: %s
Token %s`, l.Author, l.Title, l.Hash, l.Owner, l.Token))
}
