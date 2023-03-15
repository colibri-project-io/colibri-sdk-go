package sqlDB

import (
	"context"
	"database/sql"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlTransactionWithoutInitialize(t *testing.T) {
	basePath := test.MountAbsolutPath("../../../development-environment/database/sql-tx/")

	test.InitializeSqlDBTest()
	pc := test.UsePostgresContainer()

	if err := pc.Dataset(basePath, "schema.sql"); err != nil {
		logging.Fatal(err.Error())
	}

	instance = nil

	stmt1 := NewStatement(context.Background(), "", "Contact Name 1", "em@il.com")
	err := stmt1.Execute()

	assert.Error(t, err, db_not_initialized_error)
}

func TestSqlTransaction(t *testing.T) {
	type contact struct {
		Name  string
		Email string
	}

	basePath := test.MountAbsolutPath("../../../development-environment/database/sql-tx/")

	test.InitializeSqlDBTest()
	pc := test.UsePostgresContainer()

	Initialize()

	if err := pc.Dataset(basePath, "schema.sql"); err != nil {
		logging.Fatal(err.Error())
	}

	t.Run("Should return error when query is nil", func(t *testing.T) {
		if err := pc.Dataset(basePath, "clear.sql"); err != nil {
			logging.Fatal(err.Error())
		}

		stmt := NewStatement(context.Background(), "", "Contact Name 1", "em@il.com")
		err := stmt.Execute()

		assert.Error(t, err, query_is_empty_error)
	})

	t.Run("commit ok", func(t *testing.T) {
		if err := pc.Dataset(basePath, "clear.sql"); err != nil {
			logging.Fatal(err.Error())
		}

		tx := NewTransaction()
		err := tx.ExecTx(context.Background(), func(ctx context.Context) error {
			insertContact1 := "INSERT INTO contacts (name, email) VALUES ($1, $2) "
			stmt1 := NewStatement(ctx, insertContact1, "Contact Name 1", "em@il.com")
			if err := stmt1.Execute(); err != nil {
				return err
			}

			insertContact2 := "INSERT INTO contacts (name, email) VALUES ($1, $2) "
			stmt2 := NewStatement(ctx, insertContact2, "Contact Name 2", "2em@il.com")
			if err := stmt2.Execute(); err != nil {
				return err
			}

			return nil
		})
		assert.NoError(t, err)

		q1 := NewQuery[contact](context.Background(), "SELECT name, email FROM contacts WHERE email = $1", "em@il.com")
		c1, err := q1.One()
		assert.NoError(t, err)
		assert.NotNil(t, c1)
		assert.Equal(t, "em@il.com", c1.Email)

		q2 := NewQuery[contact](context.Background(), "SELECT name, email FROM contacts WHERE email = $1", "2em@il.com")
		c2, err := q2.One()
		assert.NoError(t, err)
		assert.NotNil(t, c2)
		assert.Equal(t, "2em@il.com", c2.Email)
	})

	t.Run("fail with rollback", func(t *testing.T) {
		if err := pc.Dataset(basePath, "clear.sql"); err != nil {
			logging.Fatal(err.Error())
		}

		q1 := NewQuery[contact](context.Background(), "SELECT name, email FROM contacts WHERE email = $1", "em@il.com")
		c1, err := q1.One()
		assert.NoError(t, err)
		assert.Nil(t, c1)

		tx := NewTransaction()
		err = tx.ExecTx(context.Background(), func(ctx context.Context) error {
			insertContact1 := "INSERT INTO contacts (name, email) VALUES ($1, $2) "
			stmt1 := NewStatement(ctx, insertContact1, "Contact Name 1", "em@il.com")
			if err := stmt1.Execute(); err != nil {
				return err
			}

			insertContact2 := "INSERT INTO contacts (name, email) VALUES ($1, $2) "
			stmt2 := NewStatement(ctx, insertContact2, "Contact Name 2", "em@il.com")
			if err := stmt2.Execute(); err != nil {
				return err
			}

			return nil
		})
		assert.Error(t, err)

		q1 = NewQuery[contact](context.Background(), "SELECT name, email FROM contacts WHERE email = $1", "em@il.com")
		c1, err = q1.One()
		assert.NoError(t, err)
		assert.Nil(t, c1)
	})
}

func TestSqlTransaction_isolationLevel(t *testing.T) {
	t.Run("isolation level default", func(t *testing.T) {
		tx := NewTransaction()
		assert.NotNil(t, tx)
		assert.Equal(t, sql.LevelDefault, tx.(*sqlTransaction).isolation)
	})

	t.Run("isolation level serializable", func(t *testing.T) {
		tx := NewTransaction(sql.LevelSerializable)
		assert.NotNil(t, tx)
		assert.Equal(t, sql.LevelSerializable, tx.(*sqlTransaction).isolation)
	})

	t.Run("multiple isolations, only first is used", func(t *testing.T) {
		tx := NewTransaction(sql.LevelLinearizable, sql.LevelSerializable)
		assert.NotNil(t, tx)
		assert.Equal(t, sql.LevelLinearizable, tx.(*sqlTransaction).isolation)
	})
}
