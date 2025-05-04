package stringlistutils

import (
	"database/sql/driver"
	"encoding/json"
	"log"
)

type SqlListString []string

func (s *SqlListString) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		log.Println("Error: expected []byte for JSON scan")
	}
	return json.Unmarshal(bytes, s)
}

func (s SqlListString) Value() (driver.Value, error) {
	return json.Marshal(s)
}
