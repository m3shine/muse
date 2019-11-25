package types

type QueryResLyric struct {
	Value string `json:"value"`
}

func (r QueryResLyric) String() string {
	return r.Value
}