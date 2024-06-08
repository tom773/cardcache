package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

func CreateObject(raw []byte) ([]reflect.StructField, error) {
	obj := bytes.SplitN(raw, []byte(" "), 2)
	name := string(obj[0])
	rawData := obj[1]

	var schema map[string]string
	if err := json.Unmarshal(rawData, &schema); err != nil {
		fmt.Printf("Type Declerations Must Be Surrounded By Double Quotes")
		return nil, fmt.Errorf("error: %v", err)
	}

	fields := make([]reflect.StructField, 0, len(schema))
	for field, t := range schema {
		fieldType := getFieldType(t)
		fields = append(fields, reflect.StructField{
			Name: capitalize(field),
			Type: fieldType,
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, field)),
		})
	}
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Schema: \n")
	return fields, nil
}

func getFieldType(typ string) reflect.Type {
	switch typ {
	case "string":
		return reflect.TypeOf("")
	case "int":
		return reflect.TypeOf(0)
	case "float":
		return reflect.TypeOf(0.0)
	default:
		return reflect.TypeOf("")
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(bytes.ToUpper([]byte{s[0]})) + s[1:]
}
