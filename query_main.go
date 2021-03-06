package seecool

import (
	"database/sql"
)

type query struct {
	queryType string
	table     string
	condition string
	order     string
	limit     string
	columns   []string
	groups    []string
	using     []string // for delete query
	keys      []string // for insert query
	values    []string // for insert query
}

func (q *query) String() string {
	switch q.queryType {
	case "SELECT", "SELECT DISTINCT":
		q.columns = columnCheck(q.columns)
		str := q.queryType + " " + arrStr(q.columns) + " FROM " + q.table
		if q.condition != "" {
			str += " WHERE " + q.condition
		}
		if q.groups != nil {
			str += " GROUP BY " + arrStr(q.groups)
		}
		if q.order != "" {
			str += q.order
		}
		if q.limit != "" {
			str += q.limit
		}
		return str
	case "DELETE":
		str := q.queryType + " FROM " + q.table
		if q.using != nil {
			str += " USING " + arrStr(q.using)
		}
		if q.condition != "" {
			str += " WHERE " + q.condition
		}
		return str
	case "INSERT":
		str := q.queryType + " INTO " + q.table
		if len(q.keys) > 0 && len(q.values) > 0 && len(q.keys) == len(q.values) {
			str += " (" + arrStr(q.keys) + ")"
			str += " VALUES "
			str += "(" + arrStr(inQuote(q.values)) + ")"
			return str
		}
		return errMissingKVQuery.Error()
	case "UPDATE":
		str := q.queryType + " " + q.table + " SET "
		lenk := len(q.keys)
		lenv := len(q.values)
		if lenk > 0 && lenv > 0 && lenk == lenv {
			for i := 0; i < lenk-1; i++ {
				if q.values[i][0] == '(' {
					str += q.keys[i] + " = " + q.values[i] + ", "
				} else {
					str += q.keys[i] + " = '" + q.values[i] + "', "
				}
			}
			if q.values[lenk-1][0] == '(' {
				str += q.keys[lenk-1] + " = " + q.values[lenk-1]
			} else {
				str += q.keys[lenk-1] + " = '" + q.values[lenk-1] + "'"
			}
		}
		if q.condition != "" {
			str += " WHERE " + q.condition
		}
		return str
	}
	return errMalformedQuery.Error()
}

func QueryJson(db *sql.DB, query *query) ([]byte, error) {
	rows, err := db.Query(query.String())
	if err != nil {
		return nil, err
	}

	json, err := JsonByte(rows)
	if err != nil {
		return nil, err
	}

	return json, nil
}
