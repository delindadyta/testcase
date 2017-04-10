package helper

import (
	"encoding/json"
)

func Deserialize(str string, res interface{}) error {
	err := json.Unmarshal([]byte(str), &res)
	return err
}
