package store

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/alrobwilloliver/animal-service-gin-dockertest/model"
)

func TestGetAll_sqlmock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dsn := "host=localhost port=5432 user=postgres database=postgres sslmode=disable"
	gormDb, err := gorm.Open(postgres.Dialector{
		Config: &postgres.Config{
			DriverName:           "postgres",
			DSN:                  dsn,
			PreferSimpleProtocol: true,
			WithoutReturning:     true,
			Conn:                 db,
		},
	}, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm database: %s", err)
	}

	// before we actually execute our api function, we need to expect required DB actions
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "dog").
		AddRow(2, "cat")

	// must escape special characters in regex
	mock.ExpectQuery(`SELECT \* FROM \"animals\"`).WillReturnRows(rows)

	querier := Querier{}

	animals, err := querier.GetAll(gormDb)
	if err != nil {
		t.Fatalf("failed to get animals: %s", err)
	}
	if len(animals) != 2 {
		t.Fatalf("expected 1 animal, got %d", len(animals))
	}
	if (animals)[0].ID != 1 {
		t.Fatalf("expected id %d, got %d", 1, (animals)[0].ID)
	}
	if (animals)[0].Name != "dog" {
		t.Fatalf("expected make %s, got %s", "dog", (animals)[0].Name)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreate_sqlmock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dsn := "host=localhost port=5432 user=postgres database=postgres sslmode=disable"
	gormDb, err := gorm.Open(postgres.Dialector{
		Config: &postgres.Config{
			DriverName:           "postgres",
			DSN:                  dsn,
			PreferSimpleProtocol: true,
			WithoutReturning:     true,
			Conn:                 db,
		},
	}, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm database: %s", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "dog")

	// before we actually execute our api function, we need to expect required DB actions
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "animals" ("name") VALUES ($1) RETURNING "id"`)).
		WithArgs("dog").WillReturnRows(rows)
	mock.ExpectCommit()

	querier := Querier{}

	expectedAnimal := &model.Animal{Name: "dog"}
	animal, err := querier.Create(gormDb, *expectedAnimal)
	if err != nil {
		t.Fatalf("failed to create animal: %s", err)
	}
	if animal.Name != expectedAnimal.Name {
		t.Fatalf("expected make %s, got %s", expectedAnimal.Name, animal.Name)
	}
	if animal.ID != 1 {
		t.Fatalf("expected id %d, got %d", 1, animal.ID)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
