package appziyoudao

import (
	"TianlangCapturer/src/api"
	"strings"
)

// Result Result
type Result struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// GetValue GetValue
func (r *Result) GetValue(key string) *api.AnyValue {
	mDataMap := r.Data.(map[string]interface{})
	keys := strings.Split(key, ".")
	keysLen := len(keys)
	for i, key := range keys {
		if i == 0 && strings.ToLower(key) == "data" {
			continue
		}
		if value, ok := mDataMap[key]; ok {
			if i < keysLen-1 {
				mDataMap = value.(map[string]interface{})
				continue
			}
			return api.NewAnyValue(value)
		}
		break
	}
	return nil
}

// GetStringValue GetStringValue
func (r *Result) GetStringValue(key string) string {
	mDataMap := r.Data.(map[string]interface{})
	keys := strings.Split(key, ".")
	keysLen := len(keys)
	for i, key := range keys {
		if value, ok := mDataMap[key]; ok {
			if i < keysLen-1 {
				mDataMap = value.(map[string]interface{})
				continue
			}
			return value.(string)
		}
		break
	}
	return ""
}
