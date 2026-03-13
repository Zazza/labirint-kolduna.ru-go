package tests

import (
	"gamebook-backend/providers"
	"gamebook-backend/tests/helpers"
	"log"
	"os"
	"testing"

	"github.com/samber/do"
)

var (
	testInjector *do.Injector
	testDB       *helpers.TestDB
	testServer   *helpers.TestServer
)

func TestMain(m *testing.M) {
	cfg, err := helpers.LoadConfig("../..")
	if err != nil {
		log.Fatalf("❌ Configuration error: %v", err)
	}

	testInjector = do.New()

	providers.RegisterDependencies(cfg, testInjector)

	testDB = helpers.NewTestDB(testInjector)
	testServer = helpers.NewTestServer(testInjector)

	code := m.Run()

	os.Exit(code)
}
