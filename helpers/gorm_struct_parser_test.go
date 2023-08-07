package gormStructParser

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	ID   int    `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func TestMapStructFieldsToDBFields(t *testing.T) {
	testStruct := &TestStruct{}

	result, err := MapStructFieldsToDBFields(testStruct)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := map[string]string{
		"ID":   "id",
		"Name": "name",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
