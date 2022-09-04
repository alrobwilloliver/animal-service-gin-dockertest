package store

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/alrobwilloliver/animal-service-gin-dockertest/model"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	log "github.com/sirupsen/logrus"
)

var db *sql.DB

var animals = []model.Animal{
	{
		Name: "dog",
	},
	{
		Name: "cat",
	},
}

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=postgres",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		config.PortBindings = map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostPort: "5432"}, {HostPort: "5433"}},
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	// port = resource.GetPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://postgres:secret@%s/postgres?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestGetAll(t *testing.T) {
	dsn := "host=localhost user=postgres password=secret port=5432 sslmode=disable"
	gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm database: %s", err)
	}

	gormDb.AutoMigrate(&model.Animal{})
	var animals = []model.Animal{
		{
			Name: "dog",
		},
		{
			Name: "cat",
		},
	}
	gormDb.Create(&animals)

	querier := Querier{}

	createdAnimal, err := querier.GetAll(gormDb)
	if err != nil {
		t.Fatalf("failed to get animals: %s", err)
	}
	if len(createdAnimal) != 2 {
		t.Fatalf("expected 1 animal, got %d", len(animals))
	}
	if (createdAnimal)[0].ID != 1 {
		t.Fatalf("expected id %d, got %d", 1, (animals)[0].ID)
	}
	if (createdAnimal)[0].Name != "dog" {
		t.Fatalf("expected make %s, got %s", "dog", (animals)[0].Name)
	}
	if (createdAnimal)[1].ID != 2 {
		t.Fatalf("expected id %d, got %d", 2, (animals)[1].ID)
	}
	if (createdAnimal)[1].Name != "cat" {
		t.Fatalf("expected make %s, got %s", "cat", (animals)[1].Name)
	}

	gormDb.Migrator().DropTable(&model.Animal{})
}

func TestCreate(t *testing.T) {
	dsn := "host=localhost user=postgres password=secret port=5432 sslmode=disable"
	gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm database: %s", err)
	}

	gormDb.AutoMigrate(&animals)
	gormDb.Create(&animals)

	querier := Querier{}

	expectedAnimal := model.Animal{Name: "snake"}
	animal, err := querier.Create(gormDb, expectedAnimal)
	if err != nil {
		t.Fatalf("failed to create animal: %s", err)
	}
	if animal.Name != expectedAnimal.Name {
		t.Fatalf("expected make %s, got %s", expectedAnimal.Name, animal.Name)
	}
	if animal.ID != 3 {
		t.Fatalf("expected id %d, got %d", 3, animal.ID)
	}
	gormDb.Migrator().DropTable(&model.Animal{})
}
