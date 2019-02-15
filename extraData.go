package tracefall

import (
	"encoding/json"
)

// ExtraData is data tree
type ExtraData map[string]interface{}

// Set data by key
func (e *ExtraData) Set(key string, val interface{}) *ExtraData {
	(*e)[key] = val
	return e
}

// Clear make empty data struct
func (e *ExtraData) Clear() *ExtraData {
	*e = NewExtraData()
	return e
}

// Get return value by key
func (e ExtraData) Get(key string) interface{} {
	if val, ok := e[key]; ok {
		return val
	}
	return nil
}

// ToJSON get json from data
func (e ExtraData) ToJSON() []byte {
	b, err := json.Marshal(e)
	if err != nil {
		b = []byte(`{}`)
	}
	return b
}

// FromJSON set data to struct from json
func (e *ExtraData) FromJSON(b []byte) error {
	return json.Unmarshal(b, e)
}

// NewExtraData create new Data struct
func NewExtraData() ExtraData {
	return make(ExtraData)
}
