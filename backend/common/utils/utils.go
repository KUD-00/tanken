package utils

import (
	"fmt"
	"log"
	"strconv"
)

func AddIfNotNil(mapInterface map[string]interface{}, key string, value interface{}) {
	if value != nil {
		mapInterface[key] = value
	}
}

func StringToInt64(s string, defaultValue int64) int64 {
	if val, err := strconv.ParseInt(s, 10, 64); err == nil {
		return val
	}
	return defaultValue
}

func StringToFloat64(s string, defaultVal float64) float64 {
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val
	}
	return defaultVal
}

func StringToBool(s string, defaultValue bool) bool {
	if val, err := strconv.ParseBool(s); err == nil {
		return val
	}
	return defaultValue
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func StringPToInt64P(value interface{}, defaultValue int64) *int64 {
	if value == nil {
		return &defaultValue
	}
	if str, ok := value.(string); ok {
		var intValue int64
		_, err := fmt.Sscan(str, &intValue)
		if err == nil {
			return &intValue
		}
	}
	return &defaultValue
}

func StringPToBoolP(value interface{}, defaultValue bool) *bool {
	if value == nil {
		return &defaultValue
	}
	if str, ok := value.(string); ok {
		var boolValue bool
		_, err := fmt.Sscan(str, &boolValue)
		if err == nil {
			return &boolValue
		}
	}
	return &defaultValue
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func StringPtr(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}
