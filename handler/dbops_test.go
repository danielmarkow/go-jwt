package handler

import (
	_ "github.com/mattn/go-sqlite3"
	"go-jwt/dbutils"
	"go-jwt/migrations"
	"testing"
)

func TestDbOps(t *testing.T) {
	db, dbClose := dbutils.ConnectSqlite("go-jwt-test.db")
	defer dbClose()
	defer dbutils.CleanUpSqliteDbFile("./go-jwt-test.db")

	migrations.CreateTables(db)
	appCtx := AppContext{DB: db}

	t.Run("select user when 0 rows returned", func(t *testing.T) {
		user, err := appCtx.getUserById(99999)
		if err != nil {
			t.Errorf("an error occurred retrieving the user: %s", err.Error())
		}
		if user != nil {
			t.Error("user is expected to be nil")
		}
	})

	t.Run("insert user", func(t *testing.T) {
		hashedPw, err := hashPassword("testPw")
		if err != nil {
			t.Errorf("failed to hash pw: %s \n", err.Error())
		}
		user := userIn{
			Email:    "test@test.de",
			Password: hashedPw,
		}
		err = appCtx.createUser(user)
		if err != nil {
			t.Errorf("user creation failed: %s \n", err.Error())
		}
	})
	
}
