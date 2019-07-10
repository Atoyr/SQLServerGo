package database

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

var db *sql.DB

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		os.Exit(1)
	}
	server := os.Getenv("DBSERVER")
	instance := os.Getenv("DBINSTANCE")
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASSWORD")
	database := os.Getenv("DATABASE")
	con := SqlConnection{}
	con.Server = server
	con.Instance = instance
	con.User = user
	con.Password = password
	con.Database = database
	d, err := sql.Open("sqlserver", con.ConnectionString())
	if err != nil {
		os.Exit(1)
	}
	db = d

	code := m.Run()

	os.Exit(code)
}

func TestSimple(t *testing.T) {
	ds, err := NewDatabases(db)
	if err != nil {
		t.Error(err)
	}
	if len(ds) == 0 {
		t.Error("database not get")
	}
	for d := range ds {
		t.Log(d)
	}
}
