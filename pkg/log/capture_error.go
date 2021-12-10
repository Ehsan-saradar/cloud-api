package log

import (
	"fmt"

	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
)

// CaptureError send err with some data to sentry.
func CaptureError(err error) {
	cause := errors.Cause(err)

	inters := []raven.Interface{
		raven.NewException(cause, raven.GetOrNewStacktrace(cause, 1, 3, nil)),
	}
	packet := raven.NewPacket(err.Error(), inters...)
	eventID, ch := raven.Capture(packet, nil)
	if eventID != "" {
		sentryErr := <-ch
		if sentryErr != nil {
			fmt.Printf("could not send (%s) error to sentry: %s\n", err, sentryErr)
		}
	}
}
