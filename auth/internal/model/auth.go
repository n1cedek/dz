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

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Role      Role
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type PublicInfo struct {
	ID    int64
	Name  string
	Email string
}
