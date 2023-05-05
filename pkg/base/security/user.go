package security

import "encoding/base64"

const (
	profileField  string = "profile"
	tenantIdField string = "tenantId"
)

type User struct {
	ID       string
	Email    string
	Name     string
	TenantID string
	Profile  string
}

type UserCreate struct {
	ID       string
	Email    string
	Password string
	Name     string
	TenantID string
	Profile  string
}

type UserUpdate struct {
	Email    string
	Password string
	Name     string
	Profile  string
}

func (m *UserCreate) decodeBase64Password() (string, error) {
	decodedPassword, err := base64.StdEncoding.DecodeString(m.Password)
	if err != nil {
		return "", err
	}

	return string(decodedPassword), nil
}

func (m *UserUpdate) decodeBase64Password() (string, error) {
	decodedPassword, err := base64.StdEncoding.DecodeString(m.Password)
	if err != nil {
		return "", err
	}

	return string(decodedPassword), nil
}
