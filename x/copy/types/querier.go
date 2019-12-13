package types

type QueryResLyric struct {
	Value string `json:"value"`
}

func (r QueryResLyric) String() string {
	return r.Value
}

type QueryResMusic struct {
	Value string `json:"value"`
}

func (r QueryResMusic) String() string {
	return r.Value
}