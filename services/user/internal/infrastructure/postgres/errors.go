package postgres

import (
	"database/sql"
	"strings"

	"github.com/lib/pq"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
)

// PostgreSQL Error Codes
const (
	PgErrUniqueViolation     = "23505" // 重複キー違反
	PgErrForeignKeyViolation = "23503" // 外部キー違反
	PgErrNotNullViolation    = "23502" // NOT NULL違反
	PgErrCheckViolation      = "23514" // CHECK制約違反
	PgErrInvalidTextValue    = "22P02" // 不正なテキスト表現
	PgErrDataException       = "22000" // データ例外
)

// handlePostgresError converts PostgreSQL errors to domain errors
func handlePostgresError(err error) error {
	if err == nil {
		return nil
	}

	// Handle no rows error
	if err == sql.ErrNoRows {
		return domain.NewError(domain.ErrNotFound.Error())
	}

	// Handle PostgreSQL specific errors
	if pgErr, ok := err.(*pq.Error); ok {
		switch string(pgErr.Code) {
		case PgErrUniqueViolation:
			// Check which field caused the violation
			if strings.Contains(pgErr.Detail, "email") || strings.Contains(pgErr.Constraint, "email") {
				return domain.ErrDuplicateEmail
			}
			if strings.Contains(pgErr.Detail, "id") || strings.Contains(pgErr.Constraint, "pkey") {
				return domain.ErrDuplicateID
			}
			if strings.Contains(pgErr.Detail, "name") || strings.Contains(pgErr.Constraint, "name") {
				return domain.NewError("name already exists")
			}
			return domain.ErrUserAlreadyExists

		case PgErrForeignKeyViolation:
			return domain.NewError("foreign key constraint violation")

		case PgErrNotNullViolation:
			// Check which column is null
			if strings.Contains(pgErr.Column, "email") {
				return domain.NewError("email is required")
			}
			if strings.Contains(pgErr.Column, "name") {
				return domain.NewError("name is required")
			}
			return domain.NewError("required field is missing")

		case PgErrCheckViolation:
			return domain.NewError("check constraint violation")

		case PgErrInvalidTextValue, PgErrDataException:
			return domain.NewError("invalid data format")

		default:
			// Return the original error for unknown error codes
			return err
		}
	}

	// Return the original error if it's not a PostgreSQL error
	return err
}
