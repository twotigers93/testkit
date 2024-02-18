package testkit

import (
	"database/sql"
	"fmt"
)

func Init() error {
	conn, err := GetConn()
	if err != nil {
		return err
	}
	_, err = InitMultiStatement(conn)
	if err != nil {
		return err
	}
	_, err = InitTimeZone(conn)

	if err != nil {
		return err
	}

	_, err = InitReadOnlyUser(conn)
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}

func InitTimeZone(conn *sql.DB) (sql.Result, error) {
	// SET time_zone = 'UTC';
	ret, err := conn.Exec("SET GLOBAL time_zone = 'UTC'")
	return ret, err
}

func InitMultiStatement(conn *sql.DB) (sql.Result, error) {
	ret, err := conn.Exec("SET GLOBAL tidb_multi_statement_mode='ON'")
	return ret, err
}

// init read only user
func InitReadOnlyUser(conn *sql.DB) (sql.Result, error) {
	// CREATE USER 'readonly'@'%' IDENTIFIED BY 'readonly';
	// GRANT SELECT ON *.* TO 'readonly'@'%';
	ret, err := conn.Exec(fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s'", readOnlyUser, readOnlyPass))
	if err != nil {
		return nil, err
	}
	ret, err = conn.Exec("GRANT SELECT ON *.* TO 'readonly'@'%'")
	return ret, err
}
