package testcase

import "database/sql"

// CreateTable creates a table in the database using the provided model.
func CreateTable(db *sql.DB, tableName string, model Model) error {
	sql := model.NewTableCreateSQL(tableName)
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
