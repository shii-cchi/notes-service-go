package hash

type Hasher interface {
	Hash(data string) (string, error)
	IsValidData(hashedData, data string) bool
}
