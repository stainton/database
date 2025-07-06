package testcase

import (
	"strings"
	"testing"
)

func TestModel(t *testing.T) {
	model := &Model{
		TableModel: map[string]string{
			"case_name":    "varchar(100)",
			"case_id":      "varchar(100)",
			"auto_run":     "boolean",
			"success_rate": "decimal(10, 2)",
			"runtimes":     "bigint",
		},
		Name:        "test_table",
		Verified:    true,
		KeyNums:     -1,
		columnNames: nil,
	}
	t.Run("test_create_table_sql", func(t *testing.T) {
		sql := model.NewTableCreateSQL()
		expected := "CREATE TABLE IF NOT EXISTS test_table (case_name varchar(100), case_id varchar(100), auto_run boolean, success_rate decimal(10, 2), runtimes bigint);"
		if sql != expected {
			t.Errorf("Expected %s, got %s", expected, sql)
		}
	})
	t.Run("test_column_names", func(t *testing.T) {
		columnNames := model.ColumnNames()
		expected := []string{"case_name", "case_id", "auto_run", "success_rate", "runtimes"}
		if len(columnNames) != len(expected) {
			t.Errorf("Expected %d column names, got %d", len(expected), len(columnNames))
		}
		for i, v := range expected {
			if columnNames[i] != v {
				t.Errorf("Expected column name %s, got %s", v, columnNames[i])
			}
		}
	})
	t.Run("test_new_scan_buffer", func(t *testing.T) {
		buffer := model.NewScanBuffer()
		if len(buffer) != len(model.TableModel) {
			t.Errorf("Expected buffer length %d, got %d", len(model.TableModel), len(buffer))
		}
		for i, v := range model.ColumnNames() {
			if buffer[i] == nil {
				t.Errorf("Expected buffer[%d] to be non-nil, got nil", i)
			}
			if strings.Contains(v, "int") {
				if _, ok := buffer[i].(*int); !ok {
					t.Errorf("Expected buffer[%d] to be of type *int, got %T", i, buffer[i])
				}
				continue
			} else if strings.Contains(v, "decimal") {
				if _, ok := buffer[i].(*float64); !ok {
					t.Errorf("Expected buffer[%d] to be of type *float64, got %T", i, buffer[i])
				}
				continue
			} else if strings.Contains(v, "boolean") {
				if _, ok := buffer[i].(*bool); !ok {
					t.Errorf("Expected buffer[%d] to be of type *boolean, got %T", i, buffer[i])
				}
				continue
			} else if strings.Contains(v, "char") {
				if _, ok := buffer[i].(*string); !ok {
					t.Errorf("Expected buffer[%d] to be of type *string, got %T", i, buffer[i])
				}
				continue
			}
		}
	})
	t.Run("test_unmarshal_buffer", func(t *testing.T) {
		buffer := make([]any, len(model.TableModel))
		for i, v := range model.ColumnNames() {
			if strings.Contains(model.TableModel[v], "int") {
				tv := 1
				buffer[i] = &tv
			} else if strings.Contains(model.TableModel[v], "decimal") {
				tv := 1.23
				buffer[i] = &tv
			} else if strings.Contains(model.TableModel[v], "boolean") {
				tv := true
				buffer[i] = &tv
			} else if strings.Contains(model.TableModel[v], "char") {
				tv := "test"
				buffer[i] = &tv
			}
		}
		mp := model.UnmarshalBuffer(buffer)
		if v, ok := mp["case_name"].(string); !ok || v != "test" {
			t.Errorf("Expected case_name to be 'test', got %v", mp["case_name"])
		}
		if v, ok := mp["case_id"].(string); !ok || v != "test" {
			t.Errorf("Expected case_id to be 'test', got %v", mp["case_id"])
		}
		if v, ok := mp["auto_run"].(bool); !ok || v != true {
			t.Errorf("Expected auto_run to be true, got %v", mp["auto_run"])
		}
		if v, ok := mp["success_rate"].(float64); !ok || v != 1.23 {
			t.Errorf("Expected success_rate to be 1.23, got %v", mp["success_rate"])
		}
		if v, ok := mp["runtimes"].(int); !ok || v != 1 {
			t.Errorf("Expected runtimes to be 1, got %v", mp["runtimes"])
		}
	})
}
