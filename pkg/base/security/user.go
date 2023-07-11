package security

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
	TenantID string
	Profile  string
}
