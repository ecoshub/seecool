package seecool

import "strconv"

var errorConditionType string = "Wrong type in condition statement condition can be a string or *condition."

type query struct {
	qtype string
	cols  []string
	tbl   string
	cond  string
	grp   []string
	ord   string
	lim   string
}

func (q *query) In(col string, list ...string) *query {
	lenl := len(list)
	if lenl == 0 {
		return q
	} else {

	}
	elements := "("
	for i := 0; i < lenl-1; i++ {
		elements += list[i] + `, `
	}
	elements += list[lenl-1] + ")"
	if q.cond == "" {
		q.cond = col + ` IN ` + elements
	} else {
		q.cond += " AND " + col + ` IN ` + elements
	}
	return q
}

func (q *query) And(op1, op, op2 string) *query {
	if q.cond == "" {
		q.cond = cond(op1, op, op2)
	} else {
		q.cond += " AND " + cond(op1, op, op2)
	}
	return q
}

func (q *query) Or(op1, op, op2 string) *query {
	if q.cond == "" {
		q.cond = cond(op1, op, op2)
	} else {
		q.cond += " OR " + cond(op1, op, op2)
	}
	return q
}

func (q *query) Equal(op1, op2 string) *query {
	if q.cond == "" {
		q.cond = cond(op1, "=", op2)
	} else {
		q.cond += " AND " + cond(op1, "=", op2)
	}
	return q
}

func (q *query) EqualInt(op1 string, op2i int) *query {
	op2 := strconv.Itoa(op2i)
	return q.Equal(op1, op2)
}

func (q *query) NotEqual(op1, op2 string) *query {
	if q.cond == "" {
		q.cond = cond(op1, "<>", op2)
	} else {
		q.cond += " AND " + cond(op1, "=", op2)
	}
	return q
}

func (q *query) NotEqualInt(op1 string, op2i int) *query {
	op2 := strconv.Itoa(op2i)
	return q.NotEqual(op1, op2)
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
		q.cond += " AND " + key + " BETWEEN " + "date '" + lowerLimit + "' AND date '" + upperLimit + "'"
	}
	return q
}

func (q *query) Order(col string) *query {
	q.ord = " ORDER BY " + col
	return q
}

func (q *query) OrderDesc(col string) *query {
	q.ord = " ORDER BY " + col + " DESC"
	return q
}

func (q *query) Limit(lim string) *query {
	q.lim = " LIMIT " + lim
	return q
}

func (q *query) LimitInt(lim int) *query {
	q.lim = " LIMIT " + strconv.Itoa(lim)
	return q
}

func (q *query) LimitOff(lim, off string) *query {
	q.lim = " LIMIT " + lim + " OFFSET " + off
	return q
}

func (q *query) LimitOffInt(lim, off int) *query {
	q.lim = " LIMIT " + strconv.Itoa(lim) + " OFFSET " + strconv.Itoa(off)
	return q
}

func Select(table string, cols ...string) *query {
	q := &query{qtype: "SELECT "}
	q.tbl = table
	q.cols = cols
	return q
}

func Distinct(table string, cols ...string) *query {
	q := &query{qtype: "SELECT DISTINCT "}
	q.tbl = table
	q.cols = cols
	return q
}

func Count(table string, cols ...string) *query {
	q := &query{qtype: "SELECT "}
	q.tbl = table
	if len(cols) == 0 {
		q.cols = []string{"COUNT(*)"}
	} else {
		for _, c := range cols {
			q.cols = append(q.cols, "COUNT("+c+")")
		}
	}
	return q
}

func cond(left, op, right string) string {
	return left + " " + op + " '" + right + "'"
}

func (q *query) Cond(op1, op, op2 string) *query {
	if q.cond == "" {
		q.cond = cond(op1, op, op2)
	} else {
		q.cond += " AND " + cond(op1, op, op2)
	}
	return q
}

func (q *query) CondInt(op1 string, op string, op2 int) *query {
	q.cond = op1 + " " + op + " " + strconv.Itoa(op2)
	return q
}

func (q *query) Columns(cols ...string) *query {
	q.cols = append(q.cols, cols...)
	return q
}

func (q *query) Group(cols ...string) *query {
	q.grp = cols
	for _, c := range cols {
		exs := false
		for _, g := range q.cols {
			if c == g {
				exs = true
				break
			}
		}
		if !exs {
			q.cols = append(q.cols, c)
		}
	}
	return q
}

func (q *query) String() string {
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
	leng := len(q.grp)
	if leng != 0 {
		str += " GROUP BY "
		for i := 0; i < leng-1; i++ {
			str += q.grp[i] + ", "
		}
		str += q.grp[leng-1]
	}
	if q.ord != "" {
		str += q.ord
	}
	if q.lim != "" {
		str += q.lim
	}
	str += ";"
	return str
}
