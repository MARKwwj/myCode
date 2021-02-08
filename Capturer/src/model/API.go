package model

import "encoding/json"

// Marshal Marshal
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal Unmarshal
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// UnmarshalJSON UnmarshalJSON
func UnmarshalJSON(data interface{}, v interface{}) error {
	switch value := data.(type) {
	case []byte:
		return Unmarshal(value, v)
	default:
		data, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return Unmarshal(data, v)
	}
}
