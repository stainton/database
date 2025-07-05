package testcase

import (
	"fmt"
	"testing"
)

func TestVerfication(t *testing.T) {
	fmt.Println(verifyField("varchar(100)", "test string"))
	fields := []string{
		"varchar(100)",
		"text",
		"int",
		"bigint",
		"decimal(10,2)",
		"boolean",
	}
	values := []any{
		"test string",
		"test text",
		123,
		1234567890,
		123.45,
		true,
	}
	for i, field := range fields {
		if err := verifyField(field, values[i]); err != nil {
			t.Errorf("verification failed for field %s with value %v: %v", field, values[i], err)
		} else {
			fmt.Printf("verification passed for field %s with value %v\n", field, values[i])
		}
	}
}
