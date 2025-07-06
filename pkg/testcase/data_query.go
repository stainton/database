package testcase

import (
	"database/sql"
	"fmt"
	"strings"
)

func QueryAllField(db *sql.DB, model *Model) ([]map[string]any, error) {
	columns := model.ColumnNames()
	sqlStr := fmt.Sprintf("SELECT %s FROM %s;", strings.Join(columns, ","), model.Name)
	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, fmt.Errorf("query all field error: %w", err)
	}
	defer rows.Close()
	results := make([]map[string]any, 0, int(1e4))
	for rows.Next() {
		scanBuffer := model.NewScanBuffer()
		if err := rows.Scan(scanBuffer...); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		results = append(results, model.UnmarshalBuffer(scanBuffer))
	}
	return results, nil
}
