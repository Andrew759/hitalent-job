package base

import (
	"hitalent/cmd/config"
	"hitalent/cmd/factory"
	"hitalent/cmd/service"
	"net/http/httptest"
	"testing"
)

func StartTestServer(t *testing.T) *httptest.Server {
	factory.InitViper()

	appCfg := config.AppConfiguration{}.NewAppConfiguration()
	dbb := service.InitORM(&appCfg.DatabaseConfig)

	t.Helper()
	mux := factory.BuildServer(dbb)

	server := httptest.NewServer(mux)

	t.Cleanup(func() {
		dbb.CloseDB()
		server.Close()
	})

	return server
}
