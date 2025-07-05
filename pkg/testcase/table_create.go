package testcase

import "database/sql"

// CreateTable 根据数据模型创建数据表
func CreateTable(db *sql.DB, model *Model) error {
	sql := model.NewTableCreateSQL()
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
