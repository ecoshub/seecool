package main

import (
	"database/sql"
	"fmt"
	"jin"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	_ "github.com/lib/pq"
)

type test struct {
	val *interface{}
}

func main() {
	connStr := "user=ecomain dbname=first sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := db.Query("select * from person where gender = 'Male' and birth_date > '2010-01-01'")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	json, err := rowsToJson(rows)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonb := stringToByteArray(json)

	prs, err := jin.Parse(jsonb)
	if err != nil {
		fmt.Println(err)
		return
	}
	val, err := prs.GetString("0", "country")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(val)
}

func rowsToJson(rows *sql.Rows) (string, error) {
	coltyps, err := rows.ColumnTypes()
	if err != nil {
		return "null", err
	}

	lenc := len(coltyps)
	values := make([]interface{}, lenc)
	referance := make([]interface{}, lenc)
	for i := range values {
		referance[i] = &values[i]
	}

	json := "["
	count := 0
	for rows.Next() {
		err = rows.Scan(referance...)
		if err != nil {
			return "", err
		}
		temp := "{"
		if lenc == 0 {
			temp = "{}"
			continue
		}
		for i := 0; i < lenc-1; i++ {
			temp += `"` + coltyps[i].Name() + `":`
			temp += formatType(toString(values[i])) + `,`
		}
		temp += `"` + coltyps[lenc-1].Name() + `":`
		temp += formatType(toString(values[lenc-1]))
		temp += "}"
		json += temp
		json += ","
		count++
	}
	if count > 0 {
		json = json[:len(json)-1]
	}
	json += "]"
	return json, nil
}

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

func toString(i interface{}) string {
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
