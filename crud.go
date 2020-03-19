package sqljson

import (
	"database/sql"
	"fmt"
	"jin"
)

var (
	jin_err     string = "Jin package"
	db_err      string = "Database"
	no_data_err string = "No entry found."
	conv_error  string = "JSON conversion error"
)

func errFunc(str, err string) string {
	return str + " error:" + err
}

// {"first_name":"Jim","last_name":"Carter","email":"jimmy@gmail.com","gender":"Male","ip_address":"192.168.1.108","country":"Iraq","birth_date":"1988-01-17"}
func Create(db *sql.DB, json []byte, table string) []byte {
	query := `insert into ` + table + ` (`
	keys := ``
	values := ``
	err := jin.IterateKeyValue(json, func(key, value []byte) bool {
		keys += string(key) + `,`
		values += `'` + string(value) + `',`
		return true
	})
	if err != nil {
		return jin.MakeJsonString([]string{"event", "error_text"}, []string{"error", errFunc(jin_err, err.Error())})
	}
	keys = keys[:len(keys)-1]
	values = values[:len(values)-1]
	query += keys + `) values (` + values + `)`
	_, err = db.Query(query)
	if err != nil {
		return jin.MakeJsonString([]string{"event", "error_text"}, []string{"error", errFunc(db_err, err.Error())})
	}
	return jin.MakeJsonString([]string{"event"}, []string{"OK"})
}

// {"key":"value","key2":"value2"} with and condition
func Read(db *sql.DB, json []byte, table string) []byte {
	query := `select * from ` + table + ` where `
	err := jin.IterateKeyValue(json, func(key, value []byte) bool {
		query += string(key) + ` = '` + string(value) + `'`
		query += ` and `
		return true
	})
	if err != nil {
		return jin.MakeJsonString([]string{"event", "error_text"}, []string{"error", errFunc(jin_err, err.Error())})
	}
	query = query[:len(query)-5]
	rows, err := db.Query(query)
	if err != nil {
		return jin.MakeJsonString([]string{"event", "error_text"}, []string{"error", errFunc(db_err, err.Error())})
	}
	json, err = RowsToJsonByte(rows)
	if err != nil {
		return jin.MakeJsonString([]string{"event", "error_text"}, []string{"error", errFunc(conv_error, err.Error())})
	}
	if string(json) == `[]` {
		return jin.MakeJsonString([]string{"event", "warning_text"}, []string{"warning", no_data_err})
	}
	return json
}

// {"unique_key":"last_name","value":"john","key":"email","new_value":"johnnyy@windowslive.com"}
func Update(db *sql.DB, json []byte, table string) bool {
	key, err := jin.GetString(json, "key")
	if err != nil {
		fmt.Println(err)
	}
	unique_key, err := jin.GetString(json, "unique_key")
	if err != nil {
		fmt.Println(err)
	}
	value, err := jin.GetString(json, "value")
	if err != nil {
		fmt.Println(err)
	}
	newValue, err := jin.GetString(json, "new_value")
	if err != nil {
		fmt.Println(err)
	}
	if !existCore(db, table, unique_key, value) {
		fmt.Println("entry not exist.")
		return false
	}
	query := `update ` + table + ` set ` + key + ` = '` + newValue + `' where ` + unique_key + ` = '` + value + `'`
	_, err = db.Query(query)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("entry updated.")
	return true
}

// {"unique_key":"email","value":"johnnyy@windowslive.com"}
func Delete(db *sql.DB, json []byte, table string) bool {
	key, err := jin.GetString(json, "key")
	if err != nil {
		fmt.Println(err)
	}
	value, err := jin.GetString(json, "value")
	if err != nil {
		fmt.Println(err)
	}
	if !existCore(db, table, key, value) {
		fmt.Println("entry not exist.")
		return false
	}
	query := `delete from ` + table + ` where ` + key + ` = '` + value + `'`
	fmt.Println(query)
	_, err = db.Query(query)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("entry deleted.")
	return true
}

// {"key":"value","key2":"value2"} with regex any condition
func Find(db *sql.DB, json []byte, table string) []byte {
	query := `select * from ` + table + ` where `
	in := false
	err := jin.IterateKeyValue(json, func(key, value []byte) bool {
		query += string(key) + ` ~* '` + string(value) + `'`
		query += ` and `
		in = true
		return true
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if in {
		// delete last 'and' statement.
		query = query[:len(query)-5]
	} else {
		fmt.Println("no entry.")
		return nil
	}
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	json, err = RowsToJsonByte(rows)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return json
}

// {"key":"value","key2":"value2"} with and condition
func Count(db *sql.DB, json []byte, table string) int {
	query := `select count(*) from ` + table + ` where `
	in := false
	err := jin.IterateKeyValue(json, func(key, value []byte) bool {
		query += string(key) + ` ~* '` + string(value) + `'`
		query += ` and `
		in = true
		return true
	})
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if in {
		// delete last 'and' statement.
		query = query[:len(query)-5]
	} else {
		fmt.Println("no entry.")
		return 0
	}
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	json, err = RowsToJsonByte(rows)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	val, err := jin.GetInt(json, "0", "count")
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if val > 0 {
		return val
	}
	return 0
}

func existCore(db *sql.DB, table, unique_key, value string) bool {
	query := `select count(*) from ` + table + ` where ` + unique_key + ` = '` + value + `'`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return false
	}
	json, err := RowsToJsonByte(rows)
	if err != nil {
		fmt.Println(err)
		return false
	}
	val, err := jin.GetInt(json, "0", "count")
	if err != nil {
		fmt.Println(err)
		return false
	}
	if val > 0 {
		return true
	}
	return false
}
