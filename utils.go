// Package gosocketify is a package with server side implementation for socketify
package gosocketify

import (
	"github.com/google/uuid"
)

func generateConnectionID() string {
	return uuid.New().String()
}
