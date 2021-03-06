package seecool

import (
	"reflect"
	"strconv"
	"time"
)

func CoreEncoder(i interface{}) string {
	if i == nil {
		return "null"
	}
	v := reflect.ValueOf(i)
	t := reflect.TypeOf(i)
	switch t.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		return intEncoder(i)
	case reflect.Float32, reflect.Float64:
		return floatEncoder(i)
	case reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Uint8:
			return string(v.Bytes())
		}
	case reflect.Struct:
		switch t.String() {
		case "time.Time":
			return v.Interface().(time.Time).Format("2006-01-02 15:04:05.000000000")
		}
	}
	return "null"
}

func intEncoder(i interface{}) string {
	switch i.(type) {
	case int:
		return strconv.Itoa(i.(int))
	case int8:
		return strconv.Itoa(int(i.(int8)))
	case int16:
		return strconv.Itoa(int(i.(int16)))
	case int32:
		return strconv.Itoa(int(i.(int32)))
	case int64:
		return strconv.Itoa(int(i.(int64)))
	case uint:
		return strconv.Itoa(int(i.(uint)))
	case uint8:
		return strconv.Itoa(int(i.(uint8)))
	case uint16:
		return strconv.Itoa(int(i.(uint16)))
	case uint32:
		return strconv.Itoa(int(i.(uint32)))
	case uint64:
		return strconv.Itoa(int(i.(uint64)))
	}
	return "null"
}

func floatEncoder(i interface{}) string {
	switch i.(type) {
	case float32:
		return floatEncoderCore(float64(i.(float32)))
	case float64:
		return floatEncoderCore(i.(float64))
	}
	return "null"
}

func floatEncoderCore(val float64) string {
	return strconv.FormatFloat(val, 'e', -1, 64)
}
