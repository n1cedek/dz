package model

import (
	"database/sql"
	"time"
)

type Role struct {
	Value string
}

var (
	RoleUser  = Role{Value: "user"}
	RoleAdmin = Role{Value: "admin"}
)
func (r Role) String() string {
	return r.Value
}

type User struct {
	ID        int64        `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Password  string       `db:"password"`
	Role      Role         `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type PublicInfo struct {
	ID        int64        `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
}