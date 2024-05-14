package sqlDB

import (
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
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

func InitializeSqlDBTest() {
	basePath := test.MountAbsolutPath(test.DATABASE_ENVIRONMENT_PATH)

	test.InitializeSqlDBTest()
	pc := test.UsePostgresContainer()

	if err := pc.Dataset(basePath, "schema.sql"); err != nil {
		logging.Fatal(err.Error())
	}

	datasets := []string{"clear-database.sql", "add-users.sql", "add-contacts.sql", "add-dogs.sql"}
	pc.Dataset(basePath, datasets...)

	Initialize()
}
