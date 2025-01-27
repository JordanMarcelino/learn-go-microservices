package encryptutils

import "encoding/base64"

type Base64Encryptor struct {
}

func NewBase64Encryptor() *Base64Encryptor {
	return &Base64Encryptor{}
}

func (e *Base64Encryptor) Encrypt(data string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(data)), nil
}

func (e *Base64Encryptor) Decrypt(data string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
