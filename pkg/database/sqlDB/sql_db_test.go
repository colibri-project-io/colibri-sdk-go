package sqlDB

import (
	"time"
)

const (
	query_base = "SELECT u.id, u.name, u.birthday, p.id, p.name FROM users u JOIN profiles p ON u.profile_id = p.id"
)

type Profile struct {
	Id   int
	Name string
}

type User struct {
	Id       int
	Name     string
	Birthday time.Time
	Profile  Profile
}

type Dog struct {
	ID              uint
	Name            string
	Characteristics []string
}
