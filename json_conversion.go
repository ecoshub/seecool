package seecool

import (
	"database/sql"
)

// Rows to json conversion function.
func JsonByte(rows *sql.Rows) ([]byte, error) {
	val, err := Json(rows)
	if err != nil {
		return nil, err
	}
	return stringToByteArray(val), nil
}

// Rows to json conversion function.
func Json(rows *sql.Rows) (string, error) {
	coltyps, err := rows.ColumnTypes()
	if err != nil {
		return "[]", err
	}

	lenc := len(coltyps)
	values := make([]interface{}, lenc)
	referance := make([]interface{}, lenc)
	for i := range values {
		referance[i] = &values[i]
	}

	json := "["
	count := false
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
			temp += formatType(CoreEncoder(values[i])) + `,`
		}
		temp += `"` + coltyps[lenc-1].Name() + `":`
		temp += formatType(CoreEncoder(values[lenc-1]))
		temp += "}"
		json += temp
		json += ","
		count = true
	}
	if count {
		json = json[:len(json)-1]
	}
	json += "]"
	return json, nil
}
