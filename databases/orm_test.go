package databases_test

import (
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/fernandoocampo/mockit/databases"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetPet(t *testing.T) {
	// Given
	expectedPet := databases.PetFriend{
		Name:   "uno",
		Breed:  "doberman",
		Friend: "fernando",
	}
	petID := "123-123"

	db, sqlmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error with sqlmock: %s", err)
	}
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn: db,
		// DriverName: "postgres",
	})

	expectedColumns := []string{"name", "breed", "friend"}
	expectedSQL := `SELECT "name","breed" FROM "pets" WHERE id = $1`
	sqlmock.ExpectPrepare(regexp.QuoteMeta(expectedSQL))
	sqlmock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(petID).
		WillReturnRows(
			sqlmock.NewRows(expectedColumns).
				AddRow("uno", "doberman", "fernando"),
		)

	gormDB, err := gorm.Open(dialector, &gorm.Config{PrepareStmt: true})
	if err != nil {
		t.Fatalf("unexpected error opening gorm: %s", err)
	}
	gormDB.Debug()
	databases.SetORM(gormDB)

	// When
	pet, err := databases.GetPet(petID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, &expectedPet, pet)
}

func TestGetPetAndFriend(t *testing.T) {
	// Given
	expectedPet := databases.PetFriend{
		Name:   "uno",
		Breed:  "doberman",
		Friend: "fernando",
	}
	petID := "123-123"

	db, sqlmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error with sqlmock: %s", err)
	}
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn: db,
		// DriverName: "postgres",
	})

	expectedColumns := []string{"name", "breed", "friend"}
	expectedSQL := `SELECT pets.name as name,pets.breed as breed,friends.name as friend FROM "pets" join friends on pets.friend_id = friends.id WHERE pets.id = $1`
	sqlmock.ExpectPrepare(regexp.QuoteMeta(expectedSQL))
	sqlmock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(petID).
		WillReturnRows(
			sqlmock.NewRows(expectedColumns).
				AddRow("uno", "doberman", "fernando"),
		)

	gormDB, err := gorm.Open(dialector, &gorm.Config{PrepareStmt: true})
	if err != nil {
		t.Fatalf("unexpected error opening gorm: %s", err)
	}
	gormDB.Debug()
	databases.SetORM(gormDB)

	// When
	pet, err := databases.GetPetAndFriend(petID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, &expectedPet, pet)
}
