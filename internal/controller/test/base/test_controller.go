package base

import (
	"hitalent/cmd/config"
	"hitalent/cmd/factory"
	"hitalent/cmd/service"
	"hitalent/internal/model"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gorm.io/gorm"
)

type TestContainer struct {
	service.DBDecorator
	HTTPServer *httptest.Server
	HTTPClient *http.Client
}

func PrepareTestContainer(t *testing.T) TestContainer {
	t.Helper()
	factory.InitViper()

	appCfg := config.AppConfiguration{}.NewAppConfiguration()
	dbb := service.InitORM(&appCfg.DatabaseConfig)

	t.Helper()
	mux := factory.BuildServer(dbb)

	server := httptest.NewServer(mux)
	CreateTables(dbb)

	t.Cleanup(func() {
		DropTables(dbb)
		dbb.CloseDB()
		server.Close()
	})

	return TestContainer{
		DBDecorator: service.DBDecorator{},
		HTTPServer:  server,
		HTTPClient:  InitHttpClient(),
	}
}

func GetTableName(db *gorm.DB, model interface{}) (string, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return "", err
	}
	return stmt.Schema.Table, nil
}

// InitHttpClient http клиент с предопределенными таймаутами
func InitHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: time.Second,
			}).DialContext,
			ResponseHeaderTimeout: time.Second,
		},
	}
}

func CreateTables(db service.DBDecorator) {
	if err := db.GDB().AutoMigrate(model.Department{}, model.Employee{}); err != nil {
		panic("failed to create tables: " + err.Error())
	}
}

func DropTables(db service.DBDecorator) {
	departmentTName, err := GetTableName(db.GDB(), model.Department{})
	if err != nil {
		panic("failed to get department table name: " + err.Error())
	}
	employeeTName, err := GetTableName(db.GDB(), model.Employee{})
	if err != nil {
		panic("failed to get employee table name: " + err.Error())
	}

	_, err = db.NativeDB().Exec("DROP TABLE IF EXISTS " + departmentTName + " CASCADE;")
	if err != nil {
		panic("failed to drop department table: " + err.Error())
	}
	_, err = db.NativeDB().Exec("DROP TABLE IF EXISTS " + employeeTName + " CASCADE;")
	if err != nil {
		panic("failed to drop employee table: " + err.Error())
	}
}
