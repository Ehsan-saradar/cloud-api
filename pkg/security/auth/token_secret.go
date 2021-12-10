package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"api.cloud.io/pkg/errors"
	"github.com/satori/go.uuid"
)

// GetSecret returns refresh token secret based on the following information.
func GetSecret(id uuid.UUID, ownerID uuid.UUID,username string, issuedAt time.Time, applicantIP, userAgent string) (string, error) {
	hash := sha256.New()

	// SHA256(ID + OwnerID + IssuedAt + ApplicantIP + UserAgent + 32 random bytes)
	hash.Write(id.Bytes())
	hash.Write(ownerID.Bytes())
	hash.Write([]byte(username))
	issuedAtBytes, _ := issuedAt.MarshalBinary()
	hash.Write(issuedAtBytes)
	hash.Write([]byte(applicantIP))
	hash.Write([]byte(userAgent))
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", errors.ErrInternal(err)
	}
	hash.Write(randomBytes)

	return hex.EncodeToString(hash.Sum(nil)), nil
}
