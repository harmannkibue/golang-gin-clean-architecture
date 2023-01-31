package utils

import (
	"bytes"
	"encoding/gob"
	"strconv"
	"strings"
)

func StringToInt32(str string) (int32, error) {
	// Convert first into an integer
	intValue, _ := strconv.Atoi(str)

	// Convert now the int to int 32
	int32Value := int32(intValue)
	return int32Value, nil
}

func Float64ToInt64(num float64) (int64, error) {
	intBalance := int(num)
	return int64(intBalance), nil
}

// Float32ToInt64 We are taking the float value and then multiplying by 100 i.e the nearest cent
func Float32ToInt64(num float32) (int64, error) {
	intBalance := int(num * 100)
	return int64(intBalance), nil
}

// Int64ToString converting int64 to string
func Int64ToString(num int64) string {
	int32val := int(num)
	stringVal := strconv.Itoa(int32val)
	return stringVal
}

// StringToInt64 converting string to int64
func StringToInt64(stringNumber string) int64 {
	num := strings.Split(stringNumber, ".")
	int32val, err := StringToInt32(string(num[0]))
	if err != nil {
		print("No error and we have a value ", int32val)
		return int64(int32val)
	}

	return int64(int32val)
}

// InterfaceToBytes -.
func InterfaceToBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
