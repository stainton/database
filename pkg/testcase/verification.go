package testcase

import (
	"fmt"
	"regexp"
	"strings"
)

type FieldVerifier func(string, any) error

var decimalRegex *regexp.Regexp = nil
var varcharRegex *regexp.Regexp = nil
var verifiers []FieldVerifier = nil

func init() {
	// 初始化正则表达式，用于匹配 decimal 类型
	decimalRegex, _ = regexp.Compile(`^decimal\(\d+,\d+\)$`)
	varcharRegex, _ = regexp.Compile(`^varchar\(\d+\)$`)
	verifiers = []FieldVerifier{
		verifyString,
		verifyText,
		verifyInt,
		verifyDecimal,
		verifyBoolean,
	}
}

func verifyString(fieldType string, value any) error {
	if varcharRegex == nil {
		var err error
		varcharRegex, err = regexp.Compile(`^varchar\(\d+\)$`)
		if err != nil {
			return err
		}
	}
	if !varcharRegex.MatchString(fieldType) {
		return fmt.Errorf("field type %s is not a valid varchar type", fieldType)
	}
	if _, ok := value.(string); !ok {
		return fmt.Errorf("value %v is not a valid string for field type %s", value, fieldType)
	}
	return nil
}

func verifyText(fieldType string, value any) error {
	if fieldType != "text" {
		return fmt.Errorf("field type %s is not a valid text type", fieldType)
	}
	if _, ok := value.(string); !ok {
		return fmt.Errorf("value %v is not a valid string for field type %s", value, fieldType)
	}
	return nil
}

func verifyInt(fieldType string, value any) error {
	if fieldType != "int" && fieldType != "bigint" {
		return fmt.Errorf("field type %s is not a valid integer type", fieldType)
	}
	if _, ok := value.(int); !ok {
		return fmt.Errorf("value %v is not a valid string for field type %s", value, fieldType)
	}
	return nil
}

func verifyDecimal(fieldType string, value any) error {
	if decimalRegex == nil {
		var err error
		decimalRegex, err = regexp.Compile(`^decimal\(\d+,\d+\)$`)
		if err != nil {
			return fmt.Errorf("failed to compile decimal regex: %v", err)
		}
	}
	if !decimalRegex.MatchString(fieldType) {
		return fmt.Errorf("field type %s is not a valid decimal type", fieldType)
	}
	if _, ok := value.(float64); !ok {
		return fmt.Errorf("value %v is not a valid string for field type %s", value, fieldType)
	}
	return nil
}

func verifyBoolean(fieldType string, value any) error {
	if fieldType != "boolean" {
		return fmt.Errorf("field type %s is not a valid boolean type", fieldType)
	}
	if _, ok := value.(bool); !ok {
		return fmt.Errorf("value %v is not a valid string for field type %s", value, fieldType)
	}
	return nil
}

// TODO: 区分错误返回的类型，类型不对继续往下，值不对的话不再继续后面的链子
func verifyField(fieldType string, value any) error {
	fieldType = strings.ToLower(fieldType)
	for _, fn := range verifiers {
		if err := fn(fieldType, value); err == nil {
			return nil
		}
	}
	return fmt.Errorf("unsupported field type: %s", fieldType)
}
