package sql

import (
	"reflect"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

// Queryable ...
type Queryable interface {
	Reset()
	String() string
	ProtoMessage()
}

func parseTag(fieldName, tag string) []string {
	parts := strings.Split(tag, ",")
	if len(parts) == 0 {
		return []string{fieldName}
	}
	return parts
}

// GeneratePGQuery ...
func GeneratePGQuery(tableName string, message Queryable) {

	var columns []string
	var values []interface{}
	//var key []*field
	v := reflect.Indirect(reflect.ValueOf(message))
	t := v.Type()
	count := t.NumField()
	//	keys := make([]*field, 0, 2)

	c := 1
	for i := 0; i < count; i++ {
		f := t.Field(i)
		// Skip fields with no tag.
		if len(f.Tag) == 0 {
			continue
		}
		jsonTag := f.Tag.Get("json")
		if len(jsonTag) == 0 {
			continue
		}

		parts := parseTag(f.Name, jsonTag)

		columns = append(columns, parts[0])
		values = append(values, c)
		c++

	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	//	insertQuery :=
	psql.Insert(tableName).Values(values...)

}
