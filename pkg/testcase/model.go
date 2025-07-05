package testcase

import (
	"encoding/json"
	"os"
	"strings"
)

type Model struct {
	TableModel map[string]string `json:"table-model"`
	Name       string            `json:"name"`
	KeyNums    int               `json:"-"`
	Verified   bool              `json:"verified"`
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
