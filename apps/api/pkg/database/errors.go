package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// uniqueViolationCode is the PostgreSQL SQLSTATE for a unique constraint violation.
const uniqueViolationCode = "23505"

// IsUniqueViolation reports whether err is a PostgreSQL unique-constraint
// violation (e.g. inserting a duplicate username).
func IsUniqueViolation(err error) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code == uniqueViolationCode
	}
	return false
}
