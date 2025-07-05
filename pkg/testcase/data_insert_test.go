package testcase

import (
	"database/sql"
	"errors"
	"testing"
)

type insertTester struct {
	sqlStr string
	values []any
	err    error
}

func (it *insertTester) sqlEqual(sqlStr string) bool {
	return it.sqlStr == sqlStr
}

func (it *insertTester) valuesEqual(values []any) bool {
	if len(it.values) != len(values) {
		return false
	}
	for i, v := range it.values {
		if v != values[i] {
			return false
		}
	}
	return true
}

func (it *insertTester) errEqual(err error) bool {
	return (it.err == nil && err == nil) || (it.err != nil && err != nil)
}

func (it *insertTester) equal(sqlStr string, values []any, err error) bool {
	res := it.sqlEqual(sqlStr) && it.valuesEqual(values) && it.errEqual(err)
	return res
}

func TestDataInsert(t *testing.T) {
	t.Run("test_insert_sql_generation", func(t *testing.T) {
		testCases := [][]any{
			{`{}`, &insertTester{
				sqlStr: "",
				values: nil,
				err:    errors.New("empty content"),
			}},
		}
		model, _ := ModelFromContent([]byte(`{"casename":"string","autorun":"bool","caseid":"string"}`))
		for _, content := range testCases {
			sqlStr, values, err := makeInsertSQL(model, []byte(content[0].(string)))
			if !content[1].(*insertTester).equal(sqlStr, values, err) {
				t.Errorf("Expected: %v, got: %v, values: %v, err: %v", content[1], &insertTester{sqlStr, values, err}, values, err)
			}
		}
	})

	t.Run("test_insert_one_from_json", func(t *testing.T) {
		db, err := sql.Open("mysql", "root:961110@tcp(localhost:13306)/lottery")
		if err != nil {
			t.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()
		model, _ := ModelFromContent([]byte(`{"casename":"string","autorun":"bool","caseid":"string"}`))
		err = InsertOneFromJson(db, model, []byte(`{"casename":"test kubernetes","autorun":false,"caseid":"TC_1"}`))
		if err != nil {
			t.Errorf("Failed to insert data: %v", err)
		}
	})
}
