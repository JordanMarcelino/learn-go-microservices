package encryptutils

type Hasher interface {
	Hash(password string) (string, error)
	Check(password, hash string) bool
}
