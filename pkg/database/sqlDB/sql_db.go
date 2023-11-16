package sqlDB

import (
	"database/sql"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/lib/pq"
	"reflect"
	"strings"

	"io"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/observer"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"golang.org/x/exp/slices"
)

const (
	db_connection_success    string = "SQL database connected"
	db_connection_error      string = "An error occurred while trying to connect to the database. Error: %s"
	db_migration_error       string = "An error occurred when validate database migrations: %v"
	db_not_initialized_error string = "database not initialized"
	query_is_empty_error     string = "query is empty"
	page_is_empty_error      string = "page is empty"
)

type sqlDBObserver struct{}

// Close method called when shutdown signal
func (o sqlDBObserver) Close() {
	logging.Info("closing database connection")
	if err := instance.Close(); err != nil {
		logging.Error("error when closing database connection: %v", err)
	}
}

var instance *sql.DB

// Initialize start connection with sql database and execute migration
func Initialize() {
	sqlDB, err := sql.Open(monitoring.GetSQLDBDriverName(), config.SQL_DB_CONNECTION_URI)
	if err != nil {
		logging.Fatal(db_connection_error, err)
	}

	if err = sqlDB.Ping(); err != nil {
		logging.Fatal(db_connection_error, err)
	}

	instance = sqlDB

	if err = migrations(); err != nil {
		logging.Fatal(db_migration_error, err)
	}

	observer.Attach(sqlDBObserver{})
	logging.Info(db_connection_success)
}

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

		isStruct, isTime, isNull, isSlice := reflectValueValidations(field)
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

func reflectValueValidations(value reflect.Value) (isStruct, isTime, isNull, isSlice bool) {
	isStruct = value.Kind() == reflect.Struct
	isTime = slices.Contains([]string{"time.Time", "types.IsoDate", "types.IsoTime"}, value.Type().String())
	isNull = strings.Contains(value.Type().String(), "Null")
	isSlice = value.Kind() == reflect.Slice
	return
}

func closer(o io.Closer) {
	if err := o.Close(); err != nil {
		logging.Error("could not close statement: %v", err)
	}
}
