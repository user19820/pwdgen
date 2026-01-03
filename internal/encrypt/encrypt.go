package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

const keySize = 32

func Setup() error {
	homePath, homePathErr := os.UserHomeDir()
	if homePathErr != nil {
		return homePathErr
	}

	keyPath := filepath.Join(homePath, ".local", "share", "pwdgen", "pwdgen.key")

	_, keyPathOpenErr := os.Open(keyPath)
	if keyPathOpenErr == nil || !errors.Is(keyPathOpenErr, os.ErrNotExist) {
		return errors.New("pwdgen is already initialized")
	}

	key := make([]byte, keySize)

	if _, randReadErr := rand.Read(key); randReadErr != nil {
		return randReadErr
	}

	hexKey := hex.EncodeToString(key)

	if writeFileErr := os.WriteFile(keyPath, []byte(hexKey), 0o600); writeFileErr != nil {
		return writeFileErr
	}

	return nil
}

func Encrypt(data []byte) ([]byte, error) {
	if data == nil {
		return nil, errors.New("encrypt: data should not be nil")
	}

	hexKey, loadKeyErr := loadKey()
	if loadKeyErr != nil {
		return nil, fmt.Errorf("encrypt: load key: %w", loadKeyErr)
	}

	key, decodeHexErr := hex.DecodeString(string(hexKey))
	if decodeHexErr != nil {
		return nil, fmt.Errorf("encrypt: decode key: %w", decodeHexErr)
	}

	block, cipherErr := aes.NewCipher(key)
	if cipherErr != nil {
		return nil, fmt.Errorf("encrypt: cipher creation: %w", cipherErr)
	}

	gcm, gcmErr := cipher.NewGCM(block)
	if gcmErr != nil {
		return nil, fmt.Errorf("encrypt: gcm: %w", gcmErr)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, readFullErr := io.ReadFull(rand.Reader, nonce); readFullErr != nil {
		return nil, fmt.Errorf("encrypt: nonce generation: %w", readFullErr)
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func Decrypt(data []byte) ([]byte, error) {
	hexKey, loadKeyErr := loadKey()
	if loadKeyErr != nil {
		return nil, fmt.Errorf("encrypt: load key: %w", loadKeyErr)
	}

	key, decodeHexErr := hex.DecodeString(string(hexKey))
	if decodeHexErr != nil {
		return nil, fmt.Errorf("encrypt: decode key: %w", decodeHexErr)
	}

	block, cipherErr := aes.NewCipher(key)
	if cipherErr != nil {
		return nil, fmt.Errorf("decrypt: cipher creation: %w", cipherErr)
	}

	gcm, gcmErr := cipher.NewGCM(block)
	if gcmErr != nil {
		return nil, fmt.Errorf("decrypt: gcm: %w", gcmErr)
	}

	nonceSize := gcm.NonceSize()

	if len(data) < nonceSize {
		return nil, errors.New("decrypt: data has incorrect length")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}

func loadKey() ([]byte, error) {
	homeDir, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil {
		return nil, homeDirErr
	}

	keyDir := path.Join(homeDir, ".local", "share", "pwdgen", "pwdgen.key")

	f, fOpenErr := os.Open(keyDir)
	if fOpenErr != nil {
		return nil, fOpenErr
	}
	defer f.Close()

	key, readAllErr := io.ReadAll(f)
	if readAllErr != nil {
		return nil, readAllErr
	}

	return key, nil
}
