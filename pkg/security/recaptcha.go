package security

import (
	"api.cloud.io/pkg/errors"
	"api.cloud.io/pkg/log"
	"gopkg.in/ezzarghili/recaptcha-go.v2"
)

var (
	RecaptchaSecret string
	recaptchaClient recaptcha.ReCAPTCHA
)

func init() {
	recaptchaClient, _ = recaptcha.NewReCAPTCHA(RecaptchaSecret)
}

func VerifyRecaptcha(token, clientIP string) error {
	success, err := recaptchaClient.Verify(token, clientIP)
	if err != nil {
		log.CaptureError(err)
	}
	if !success {
		return errors.ErrRecaptchaFailed
	}
	return nil
}
