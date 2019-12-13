package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
)

const TokenAmount int64 = 10e+7

//歌词
type Lyric struct {
	BCICode string         `json:"bci_code"`
	Author  string         `json:"author"`
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Hash    []byte         `json:"hash"`
	Owner   sdk.AccAddress `json:"owner"`
	Token   sdk.Coins      `json:"token"`
}

func NewLyric(bciCode, author, title, content string, hash []byte, owner sdk.AccAddress, token sdk.Coins) Lyric {
	return Lyric{
		BCICode: bciCode,
		Author:  author,
		Title:   title,
		Content: content,
		Hash:    hash,
		Owner:   owner,
		Token:   token,
	}
}

func (l Lyric) String() string {
	return fmt.Sprintf(`"bci_code":"%s","author":"%s","title":"%s","content":"%s","hash":"%s","owner":"%s","token":"%v"`, l.BCICode, l.Author, l.Title, l.Content, l.Hash, l.Owner, l.Token)
}

//音乐作品
type MusicWorks struct {
	BCICode      string        `json:"bci_code"`
	Title        string        `json:"title"`
	Hash         []byte        `json:"hash"`
	LyricCode    string        `json:"lyric_code,omitempty"`
	Stakeholders []Stakeholder `json:"stakeholders"`
	WorksID      string        `json:"works_id,omitempty"` //ISWC、DCI等
}

type Stakeholder struct {
	Name     string         `json:"name"`
	IDI      string         `json:"idi,omitempty"`  //现实身份标识
	Type     string         `json:"type,omitempty"` //词曲作者、出版商、发行人、演出团体
	Address  sdk.AccAddress `json:"address"`
	Describe string         `json:"describe,omitempty"`
	Weights  big.Float      `json:"weights"` //权重
}

func NewMusicWorks(bciCode string, title string, hash []byte, holders []Stakeholder, worksId string) MusicWorks {
	return MusicWorks{
		BCICode:      bciCode,
		Title:        title,
		Hash:         hash,
		Stakeholders: holders,
		WorksID:      worksId,
	}
}

func (m MusicWorks) String() string {
	return fmt.Sprintf(`"bci_code":"%s","title":"%s","hash":"%s","stakeholders":"%v","works_id":"%s"`, m.BCICode, m.Title, m.Hash, m.Stakeholders, m.WorksID)
}
