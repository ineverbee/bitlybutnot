package store

//Store interface describes main functionality of storage
type Store interface {
	Set(key uint32, shortURL, longURL string) error
	Get(key uint32) (string, error)
}
