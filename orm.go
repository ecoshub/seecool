package sqljson

var errorConditionType string = "Wrong type in condition statement condition can be a string or *condition."

type query struct {
	start string
	qtype string
	cols  []string
	tbl   string
	cond  string
	ord   string
}

func cond(left, right, op string) string {
	return left + " " + op + " '" + right + "'"
}

func (q *query) And(op1, op2, op string) *query {
	if q.cond == "" {
		q.cond = cond(op1, op2, op)
	} else {
		q.cond += " AND " + cond(op1, op2, op)
	}
	return q
}

func (q *query) Or(op1, op2 string, op string) *query {
	if q.cond == "" {
		q.cond = cond(op1, op2, op)
	} else {
		q.cond += " OR " + cond(op1, op2, op)
	}
	return q
}

func (q *query) Between(key, lowerLimit, upperLimit string) *query {
	if q.cond == "" {
		q.cond = key + " BETWEEN " + "'" + lowerLimit + "' AND '" + upperLimit + "'"
	} else {
		q.cond += " AND " + key + " BETWEEN " + "'" + lowerLimit + "' AND '" + upperLimit + "'"
	}
	return q
}

func (q *query) BetweenDate(key, lowerLimit, upperLimit string) *query {
	if q.cond == "" {
		q.cond = key + " BETWEEN " + "date '" + lowerLimit + "' AND date '" + upperLimit + "'"
	} else {
		// err
	}
	return q
}

func (q *query) Table(tbl string) *query {
	q.tbl = tbl
	return q
}

func Select(table ...string) *query {
	q := &query{qtype: "SELECT "}
	lent := len(table)
	if lent == 1 {
		q.tbl = table[0]
	}
	if lent > 1 {
		q.tbl = table[lent-1]
		q.cols = table[:lent-1]
	}
	return q
}

func Count(table string) *query {
	q := &query{qtype: "SELECT "}
	q.cols = append(q.cols, "COUNT(*)")
	q.tbl = table
	return q
}

func (q *query) Cond(op1, op2, op string) *query {
	q.cond = cond(op1, op2, op)
	return q
}

func (q *query) ToString() (string, error) {
	str := q.qtype
	lenc := len(q.cols)
	if lenc == 0 {
		str += "*"
	} else {
		for i := 0; i < lenc-1; i++ {
			str += q.cols[i] + ", "
		}
		str += q.cols[lenc-1]
	}
	str += " FROM " + q.tbl
	if q.cond != "" {
		str += " WHERE " + q.cond
	}
	if q.ord != "" {
		str += " ORDER BY " + q.ord
	}
	str += ";"
	return str, nil
}
