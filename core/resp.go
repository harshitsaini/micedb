package core

import (
	"fmt"
	"strconv"
)

func decodeSimpleString(message string) (string, int, error) {
	var idx int = 0

	for {
		if message[idx] == '\r' {
			return message[:idx], idx + 2, nil
		}
		idx += 1
	}

	return "", idx + 2, nil
}

func decodeInt64(message string) (int64, int, error) {
	value, next_pos, _ := decodeSimpleString(message)
	int_value, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		// ... handle error
		panic(err)
	}

	return int_value, next_pos, nil
}

func decodeBulkString(message string) (string, int, error) {

	string_len, next_pos, err := decodeInt64(message)
	last_pos := next_pos + int(string_len)

	if err != nil {
		// ... handle error
		panic(err)
	}
	return message[next_pos:last_pos], last_pos + 2, nil
}

func decodeArray(message string) ([]interface{}, int, error) {
	array_len, base_pos, err := decodeInt64(message)
	if err != nil {
		// ... handle error
		panic(err)
	}
	var arr = make([]interface{}, int(array_len))
	for idx := 1; idx <= int(array_len); idx++ {
		var val interface{}
		var next_pos int
		switch message[base_pos] {
		case '+':
			val, next_pos, _ = decodeSimpleString(message[base_pos+1:])
		case ':':
			val, next_pos, _ = decodeInt64(message[base_pos+1:])
		case '$':
			val, next_pos, _ = decodeBulkString(message[base_pos+1:])
		case '*':
			val, next_pos, _ = decodeArray(message[base_pos+1:])
		case '-':
			val, next_pos, _ = decodeSimpleString(message[base_pos+1:])
		default:
			fmt.Println("Invalid")
		}
		base_pos = base_pos + next_pos + 1
		arr[idx-1] = val
	}

	return arr, base_pos, nil
}

func Decode(k []byte) (interface{}, error) {
	var message string = string(k)
	var val interface{}
	switch message[0] {
	case '+':
		val, _, _ = decodeSimpleString(message[1:])
	case ':':
		val, _, _ = decodeInt64(message[1:])
	case '$':
		val, _, _ = decodeBulkString(message[1:])
	case '*':
		val, _, _ = decodeArray(message[1:])
	case '-':
		val, _, _ = decodeSimpleString(message[1:])
	default:
		fmt.Println("Invalid")
	}

	return val, nil
}
