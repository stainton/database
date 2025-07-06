package testcase

import (
	"database/sql"
	"testing"
)

func TestDataQuery(t *testing.T) {
	model := &Model{
		TableModel: map[string]string{
			"casename": "varchar(100)",
			"caseid":   "varchar(100)",
			"autorun":  "boolean",
		},
		Name:        "test_case",
		Verified:    true,
		KeyNums:     -1,
		columnNames: nil,
	}
	t.Run("test_query_all", func(t *testing.T) {
		db, err := sql.Open("mysql", "root:961110@tcp(localhost:13306)/lottery")
		if err != nil {
			t.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()
		mp, err := QueryAllField(db, model)
		if err != nil {
			t.Fatalf("QueryAllField failed: %v", err)
		}
		if len(mp) == 0 {
			t.Fatalf("Expected non-empty result, got empty")
		}
		for _, m := range mp {
			if v, ok := m["casename"]; !ok {
				t.Errorf("Expected 'casename' field in result, got %v", m)
			} else {
				if _, ok := v.(string); !ok {
					t.Errorf("Expected 'casename' to be a string, got %T", v)
				}
			}
			if _, ok := m["caseid"]; !ok {
				t.Errorf("Expected 'caseid' field in result, got %v", m)
			} else {
				if _, ok := m["caseid"].(string); !ok {
					t.Errorf("Expected 'caseid' to be a string, got %T", m["caseid"])
				}
			}
			if _, ok := m["autorun"]; !ok {
				t.Errorf("Expected 'autorun' field in result, got %v", m)
			} else {
				if _, ok := m["autorun"].(bool); !ok {
					t.Errorf("Expected 'autorun' to be a boolean, got %T", m["autorun"])
				}
			}
		}
	})

	t.Run("test_query_specify_one", func(t *testing.T) {})
}
