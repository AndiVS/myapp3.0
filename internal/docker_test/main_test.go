package docker_test

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log"
	"myapp3.0/internal/model"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
)

var (
	db           *sql.DB
	mockContext  = mock.Anything
	errSomeError = errors.New("some error")
	firstC       = model.Record{
		ID:   uuid.New(),
		Name: "firstCat",
		Type: "firstType",
	}
	secondC = model.Record{
		ID:   uuid.New(),
		Name: "secondCat",
		Type: "secondType",
	}
	cats = []*model.Record{&firstC, &secondC}
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestSomething(t *testing.T) {
	// db.Query()
	id := uuid.New()
	row := db.QueryRow(mockContext,
		"INSERT INTO cats (_id, name, type) VALUES ($1, $2, $3) RETURNING _id", id, firstC.Name, firstC.Type)

	err := row.Scan(id)

	require.NoError(t, err)
}
