package util

type LangString struct {
	De string `json:"de"`
	En string `json:"en"`
}

func (ls LangString) Get() string {
	return ls.De
}
