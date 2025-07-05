package testcase

import (
	"encoding/json"
	"os"
	"strings"
)

type Model map[string]string

func ModelFromFile(filename string) (Model, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	md := make(Model, 1000)
	err = json.Unmarshal(content, &md)
	if err != nil {
		return nil, err
	}
	return md, nil
}

func ModelFromContent(content []byte) (Model, error) {
	md := make(Model, 1000)
	err := json.Unmarshal(content, &md)
	if err != nil {
		return nil, err
	}
	return md, nil
}

func (m Model) NewTableCreateSQL(tableName string) string {
	sb := strings.Builder{}
	sb.WriteString("CREATE TABLE IF NOT EXISTS ")
	sb.WriteString(tableName)
	sb.WriteString(" (")
	// sql := "CREATE TABLE IF NOT EXISTS " + tableName + " ("
	for k, v := range m {
		sb.WriteString(k)
		sb.WriteByte(' ')
		sb.WriteString(v)
		sb.WriteString(", ")
	}
	sb.WriteString(");")
	sql := sb.String()
	return strings.ReplaceAll(sql, ", );", ");")
}
