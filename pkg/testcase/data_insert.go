package testcase

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var decimalRegex *regexp.Regexp = nil

// verifiedMap 根据model的参数验证JSON的内容是否满足数据类型的要求
// TODO：验证的时候需要检查数据类型是不是SQL支持的数据类型
func verifiedMap(model *Model, content []byte) (map[string]any, error) {
	data := make(map[string]any, model.KeyNums)
	err := json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	if !model.Verified {
		return data, nil
	}
	for key, fieldType := range model.TableModel {
		value, ok := data[key]
		if !ok {
			return nil, fmt.Errorf("missing required field: %s", key)
		}
		switch strings.ToLower(fieldType) {
		case "varchar", "text":
			if _, ok := value.(string); !ok {
				return nil, fmt.Errorf("field %s must be a string", key)
			}
		case "int", "bigint":
			if _, ok := value.(int); !ok {
				return nil, fmt.Errorf("field %s must be an integer", key)
			}
		case "boolean":
			if _, ok := value.(bool); !ok {
				return nil, fmt.Errorf("field %s must be a boolean", key)
			}
		case "float", "double":
			if _, ok := value.(float64); !ok {
				return nil, fmt.Errorf("field %s must be a float", key)
			}
		default:
			if decimalRegex == nil {
				decimalRegex, err = regexp.Compile(`^decimal\(\d+,\d+\)$`)
				if err != nil {
					return nil, fmt.Errorf("failed to compile decimal regex: %v", err)
				}
			}
			if !decimalRegex.MatchString(fieldType) {
				return nil, fmt.Errorf("unsupported field type %s for field %s", fieldType, key)
			}
		}
	}
	return data, nil
}

// makeInsertSQL 创建用于插入数据的SQL语句
func makeInsertSQL(model *Model, content []byte) (string, []any, error) {
	mp, err := verifiedMap(model, content)
	if err != nil {
		return "", nil, err
	}
	keyNums := model.KeyNums
	if len(mp) < keyNums {
		keyNums = len(mp)
	}
	sb := strings.Builder{}
	sb.WriteString("INSERT INTO ")
	sb.WriteString(model.Name)
	sb.WriteString(" (")
	values := make([]any, 0, keyNums)
	for k, v := range mp {
		sb.WriteString(k)
		values = append(values, v)
		if len(values) < keyNums {
			sb.WriteString(",")
		}
	}
	sb.WriteString(") VALUES (")
	for i := range keyNums {
		sb.WriteByte('?')
		if i < keyNums-1 {
			sb.WriteByte(',')
		}
	}
	sb.WriteString(");")
	return sb.String(), values, nil
}

// InsertOneFromJson 从JSON提取内容插入一条到数据库
func InsertOneFromJson(db *sql.DB, model *Model, content []byte) error {
	sqlStr, values, err := makeInsertSQL(model, content)
	if err != nil {
		return err
	}
	fmt.Println("SQL Statement:", sqlStr)
	res, err := db.Exec(sqlStr, values...)
	if err != nil {
		return err
	}
	if _, err = res.LastInsertId(); err != nil {
		return err
	}
	return nil
}
