package testcase

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestCreateTable(t *testing.T) {
	db, err := sql.Open("mysql", "root:961110@tcp(localhost:3306)/testdb")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	content := `{
		"caseid": "varchar(100)",
		"casename": "varchar(100)",
		"autorun": "boolean"
	}`
	model, err := ModelFromContent([]byte(content))
	if err != nil {
		t.Fatalf("Failed to parse model: %v", err)
	}
	err = CreateTable(db, "test_case", model)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
}
