package dotastats

import (
	"errors"
	"time"
)

// Errors
var ErrNoRows = errors.New("db: no rows in result set")
var ErrDuplicateRow = errors.New("db: duplicate row found for unique constraint")

func TimeNow() time.Time {
	return time.Now().UTC()
}
