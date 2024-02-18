package testkit

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/twotigers93/tidb/domain"
	"github.com/twotigers93/tidb/kv"
	"github.com/twotigers93/tidb/session"
	"github.com/twotigers93/tidb/store/mockstore"
)

// CreateMockStoreAndDomain return a new mock kv.Storage and *domain.Domain.
func CreateMockStoreAndDomain(opts ...mockstore.MockTiKVStoreOption) (kv.Storage, *domain.Domain, error) {
	store, err := mockstore.NewMockStore(opts...)
	if err != nil {
		return nil, nil, err
	}
	dom, err := bootstrap(store, 500*time.Millisecond)
	if err != nil {
		return nil, nil, err
	}
	return store, dom, nil
}

func bootstrap(store kv.Storage, lease time.Duration) (*domain.Domain, error) {
	session.SetSchemaLease(lease)
	session.DisableStats4Test()
	dom, err := session.BootstrapSession(store)
	if err != nil {
		return nil, err
	}
	dom.SetStatsUpdating(true)
	return dom, nil
}

// check err is mysql 1050 error
func IsAlreadyExistsError(err error) bool {
	var e *mysql.MySQLError
	ok := errors.As(err, &e)
	if ok {
		if e.Number == 1050 {
			return true
		}
	}
	return false
}
func GetAllTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func TruncateTable(db *sql.DB, table string) error {
	_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
	return err
}

// TruncateAllTable truncate all tables
func TruncateAllTable(db *sql.DB) error {
	tables, err := GetAllTables(db)
	if err != nil {
		return err
	}
	for _, table := range tables {
		err := TruncateTable(db, table)
		if err != nil {
			return err
		}
	}
	return nil
}

// DropTable drop table
func DropTable(db *sql.DB, table string) error {
	_, err := db.Exec(fmt.Sprintf("DROP TABLE %s", table))
	return err
}

// DropAllTable drop all tables
func DropAllTable(db *sql.DB) error {
	tables, err := GetAllTables(db)
	if err != nil {
		return err
	}
	for _, table := range tables {
		err := DropTable(db, table)
		if err != nil {
			return err
		}
	}
	return nil
}

// ExecFile execute sql file
func ExecFile(db *sql.DB, file string) error {
	sqlContent, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// execute sql
	_, err = db.Exec(string(sqlContent))
	if err != nil {
		return err
	}
	return nil
}

func ExecSql(db *sql.DB, sql string) error {
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func SetUTCZone() {
	time.Local = time.UTC
}
