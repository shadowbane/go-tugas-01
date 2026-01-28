package helpers

import (
	cryptorand "crypto/rand"

	"github.com/oklog/ulid/v2"
)

func NewULID() string {
	return ulid.MustNew(ulid.Now(), cryptorand.Reader).String()
}
