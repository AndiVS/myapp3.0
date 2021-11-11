package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"myapp3.0/internal/model"
	"myapp3.0/internal/repository"
	"net/http"
	"net/http/httptest"
	"os/exec"

	"os"
	"strings"
	"testing"
	"time"
)

var poll *pgxpool.Pool

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=e3cr3t",
			"POSTGRES_USER=user",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://user:e3cr3t@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseURL)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		poll, err = pgxpool.Connect(context.Background(), databaseURL)
		if err != nil {
			return err
		}
		return poll.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	/*	_, err = pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "flyway/flyway",
			Tag:        "6.3.1",
			Cmd: []string{
				"-url=jdbc:postgresql://" + hostAndPort + "/dbname",
				"-user=user",
				"-password=e3cr3t",
				"migrate",
			},
		}, func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})
		if err != nil {
			log.Fatalf("Could not start resource: %s", err)
		}*/

	cmd := exec.Command("./flyway", "-url=jdbc:postgresql://"+hostAndPort+"/dbname", "-user=user", "-password=e3cr3t", "migrate")

	cmd.Dir = "/home/andeisaldyun/flyway-8.0.2"
	err = cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

var (
	firstC = model.Cat{
		ID:   uuid.New(),
		Name: "firstCat",
		Type: "firstType",
	}
	secondC = model.Cat{
		ID:   uuid.New(),
		Name: "secondCat",
		Type: "secondType",
	}
	cats = []*model.Cat{&firstC, &secondC}

	firstU = model.User{
		Username: "firstUser",
		Password: "firstUser",
		IsAdmin:  false,
	}
	secondU = model.User{
		Username: "secondUser",
		Password: "secondUser",
		IsAdmin:  false,
	}
	users = []*model.User{&firstU, &secondU}
)

func TestPostgresRepository(t *testing.T) {
	rep := repository.NewRepository(poll)

	// InsertC
	ctx, _ := setup(http.MethodPost, &firstC)
	id, err := rep.InsertCat(ctx.Request().Context(), &firstC)
	require.NoError(t, err)
	firstC.ID = id

	ctx, _ = setup(http.MethodPost, &secondC)
	id, err = rep.InsertCat(ctx.Request().Context(), &secondC)
	require.NoError(t, err)
	secondC.ID = id

	// SelectAllC
	resa, err := rep.SelectAllCat(ctx.Request().Context())
	require.NoError(t, err)
	require.Equal(t, cats[0].ID, resa[0].ID)
	require.Equal(t, cats[1].ID, resa[1].ID)

	// SelectC
	res, err := rep.SelectCat(ctx.Request().Context(), firstC.ID)
	require.NoError(t, err)
	require.Equal(t, firstC.ID, res.ID)
	require.Equal(t, firstC.Name, res.Name)

	// UpdateC
	thirdC := model.Cat{ID: firstC.ID, Name: "thirdCat", Type: "thirdType"}
	ctx, _ = setup(http.MethodPost, thirdC)
	err = rep.UpdateCat(ctx.Request().Context(), &thirdC)
	require.NoError(t, err)
	res, err = rep.SelectCat(ctx.Request().Context(), firstC.ID)
	require.NoError(t, err)
	require.Equal(t, firstC.ID, res.ID)
	require.Equal(t, thirdC.Name, res.Name)

	// DeleteC
	err = rep.DeleteCat(ctx.Request().Context(), firstC.ID)
	require.NoError(t, err)
	_, err = rep.SelectCat(ctx.Request().Context(), firstC.ID)
	require.Error(t, err)
}

func setup(method string, body interface{}) (echo.Context, *httptest.ResponseRecorder) {
	jsonBody := ""
	if body != nil {
		jsonBody = mustEncodeJSON(body)
	}
	request := httptest.NewRequest(method, "/", strings.NewReader(jsonBody))
	if body != nil {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)
	return c, recorder
}

func mustEncodeJSON(data interface{}) string {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}
