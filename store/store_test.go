package store

import (
	"errors"
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
			"5432/tcp": {{HostPort: "5432"}},
		}
		// store the data in memory to speed up tests
		config.Tmpfs = map[string]string{
			"/var/lib/postgresql/data": "rw",
		}

	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
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
	t.Run("should successfully return 2 animals", func(t *testing.T) {
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
			t.Fatalf("expected 2 animals, got %d", len(createdAnimal))
		}
		if (createdAnimal)[0].ID != 1 {
			t.Fatalf("expected id %d, got %d", 1, (createdAnimal)[0].ID)
		}
		if (createdAnimal)[0].Name != "dog" {
			t.Fatalf("expected make %s, got %s", "dog", (createdAnimal)[0].Name)
		}
		if (createdAnimal)[1].ID != 2 {
			t.Fatalf("expected id %d, got %d", 2, (createdAnimal)[1].ID)
		}
		if (createdAnimal)[1].Name != "cat" {
			t.Fatalf("expected make %s, got %s", "cat", (createdAnimal)[1].Name)
		}

		gormDb.Migrator().DropTable(&model.Animal{})
	})
	t.Run("should fail and return no animals", func(t *testing.T) {
		dsn := "host=localhost user=postgres password=secret port=5432 sslmode=disable"
		gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to open gorm database: %s", err)
		}

		gormDb.AutoMigrate(&model.Animal{})
		// simulate error
		gormDb.Error = errors.New("failed to get animals")

		querier := Querier{}

		_, err = querier.GetAll(gormDb)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if err.Error() != "failed to get animals" {
			t.Fatalf("expected error %s, got %s", "failed to get animals", err.Error())
		}

		gormDb.Migrator().DropTable(&model.Animal{})
	})
}

func TestCreate(t *testing.T) {
	t.Run("should successfully create an animal", func(t *testing.T) {

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
	})
	t.Run("should fail and return no animals", func(t *testing.T) {
		dsn := "host=localhost user=postgres password=secret port=5432 sslmode=disable"
		gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to open gorm database: %s", err)
		}

		gormDb.AutoMigrate(&animals)
		gormDb.Create(&animals)
		// simulate error
		gormDb.Error = errors.New("failed to create animal")

		querier := Querier{}

		expectedAnimal := model.Animal{Name: "snake"}
		_, err = querier.Create(gormDb, expectedAnimal)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if err.Error() != "failed to create animal" {
			t.Fatalf("expected error %s, got %s", "failed to create animal", err.Error())
		}
		gormDb.Migrator().DropTable(&model.Animal{})
	})
}
