package testcase

import (
	"encoding/json"
	"os"
	"strings"
)

type Model struct {
	TableModel  map[string]string `json:"table-model"`
	Name        string            `json:"name"`
	Verified    bool              `json:"verified"`
	KeyNums     int               `json:"-"`
	columnNames []string          `json:"-"`
}

func ModelFromFile(filename string) (*Model, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ModelFromContent(content)
}

func ModelFromContent(content []byte) (*Model, error) {
	model := make(map[string]any, 1000)
	err := json.Unmarshal(content, &model)
	if err != nil {
		return nil, err
	}

	// 验证表名参数正确
	tn, ok := model["name"]
	tableName, typeOK := tn.(string)
	if !ok || !typeOK || tableName == "" {
		return nil, os.ErrInvalid
	}

	// 验证验证参数正确
	v, ok := model["verified"]
	verified, typeOK := v.(bool)
	if !ok || !typeOK {
		return nil, os.ErrInvalid
	}

	// 验证表模型正确
	tm, ok := model["table-model"]
	tableModel, typeOK := tm.(map[string]string)
	if !ok || !typeOK {
		return nil, os.ErrInvalid
	}

	return &Model{
		TableModel: tableModel,
		Name:       tableName,
		KeyNums:    len(tableModel),
		Verified:   verified,
	}, nil
}

// NewTableCreateSQL 生成创建表的SQL语句
func (m *Model) NewTableCreateSQL() string {
	sb := strings.Builder{}
	sb.WriteString("CREATE TABLE IF NOT EXISTS ")
	sb.WriteString(m.Name)
	sb.WriteString(" (")
	for k, v := range m.TableModel {
		sb.WriteString(k)
		sb.WriteByte(' ')
		sb.WriteString(v)
		sb.WriteString(", ")
	}
	sb.WriteString(");")
	sql := sb.String()
	return strings.ReplaceAll(sql, ", );", ");")
}

func (m *Model) ColumnNames() []string {
	if m.columnNames != nil {
		return m.columnNames
	}
	names := make([]string, 0, len(m.TableModel))
	for k := range m.TableModel {
		names = append(names, k)
	}
	m.columnNames = names
	return names
}

func (m *Model) NewScanBuffer() []any {
	buffer := make([]any, 0, len(m.TableModel))
	for _, name := range m.ColumnNames() {
		valueType := strings.ToLower(m.TableModel[name])
		if strings.Contains(valueType, "int") {
			var value int = 0
			buffer = append(buffer, &value)
			continue
		}
		if strings.Contains(valueType, "decimal") {
			var value float64 = 0.0
			buffer = append(buffer, &value)
			continue
		}
		if strings.Contains(valueType, "char") || strings.Contains(valueType, "text") {
			var value string = ""
			buffer = append(buffer, &value)
			continue
		}
		if strings.Contains(valueType, "bool") {
			var value bool = false
			buffer = append(buffer, &value)
			continue
		}
	}
	return buffer
}

func (m *Model) UnmarshalBuffer(buffer []any) map[string]any {
	data := make(map[string]any, len(m.TableModel))
	for i, name := range m.ColumnNames() {
		value := buffer[i]
		if value == nil {
			data[name] = nil
			continue
		}
		switch v := value.(type) {
		case *int:
			data[name] = *v
		case *float64:
			data[name] = *v
		case *string:
			data[name] = *v
		case *bool:
			data[name] = *v
		default:
			data[name] = v // 其他类型直接使用
		}
	}
	return data
}
