package util

type LangString struct {
	De string
	En string
}

func (ls LangString) Get() string {
	return ls.De
}
