package gormStructParser

import (
	"sync"

	"gorm.io/gorm/schema"
)

func MapStructFieldsToDBFields(dest interface{}) (map[string]string, error) {
	s, err := schema.Parse(dest, &sync.Map{}, schema.NamingStrategy{})
	m := make(map[string]string)

	if err != nil {
		return m, err
	}

	for _, field := range s.Fields {
		dbName := field.DBName
		modelName := field.Name
		m[modelName] = dbName
	}
	return m, nil
}
