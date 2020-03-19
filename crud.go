package sqljson

import (
	"database/sql"
	"errors"
	"fmt"
	"jin"
)

var entryNotExiest string = "Entry not exist."

// {"first_name":"Jim","last_name":"Carter","email":"jimmy@gmail.com","gender":"Male","ip_address":"192.168.1.108","country":"Iraq","birth_date":"1988-01-17"}
func Create(db *sql.DB, json []byte, table string) error {
	query := `insert into ` + table + ` (`
	keys := ``
	values := ``
	err := jin.IterateKeyValue(json, func(key, value []byte) bool {
		keys += string(key) + `,`
		values += `'` + string(value) + `',`
		return true
	})
	if err != nil {
		return err
	}
	keys = keys[:len(keys)-1]
	values = values[:len(values)-1]
	query += keys + `) values (` + values + `)`
	_, err = db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

// {"key":"value","key2":"value2"} with and condition
func Read(db *sql.DB, json []byte, table string) ([]byte, error) {
	query := `select * from ` + table + ` where `
	err := jin.IterateKeyValue(json, func(key, value []byte) bool {
		query += string(key) + ` = '` + string(value) + `'`
		query += ` and `
		return true
	})
	if err != nil {
		return nil, err
	}
	query = query[:len(query)-5]
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	json, err = RowsToJsonByte(rows)
	if err != nil {
		return nil, err
	}
	if string(json) == `[]` {
		return nil, errors.New(entryNotExiest)
	}
	fmt.Println(json)
	return json, nil
}

// {"unique_key":"last_name","value":"john","key":"email","new_value":"johnnyy@windowslive.com"}
func Update(db *sql.DB, json []byte, table string) error {
	key, err := jin.GetString(json, "key")
	if err != nil {
		return err
	}
	unique_key, err := jin.GetString(json, "unique_key")
	if err != nil {
		return err
	}
	value, err := jin.GetString(json, "value")
	if err != nil {
		return err
	}
	newValue, err := jin.GetString(json, "new_value")
	if err != nil {
		return err
	}
	exist, err := existCore(db, table, unique_key, value)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New(entryNotExiest)
	}
	query := `update ` + table + ` set ` + key + ` = '` + newValue + `' where ` + unique_key + ` = '` + value + `'`
	_, err = db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

// {"unique_key":"email","value":"johnnyy@windowslive.com"}
func Delete(db *sql.DB, json []byte, table string) error {
	key, err := jin.GetString(json, "key")
	if err != nil {
		return err
	}
	value, err := jin.GetString(json, "value")
	if err != nil {
		return err
	}
	exist, err := existCore(db, table, key, value)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New(entryNotExiest)
	}
	query := `delete from ` + table + ` where ` + key + ` = '` + value + `'`
	_, err = db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

// {"key":"value","key2":"value2"} with regex any condition
func Find(db *sql.DB, json []byte, table string) ([]byte, error) {
	query := `select * from ` + table + ` where `
	in := false
	err := jin.IterateKeyValue(json, func(key, value []byte) bool {
		query += string(key) + ` ~* '` + string(value) + `'`
		query += ` and `
		in = true
		return true
	})
	if err != nil {
		return nil, err
	}
	if in {
		// delete last 'and' statement.
		query = query[:len(query)-5]
	} else {
		return nil, errors.New(entryNotExiest)
	}
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	json, err = RowsToJsonByte(rows)
	if err != nil {
		return nil, err
	}
	return json, nil
}

// {"key":"value","key2":"value2"} with and condition
func Count(db *sql.DB, json []byte, table string) (int, error) {
	query := `select count(*) from ` + table + ` where `
	in := false
	err := jin.IterateKeyValue(json, func(key, value []byte) bool {
		query += string(key) + ` ~* '` + string(value) + `'`
		query += ` and `
		in = true
		return true
	})
	if err != nil {
		return 0, err
	}
	if in {
		// delete last 'and' statement.
		query = query[:len(query)-5]
	} else {
		fmt.Println("no entry.")
		return 0, errors.New(entryNotExiest)
	}
	rows, err := db.Query(query)
	if err != nil {
		return 0, err
	}
	json, err = RowsToJsonByte(rows)
	if err != nil {
		return 0, err
	}
	val, err := jin.GetInt(json, "0", "count")
	if err != nil {
		return 0, err
	}
	return val, nil
}

func existCore(db *sql.DB, table, unique_key, value string) (bool, error) {
	query := `select count(*) from ` + table + ` where ` + unique_key + ` = '` + value + `'`
	rows, err := db.Query(query)
	if err != nil {
		return false, err
	}
	json, err := RowsToJsonByte(rows)
	if err != nil {
		return false, err
	}
	val, err := jin.GetInt(json, "0", "count")
	if err != nil {
		return false, err
	}
	if val > 0 {
		return true, nil
	}
	return false, nil
}
