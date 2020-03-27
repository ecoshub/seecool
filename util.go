package seecool

import (
	"errors"
	"fmt"
	"jin"
	"os"
	"penman"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var (
	malformedEnv error = errors.New("Malformed .env file. line format must be 'key' = 'value'.")
	emptyFile    error = errors.New("File is empty.")
)

func formatType(val string) string {
	if len(val) > 0 {
		if isBool(val) {
			return val
		}
		if isInt(val) {
			if val[0] == 48 && len(val) > 1 {
				return `"` + val + `"`
			}
			return val
		}
		if isFloat(val) {
			return val
		}
		if val == "null" {
			return val
		}
		start := val[0]
		end := val[len(val)-1]
		if (start == 34 && end == 34) || (start == 91 && end == 93) || (start == 123 && end == 125) {
			return val
		}
		return `"` + val + `"`
	}
	return `""`
}

func isBool(val string) bool {
	return val == "true" || val == "false"
}

func isFloat(val string) bool {
	_, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return false
	}
	return true
}

func isInt(val string) bool {
	_, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return false
	}
	return true
}

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
			return v.Interface().(time.Time).String()
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

func stringToByteArray(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}

func CsvToJson(file string) ([]byte, error) {
	rl, err := penman.Reader(file)
	if err != nil {
		return nil, err
	}
	defer rl.Close()

	line := rl.Next()

	if line == nil {
		return nil, emptyFile
	}

	columns := strings.Split(string(line), ",")
	columScheme := jin.MakeScheme(columns...)

	arr := jin.MakeEmptyArray()
	line = rl.Next()
	for line != nil {
		cols := strings.Split(string(line), ",")
		json := columScheme.MakeJsonString(cols...)
		arr, err = jin.Add(arr, json)
		if err != nil {
			return nil, err
		}
		line = rl.Next()
	}
	return arr, nil
}

func CsvToJsonNoHeader(file string) ([]byte, error) {
	rl, err := penman.Reader(file)
	if err != nil {
		return nil, err
	}
	defer rl.Close()
	line := rl.Next()

	if line == nil {
		return nil, emptyFile
	}

	firstLine := strings.Split(string(line), ",")

	columns := make([]string, len(firstLine))
	temp := "column_"
	for i := 0; i < len(firstLine); i++ {
		columns[i] = fmt.Sprintf("%v%v", temp, i+1)
	}
	columScheme := jin.MakeScheme(columns...)
	arr := jin.MakeEmptyArray()

	cols := strings.Split(string(line), ",")
	json := columScheme.MakeJsonString(cols...)
	arr, err = jin.Add(arr, json)

	line = rl.Next()
	for line != nil {
		cols := strings.Split(string(line), ",")
		json := columScheme.MakeJsonString(cols...)
		arr, err = jin.Add(arr, json)
		if err != nil {
			return nil, err
		}
		line = rl.Next()
	}
	return arr, nil
}

func GetEnv() (map[string]string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dir = dir + penman.Sep() + ".env"

	rl, err := penman.Reader(dir)
	if err != nil {
		return nil, err
	}
	defer rl.Close()
	mp := make(map[string]string)
	line := rl.Next()
	for line != nil {
		tokens := wordSplit(string(line))
		if len(tokens) > 3 || len(tokens) < 3 {
			return nil, malformedEnv
		}
		if tokens[1] != "=" {
			return nil, malformedEnv
		}
		key := strings.TrimSpace(tokens[0])
		value := strings.TrimSpace(tokens[2])
		mp[key] = value
		line = rl.Next()
	}
	return mp, nil
}

func wordSplit(line string) []string {
	lenl := len(line)
	if lenl < 5 {
		return nil
	}
	tokens := make([]string, 0, 3)
	onWord := false
	temp := ""
	for i := 0; i < lenl; i++ {
		curr := line[i]
		if !space(curr) {
			temp += string(curr)
			if !onWord {
				onWord = true
				continue
			}
		} else {
			if onWord {
				tokens = append(tokens, temp)
				temp = ""
				onWord = false
				continue
			}
		}
	}
	if temp != "" {
		tokens = append(tokens, temp)
	}
	return tokens
}

func space(curr byte) bool {
	// space
	if curr == 32 {
		return true
	}
	// tab
	if curr == 9 {
		return true
	}
	// new line NL
	if curr == 10 {
		return true
	}
	// return CR
	if curr == 13 {
		return true
	}
	return false
}

func arrStr(arr []string) string {
	lena := len(arr)
	switch lena {
	case 0:
		return "*"
	case 1:
		return arr[0]
	default:
		str := ""
		for i := 0; i < lena-1; i++ {
			str += arr[i] + ", "
		}
		str += arr[lena-1]
		return str
	}
}

func inQuote(arr []string) []string {
	for i := 0; i < len(arr); i++ {
		if arr[i][0] != '(' {
			arr[i] = `'` + arr[i] + `'`
		}
	}
	return arr
}

func columnCheck(cols []string) []string {
	for i, c := range cols {
		prefix, rest := astrixCheck(c)
		switch prefix {
		case "":
		case "c":
			cols[i] = "COUNT(" + rest + ")"
		case "a":
			cols[i] = "AVG(" + rest + ")"
		case "mi":
			cols[i] = "MIN(" + rest + ")"
		case "ma":
			cols[i] = "MAX(" + rest + ")"
		case "s":
			cols[i] = "SUM(" + rest + ")"
		default:
			cols[i] = strings.ToUpper(prefix) + "(" + rest + ")"
		}
	}
	return cols
}

func astrixCheck(str string) (string, string) {
	for i := 0; i < len(str); i++ {
		if str[i] == 42 {
			return str[:i], str[i+1:]
		}
	}
	return "", str
}
