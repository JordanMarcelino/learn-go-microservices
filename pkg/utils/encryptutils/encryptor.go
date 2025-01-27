package encryptutils

type Encryptor interface {
	Encrypt(data string) (string, error)
	Decrypt(data string) (string, error)
}
