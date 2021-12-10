package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

// Session is the token model.
type Session struct {
	ID                 uuid.UUID `sql:"type:uuid"`
	OwnerID            uuid.UUID `sql:"type:uuid,notnull"`
	Username           string
	OwnerType          string `sql:",notnull"`
	IssuedAt           time.Time
	ExpirationTime     time.Time
	Scopes             []string `pg:",array"`
	ApplicantIP        string
	UserAgent          string
	PushNotificationID string
	Secret             string `sql:"type:varchar(64),unique"` // Hex(SHA512(ID + OwnerID + CreatedAt + ApplicantIP + UserAgent + 32 random bytes))
}
