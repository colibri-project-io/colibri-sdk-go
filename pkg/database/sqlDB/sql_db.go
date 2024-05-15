package sqlDB

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/lib/pq"

	"io"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/observer"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"golang.org/x/exp/slices"
)

const (
	db_connection_success    string = "%s database connected"
	db_connection_error      string = "An error occurred while trying to connect to the %s database. Error: %s"
	db_migration_error       string = "An error occurred when validate database migrations: %v"
	db_not_initialized_error string = "database not initialized"
	query_is_empty_error     string = "query is empty"
	page_is_empty_error      string = "page is empty"
)

// sqlDBObserver is a struct for sql database observer.
type sqlDBObserver struct {
	name     string
	instance *sql.DB
}

// sqlDBInstance is a pointer to sql.DB
var sqlDBInstance *sql.DB

// Initialize start connection with sql database and execute migration.
//
// No parameters.
// No return values.
func Initialize() {
	sqlDB := NewSQLDatabaseInstance("SQL", config.SQL_DB_CONNECTION_URI)
	sqlDB.SetMaxOpenConns(config.SQL_DB_MAX_OPEN_CONNS)
	sqlDB.SetMaxIdleConns(config.SQL_DB_MAX_IDLE_CONNS)

	if err := executeDatabaseMigration(sqlDB); err != nil {
		logging.Fatal(db_migration_error, err)
	}

	sqlDBInstance = sqlDB
}

// NewSQLDatabaseInstance creates a new SQL database instance.
//
// Parameters:
// - name: a string representing the name of the database.
// - databaseURL: a string representing the URL of the database.
// Returns a pointer to sql.DB.
func NewSQLDatabaseInstance(name, databaseURL string) *sql.DB {
	sqlDB, err := sql.Open(monitoring.GetSQLDBDriverName(), databaseURL)
	if err != nil {
		logging.Fatal(db_connection_error, name, err)
	}

	if err = sqlDB.Ping(); err != nil {
		logging.Fatal(db_connection_error, name, err)
	}

	logging.Info(db_connection_success, name)
	observer.Attach(sqlDBObserver{name, sqlDB})
	return sqlDB
}

// Close finalize sql database connection
//
// No parameters.
// No return values.
func (o sqlDBObserver) Close() {
	logging.Info("waiting to safely close the %s database connection", o.name)
	if observer.WaitRunningTimeout() {
		logging.Warn("WaitGroup timed out, forcing close the %s database connection", o.name)
	}
	logging.Info("closing %s database connection", o.name)
	if err := o.instance.Close(); err != nil {
		logging.Error("error when closing %s database connection: %+v", o.name, err)
	}
}

// getDataList retrieves a list of items from the given sql.Rows object.
//
// It takes a sql.Rows object as input and returns a list of items of type T and an error.
func getDataList[T any](rows *sql.Rows) ([]T, error) {
	list := make([]T, 0)
	for rows.Next() {
		model := new(T)
		err := rows.Scan(reflectCols(model)...)
		if err != nil {
			return nil, err
		}

		list = append(list, *model)
	}

	return list, nil
}

// reflectCols generates a list of column values from the provided model.
//
// model: the model to reflect columns from
// []any: a list of column values
func reflectCols(model any) (cols []any) {
	typeOf := reflect.TypeOf(model).Elem()
	valueOf := reflect.ValueOf(model).Elem()

	isStruct, isTime, isNull, isSlice := reflectValueValidations(valueOf)
	if isSlice {
		cols = append(cols, pq.Array(valueOf.Addr().Interface()))
	} else if !isStruct || isTime || isNull {
		cols = append(cols, valueOf.Addr().Interface())
		return
	}

	for i := 0; i < typeOf.NumField(); i++ {
		field := valueOf.Field(i)

		isStruct, isTime, isNull, isSlice = reflectValueValidations(field)
		if isSlice {
			cols = append(cols, pq.Array(field.Addr().Interface()))
		} else if isStruct && !isTime && !isNull {
			cols = append(cols, reflectCols(field.Addr().Interface())...)
		} else {
			cols = append(cols, field.Addr().Interface())
		}
	}

	return cols
}

// reflectValueValidations validates the type of the provided value.
//
// value: the value to validate
// (isStruct, isTime, isNull, isSlice) : returns booleans indicating if the value is a struct, time type, null type, or a slice.
func reflectValueValidations(value reflect.Value) (isStruct, isTime, isNull, isSlice bool) {
	isStruct = value.Kind() == reflect.Struct
	isTime = slices.Contains([]string{"time.Time", "types.IsoDate", "types.IsoTime"}, value.Type().String())
	isNull = strings.Contains(value.Type().String(), "Null")
	isSlice = value.Kind() == reflect.Slice
	return
}

// closer closes the provided io.Closer interface and logs an error if closing fails.
//
// o: the io.Closer interface to be closed
// Error: returns any error encountered during closing.
func closer(o io.Closer) {
	if err := o.Close(); err != nil {
		logging.Error("could not close statement: %v", err)
	}
}
