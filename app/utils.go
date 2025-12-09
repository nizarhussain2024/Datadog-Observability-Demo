package main

import (
	"encoding/json"
)

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

