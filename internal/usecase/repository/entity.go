// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserRoles string

const (
	UserRolesAuthor UserRoles = "author"
	UserRolesReader UserRoles = "reader"
)

func (e *UserRoles) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserRoles(s)
	case string:
		*e = UserRoles(s)
	default:
		return fmt.Errorf("unsupported scan type for UserRoles: %T", src)
	}
	return nil
}

type NullUserRoles struct {
	UserRoles UserRoles `json:"userRoles"`
	Valid     bool      `json:"valid"` // Valid is true if UserRoles is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserRoles) Scan(value interface{}) error {
	if value == nil {
		ns.UserRoles, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserRoles.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserRoles) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserRoles), nil
}

func AllUserRolesValues() []UserRoles {
	return []UserRoles{
		UserRolesAuthor,
		UserRolesReader,
	}
}

type Blog struct {
	ID           uuid.UUID      `json:"id"`
	Descriptions sql.NullString `json:"descriptions"`
	UserRole     UserRoles      `json:"userRole"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    sql.NullTime   `json:"updatedAt"`
}