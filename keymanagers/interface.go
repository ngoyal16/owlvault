package keymanagers

// KeyManagementService defines methods for encryption and decryption.
type KeyManagementService interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}
