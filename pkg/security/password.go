package security

import (
	"bytes"
	"crypto/rand"

	"api.cloud.io/pkg/errors"
	"golang.org/x/crypto/scrypt"
)

const (
	// DefaultSaltLen of SecurePassword.
	DefaultSaltLen = 32
	// DefaultKeyLen of SecurePassword.
	DefaultKeyLen = 32
)

// SecurePassword get password and secure it with scrypt key derivation function.
func SecurePassword(password []byte, saltLen, keyLen int) (key, salt []byte, err error) {
	salt = make([]byte, saltLen)
	_, err = rand.Read(salt)
	if err != nil {
		return nil, nil, errors.ErrInternal(err)
	}

	key, err = scrypt.Key(password, salt, 1<<15, 8, 1, keyLen)
	if err != nil {
		return nil, nil, errors.ErrInternal(err)
	}

	return
}

// VerifyPassword compares key of password with key.
func VerifyPassword(password, key, salt []byte) error {
	newKey, err := scrypt.Key(password, salt, 1<<15, 8, 1, len(key))
	if err != nil {
		return errors.ErrInternal(err)
	}

	if !bytes.Equal(key, newKey) {
		return errors.ErrWrongPassword
	}

	return nil
}
