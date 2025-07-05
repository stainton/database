package testcase

import (
	"database/sql"
	"regexp"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestCreateTable(t *testing.T) {
	db, err := sql.Open("mysql", "root:961110@tcp(localhost:13306)/lottery")
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
	err = CreateTable(db, model)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
}

func TestRegexp(t *testing.T) {
	reg, err := regexp.Compile(`^decimal\(\d+,\d+\)$`)
	if err != nil {
		t.Fatalf("Failed to compile regex: %v", err)
	}
	if reg.MatchString("decimal(10,2)") {
		t.Logf("Regex matched: decimal(10,2)")
	} else {
		t.Errorf("Regex did not match: decimal(10,2)")
	}
}
