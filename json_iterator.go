package seecool

import "database/sql"

type iter struct {
	rows      *sql.Rows
	coltyps   []*sql.ColumnType
	referance []interface{}
	values    []interface{}
	lenc      int
}

// RowIter usage:
// iter, err := seecool.RowIter(rows)
// if err != nil {
// 	fmt.Println(err)
// 	return
// }
// val, err := iter.Next()
// for err == nil {
// 	fmt.Println(string(val))
// 	val, err = iter.Next()
// }
func RowIter(rows *sql.Rows) (*iter, error) {
	coltyps, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	lenc := len(coltyps)
	values := make([]interface{}, lenc)
	referance := make([]interface{}, lenc)
	for i := range values {
		referance[i] = &values[i]
	}
	return &iter{rows: rows, referance: referance, values: values, coltyps: coltyps, lenc: lenc}, nil
}

func (it *iter) Next() ([]byte, error) {
	it.rows.Next()
	err := it.rows.Scan(it.referance...)
	if err != nil {
		return nil, err
	}
	temp := "{"
	if it.lenc == 0 {
		temp = "{}"
		return []byte(temp), nil
	}
	for i := 0; i < it.lenc-1; i++ {
		temp += `"` + it.coltyps[i].Name() + `":`
		temp += formatType(CoreEncoder(it.values[i])) + `,`
	}
	temp += `"` + it.coltyps[it.lenc-1].Name() + `":`
	temp += formatType(CoreEncoder(it.values[it.lenc-1]))
	temp += "}"
	return []byte(temp), nil
}

// seecool.IterRowsJson(rows, func(value []byte) bool {
// 	fmt.Println(string(value))
// 	return true
// })
func IterRowsJson(rows *sql.Rows, callback func(value []byte) bool) error {
	coltyps, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	lenc := len(coltyps)
	values := make([]interface{}, lenc)
	referance := make([]interface{}, lenc)
	for i := range values {
		referance[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(referance...)
		if err != nil {
			return err
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
		if !callback([]byte(temp)) {
			return nil
		}
	}
	return nil
}
