package seecool

import "strconv"

func cond(left, op, right string) string {
	return left + " " + op + " '" + right + "'"
}

func (q *query) Cond(op1, op, op2 string) *query {
	if q.condition == "" {
		q.condition = cond(op1, op, op2)
	} else {
		q.condition += " AND " + cond(op1, op, op2)
	}
	return q
}

func (q *query) CondInt(op1 string, op string, op2 int) *query {
	return q.Cond(op1, op, strconv.Itoa(op2))
}

func (q *query) Equal(op1, op2 string) *query {
	return q.Cond(op1, "=", op2)
}

func (q *query) EqualInt(op1 string, op2i int) *query {
	return q.Equal(op1, strconv.Itoa(op2i))
}

func (q *query) NotEqual(op1, op2 string) *query {
	return q.Cond(op1, "<>", op2)
}

func (q *query) NotEqualInt(op1 string, op2i int) *query {
	return q.NotEqual(op1, strconv.Itoa(op2i))
}

func (q *query) Between(key, lowerLimit, upperLimit string) *query {
	if q.condition == "" {
		q.condition = key + " BETWEEN " + "'" + lowerLimit + "' AND '" + upperLimit + "'"
	} else {
		q.condition += " AND " + key + " BETWEEN " + "'" + lowerLimit + "' AND '" + upperLimit + "'"
	}
	return q
}

func (q *query) BetweenDate(key, lowerLimit, upperLimit string) *query {
	if q.condition == "" {
		q.condition = key + " BETWEEN " + "date '" + lowerLimit + "' AND date '" + upperLimit + "'"
	} else {
		q.condition += " AND " + key + " BETWEEN " + "date '" + lowerLimit + "' AND date '" + upperLimit + "'"
	}
	return q
}

func (q *query) Like(col, cond string) *query {
	return q.Cond(col, "LIKE", cond)
}

func (q *query) Ilike(col, cond string) *query {
	return q.Cond(col, "ILIKE", cond)
}

func (q *query) Or(op1, op, op2 string) *query {
	if q.condition == "" {
		q.condition = cond(op1, op, op2)
	} else {
		q.condition += " OR " + cond(op1, op, op2)
	}
	return q
}

// SELECT * FROM items WHERE price = '10' AND item_id IN (0, 1, 2, 3)
func (q *query) In(col string, list ...string) *query {
	lenl := len(list)
	if lenl == 0 {
		return q
	}
	elements := "("
	for i := 0; i < lenl-1; i++ {
		elements += list[i] + `, `
	}
	elements += list[lenl-1] + ")"
	if q.condition == "" {
		q.condition = col + ` IN ` + elements
	} else {
		q.condition += " AND " + col + ` IN ` + elements
	}
	return q
}
