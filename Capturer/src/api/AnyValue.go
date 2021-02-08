package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// AnyValue AnyValue
type AnyValue struct {
	mValue interface{}
	mMap   map[string]interface{}
}

// NewAnyValue NewAnyValue
func NewAnyValue(value interface{}) *AnyValue {
	return new(AnyValue).Reset(value)
}

// IsNull IsNull
func (v *AnyValue) IsNull() bool {
	return v.mValue == nil
}

// Reset Reset
func (v *AnyValue) Reset(value interface{}) *AnyValue {
	v.mValue = value
	v.mMap, _ = v.mValue.(map[string]interface{})
	return v
}

// Find Find
func (v *AnyValue) Find(key string) *AnyValue {
	if v.mMap != nil {
		if value, ok := v.mMap[key]; ok {
			return NewAnyValue(value)
		}
	}
	return nil
}

// Dump Dump
func (v *AnyValue) Dump() string {
	if v.mMap != nil {
		data, _ := json.Marshal(&v.mMap)
		return string(data)
	}
	return fmt.Sprintf("%v", v.mValue)
}

// AsString AsString
func (v *AnyValue) AsString() (result string) {
	result, _ = v.mValue.(string)
	return
}

// AsInt AsInt
func (v *AnyValue) AsInt() int {
	return int(v.mValue.(float64))
}

// AsInt32 AsInt32
func (v *AnyValue) AsInt32() int32 {
	return int32(v.mValue.(float64))
}

// AsInt64 AsInt64
func (v *AnyValue) AsInt64() int64 {
	return int64(v.mValue.(float64))
}

// AsRawSlice AsRawSlice
func (v *AnyValue) AsRawSlice() []interface{} {
	return v.mValue.([]interface{})
}

// AsSlice AsSlice
func (v *AnyValue) AsSlice() []AnyValue {
	if v.mValue == nil {
		return nil
	}
	array := v.mValue.([]interface{})
	result := make([]AnyValue, len(array))
	for i, value := range array {
		result[i].Reset(value)
	}
	return result
}

// AsStringSlice AsStringSlice
func (v *AnyValue) AsStringSlice() []string {
	if v.mValue == nil {
		return nil
	}
	array := v.mValue.([]interface{})
	result := make([]string, len(array))
	for i, value := range array {
		result[i] = strings.TrimSpace(value.(string))
	}
	return result
}
