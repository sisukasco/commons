package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	data, err := json.Marshal(j)
	if err != nil {
		return driver.Value(""), err
	}
	return driver.Value(string(data)), nil
}

func (j JSONMap) Scan(src interface{}) error {
	var source []byte
	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	case nil:
		source = []byte("")
	default:
		e := fmt.Sprintf("Invalid data type for JSONMap %v", v)
		return errors.New(e)
	}

	if len(source) == 0 {
		source = []byte("{}")
	}
	return json.Unmarshal(source, &j)
}
