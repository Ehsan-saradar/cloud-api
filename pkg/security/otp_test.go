package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhoneOTP(t *testing.T) {
	phone := "12345678911"
	pass := GenerateNumericOTP(1000, 9999)

	_, err := SetPhoneOTP(phone, pass)
	assert.NoError(t, err)

	err = VerifyPhoneOTP(phone, pass)
	assert.NoError(t, err)
}

func TestPhoneToken(t *testing.T) {
	phone := "12345678911"
	tok, err := GeneratePhoneToken()
	assert.NoError(t, err)

	SetPhoneToken(phone, tok)

	err = VerifyPhoneToken(phone, tok)
	assert.NoError(t, err)
}
