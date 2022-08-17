package encrypt

type Encryptor interface {
	Encrypt(string) (string, error)
	Check(realString, encryptedString string) (bool, error)
}
