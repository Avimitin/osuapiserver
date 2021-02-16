package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
)

var (
	db       *MySQLDataStore
	fixtures *testfixtures.Loader
)

func InitTestfixtures(dir string) error {
	var err error
	db, err = NewMySQLStore(
		"osu_test:osu_test@tcp(127.0.0.1:3306)/osu_test?charset=utf8mb4",
	)
	if err != nil {
		return fmt.Errorf("connect to database: %v", err)
	}
	err = db.CheckHealth()
	if err != nil {
		return fmt.Errorf("check health: %v", err)
	}
	fixtures, err = testfixtures.New(
		testfixtures.Database(db.db),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory(dir),
	)
	return err
}

func prepareTestDatabase(t testing.TB) {
	t.Helper()
	if err := fixtures.Load(); err != nil {
		t.Fatalf("loading fixtures: %v", err)
	}
}

func TestMain(m *testing.M) {
	var err error

	err = InitTestfixtures(os.Getenv("osuapi_project_root") + "/internal/database/fixtures/")
	if err != nil {
		fatalF("init test fixtures: %v", err)
	}

	os.Exit(m.Run())
}

func fatalF(context string, args ...interface{}) {
	fmt.Printf(context, args...)
	os.Exit(1)
}