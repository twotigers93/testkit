package testkit

import (
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/twotigers93/tidb/config"
	"github.com/twotigers93/tidb/kv"
	dbServer "github.com/twotigers93/tidb/server"
	"github.com/twotigers93/tidb/sessionctx/variable"
	"github.com/twotigers93/tidb/store/mockstore"
	"github.com/twotigers93/tidb/util/logutil"
)

var (
	defaultSocket   = "/tmp/tidb-unittest-socket"
	defaultConnOpts = "charset=utf8mb4&parseTime=true&loc=UTC"
	server          *dbServer.Server
	//dbDir           = "/tmp/tidb/db"
	dbDir         = ""
	lock          sync.Mutex
	closeL        sync.Mutex
	serverLogLeve = "info"
	logFile       = "/tmp/tidb/exec.log"
	store         kv.Storage
	readOnlyUser  = "readonly"
	readOnlyPass  = "readonly"
)

const (
	defaultRetryTime = 10
)

func newTestConfig() (*config.Config, error) {
	// 如果运行时是 macos 系统
	if runtime.GOOS == "darwin" {
		// 获取当前文件目录
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		defaultSocket = fmt.Sprintf("%s/%s", dir, "tidb-unittest-socket")
	}

	cfg := config.GetGlobalConfig()
	cfg.Port = 0
	cfg.Socket = defaultSocket
	cfg.Status.ReportStatus = false
	cfg.Security.AutoTLS = false
	variable.ProcessGeneralLog.Store(true)
	cfg.Log.Level = serverLogLeve
	cfg.Log.File.Filename = logFile
	cfg.Log.Format = "json"
	err := logutil.InitLogger(cfg.Log.ToLogConfig())
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func GetDsn() string {
	dbDSN := fmt.Sprintf("%s:%s@unix(%s)/?%s", "root", "", defaultSocket, defaultConnOpts)
	return dbDSN
}

func GetDsnWithDB(db string) string {
	dbDSN := fmt.Sprintf("%s:%s@unix(%s)/%s?%s", "root", "", defaultSocket, db, defaultConnOpts)
	return dbDSN
}

func GetReadOnlyDsnWithDB(db string) string {
	dbDSN := fmt.Sprintf("%s:%s@unix(%s)/%s?%s", readOnlyUser, readOnlyPass, defaultSocket, db, defaultConnOpts)
	return dbDSN
}

func GetConn() (*sql.DB, error) {
	dbDSN := GetDsn()
	return GetConnWithDsn(dbDSN)
}

// GetConnWithDB
func GetConnWithDB(db string) (*sql.DB, error) {
	dbDSN := GetDsnWithDB(db)
	return GetConnWithDsn(dbDSN)
}

func GetConnWithDsn(dbDSN string) (*sql.DB, error) {
	var (
		dbConn *sql.DB
		err    error
	)
	for i := 0; i < defaultRetryTime; i++ {
		dbConn, err = sql.Open("mysql", dbDSN)
		if err == nil {
			return dbConn, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return dbConn, err
}

func StartServer() error {
	lock.Lock()
	defer lock.Unlock()
	if server != nil {
		return nil
	}
	// if socket exist , remove it
	if _, err := os.Stat(defaultSocket); err == nil {
		err = os.Remove(defaultSocket)
		if err != nil {
			return err
		}
	}

	cfg, err := newTestConfig()
	if err != nil {
		return err
	}
	store, _, err = CreateMockStoreAndDomain(mockstore.WithPath(dbDir))
	if err != nil {
		return err
	}
	driver := dbServer.NewTiDBDriver(store)

	server, err = dbServer.NewServer(cfg, driver)
	if err != nil {
		return err
	}
	go func() {
		err = server.Run()
		if err != nil {
			server = nil
		}
	}()
	time.Sleep(100 * time.Millisecond)
	err = Init()
	if err != nil {
		return err
	}
	return err
}

func CloseServer() {
	closeL.Lock()
	defer closeL.Unlock()
	if server != nil {
		server.Close()
		server = nil
	}
	store.Close()
}
