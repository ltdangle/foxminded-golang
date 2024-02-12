package model

import (
	"database/sql"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestSqlclRepo(t *testing.T) {
	// Truncate db and init repo.
	_, repo := setUpSqlcRepo(t)

	// Create new user.
	user := NewUser()
	user.Email = "testuser@domain.net"
	err := repo.SaveOrUpdate(user)
	if err != nil {
		t.Fatal(err)
	}

	user.Email = "testuser2@domain.net"
	err = repo.SaveOrUpdate(user)
	if err != nil {
		t.Fatal(err)
	}

	updatedUser := NewUser()
	updatedUser.Email = "testuser3@domain.net"
	err = repo.SaveOrUpdate(updatedUser)
	if err != nil {
		t.Fatal(err)
	}

	// Assert user is found by uuid.
	foundUser := repo.Find(user.Uuid)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Uuid, foundUser.Uuid)

	// Assert nil is returned if user is not found.
	assert.Nil(t, repo.Find("wrong uuid"))

}

func TestFinders(t *testing.T) {
	mysql, repo := setUpSqlcRepo(t)

	bytes, err := os.ReadFile("testdata/users.sql")
	if err != nil {
		t.Fatal(err)
	}
	testData := string(bytes)

	_, err = mysql.Exec(testData)
	if err != nil {
		t.Fatal("failed to insert test data")
	}

	users := repo.FindAllUsers(0, 10)
	assert.Len(t, users, 2)

	foundUser := repo.FindByEmailPass("testuser1@domain.net", "464рурукрук443рр")
	assert.NotNil(t, foundUser)
}

func setUpSqlcRepo(t *testing.T) (*sql.DB, *sqlcUserRepo) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	mysqlDsn := os.Getenv("MYSQL_DSN_TEST")
	if mysqlDsn == "" {
		t.Fatal("MYSQL_DSN_TEST env variable not set.")
	}

	// Truncate table users using native sql package.
	mysql, err := sql.Open("mysql", mysqlDsn)
	if err != nil {
		t.Fatal("failed to connect database")
	}

	query := "TRUNCATE TABLE users"
	_, err = mysql.Exec(query)
	if err != nil {
		t.Fatal("failed to truncate table users")
	}

	// Init repo.
	repo, err := NewSqlcRepo(mysqlDsn)
	if err != nil {
		t.Fatal("failed to connect database")
	}

	return mysql, repo
}
