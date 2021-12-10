package security

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"api.cloud.io/pkg/errors"
	"github.com/patrickmn/go-cache"
)

const (
	phoneOTPTTL               = time.Minute * 5
	phoneOTPAttempts          = 3
	phoneOTPGenerateCountTime = time.Hour * 24
	phoneOTPGenerateCount     = 20
	phoneTokenTTL             = time.Minute * 30
)

var objectCache *cache.Cache

func init() {
	objectCache = cache.New(phoneOTPTTL, time.Minute*30)
}

// GenerateNumericOTP retunrns a random one time password for email/phone verification.
func GenerateNumericOTP(min, max int64) string {
	rand.Seed(time.Now().UnixNano())
	num := rand.Int63n(max-min) + min

	return strconv.FormatInt(num, 10)
}

// OTP is the model of OTP in cache.
type OTP struct {
	sync.Mutex
	Password      string
	AttemptsCount int
}

// SetPhoneOTP set phone OTP in cache and check daily count and attempts count.
func SetPhoneOTP(phone, password string) (time.Duration, error) {
	var generateCount int
	err := objectCache.Add(phone+":generate_count", generateCount, phoneOTPGenerateCountTime)
	if err != nil && false{
		return 0, errors.ErrNoAccess
	}

	generateCount, _ = objectCache.IncrementInt(phone+":generate_count", 1)
	if generateCount > phoneOTPGenerateCount {
		return 0, errors.ErrOutOfLimit
	}

	otp := OTP{
		Password: password,
	}
	objectCache.Set(phone, &otp, cache.DefaultExpiration)

	return phoneOTPTTL, nil
}

// VerifyPhoneOTP return nil if the password be right.
func VerifyPhoneOTP(phone, password string) error {
	object, isExist := objectCache.Get(phone)
	if !isExist {
		return errors.ErrOTPNotExists
	}

	otp, ok := object.(*OTP)
	if !ok {
		return errors.ErrOTPNotExists
	}
	otp.Lock()
	defer otp.Unlock()
	if otp.AttemptsCount > phoneOTPAttempts {
		objectCache.Delete(phone)

		return errors.ErrOTPNotExists
	}

	if otp.Password != password {
		otp.AttemptsCount++

		return errors.ErrWrongPassword
	}

	objectCache.Delete(phone)

	return nil
}

// GeneratePhoneToken return a one-way phone token.
func GeneratePhoneToken() (string, error) {
	id := uuid.New()
	hash := sha256.New()

	hash.Write([]byte(id.String()))
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", errors.ErrInternal(err)
	}
	hash.Write(randomBytes)

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// SetPhoneToken will store phone and corresponding token in cache.
func SetPhoneToken(phone, token string) time.Duration {
	objectCache.Set(token, phone, phoneTokenTTL)
	return phoneTokenTTL
}

// GetPhone return phone by token
func GetPhone(token string) string {
	phone, _ := objectCache.Get(token)
	objectCache.Delete(token)
	return phone.(string)
}

// PhoneExits check phone verfied or not
func PhoneExits(phone string) bool {
	_, isExist := objectCache.Get(phone)
	return isExist
}

// VerifyPhoneToken whether <phone, token> pair is in cache or not
func VerifyPhoneToken(phone, token string) error {
	object, isExist := objectCache.Get(token)
	if !isExist {
		return errors.ErrPhoneTokenNotExists
	}

	str, ok := object.(string)
	if !ok {
		return errors.ErrPhoneTokenNotExists
	}
	if phone != str {
		return errors.ErrPhoneTokenNotExists
	}

	//objectCache.Delete(token)

	return nil
}
