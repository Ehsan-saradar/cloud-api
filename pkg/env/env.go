package env

import (
	"log"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// GetBool retrieves the value of the environment variable named by the key in bool type.
func GetBool(key string) bool {
	str := os.Getenv(key)
	b, err := strconv.ParseBool(str)
	if err != nil {
		log.Fatalln(errors.Wrapf(err, "could not parse environment %s", key))
	}
	return b
}
